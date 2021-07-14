package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

/*
 * Client
 *  - ID:	Client ID,
 *  - Conn: Reference to websocket connection
 *  - Lobby: Reference to lobby
 */
type Client struct {
	ID         string
	PublicInfo *ClientPublicInfo
	Conn       *websocket.Conn
	Lobby      *Lobby
	mu         sync.Mutex
}

type ClientPublicInfo struct {
	Name   string
	Ready  bool
	Points int
}

/*
 * Message
 *  - Type: 0 if bytes, 1 if string (I think)
 *  - Body: String body containing content of message
 */
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

/*
 * MessageContent
 *  - Type: 		String containing type of data (eg. textMsg)
 *  - Content:	Content struct
 */
type MessageContent struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ReadyMessage struct {
	Status bool `json:"status"`
}

/*
 * MessageToClient
 *  - Type: 		The type of response
 *  - Response: Content of the response (not always there)
 */
type MessageToClient struct {
	Type     string    `json:"type"`
	Response *Response `json:"response"`
}

/*
 * Response
 *  - Lobby:	Tells the client where it has been moved (likely lobby or game)
 */
type Response struct {
	Lobby string `json:"lobby,omitempty"`
}

/*
 * Content
 *  - TextMsg: If is of textMsg type,
 */
type Content struct {
	TextMsg string `json:"textMsg,omitempty"`
	Song    string `json:"songSearch,omitempty"`
	SongID  string `json:"songID,omitempty"`
}

type ContentClient struct {
	Client  *Client
	Content *Content
}

/*
 * CreateConversation
 *  - Participants: A string delimited by | containing a list of participants in the conversation
 *  - Name:					The name of the conversation
 */
type CreateConversation struct {
	Participants string `json:"participants"`
	Name         string `json:"name"`
}

/*
 * GetConversation
 *  - ConversationID: The hash id of the conversation
 *  - Offset:					Integer offset of range of messages you're grabbing
 *  - ClientID:				ID of the client making the get request
 */
type GetConversation struct {
	ConversationID string `json:"conversationID"`
	Offset         int    `json:"offset"`
	ClientID       string `json:"clientID"`
}

func (c *Client) Send(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn.WriteJSON(v)
}

// Read function
func (c *Client) Read() {
	defer func() {
		c.Lobby.Unregister <- c

		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}
		fmt.Println(message)

		messageContent := &MessageContent{}

		err = json.Unmarshal(p, &messageContent)
		if err != nil {
			log.Println(err)
			return
		}

		switch messageContent.Type {
		case 0:
			readyMessage := &ReadyMessage{}

			err = json.Unmarshal([]byte(messageContent.Content), &readyMessage)
			if err != nil {
				log.Println(err)
				return
			}
			c.PublicInfo.Ready = readyMessage.Status
			c.Lobby.updateClientsStatus()

			allPlayersReady := true

			for player := range c.Lobby.Clients {
				if !player.PublicInfo.Ready {
					allPlayersReady = false
				}
			}

			if allPlayersReady {
				for player := range c.Lobby.Clients {
					player.Send(Message{Type: 3, Body: "Game Starting"})
				}
			}
		}
		fmt.Println("Type:", messageContent.Type)
		fmt.Println("Content:", messageContent.Content)
	}
}
