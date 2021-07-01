package websocket

import (
	"fmt"
	"sync"
)

var once sync.Once

// Router provides a type
type Router struct {
	Register           chan *Client
	Unregister         chan *Client
	CreateConversation chan *CreateConversation
	GetConversation    chan *GetConversation
	GetConversations   chan *Client
	Clients            map[*Client]bool
	SendMessage        chan Message
}

// variavel Global
var instance *Router

// Connect provides a singleton pattern for router
func Connect() *Router {

	once.Do(func() {
		instance = &Router{
			Register:           make(chan *Client),
			Unregister:         make(chan *Client),
			CreateConversation: make(chan *CreateConversation),
			GetConversation:    make(chan *GetConversation),
			GetConversations:   make(chan *Client),
			Clients:            make(map[*Client]bool),
			SendMessage:        make(chan Message),
		}
	})

	return instance
}

// Start function starts and handles channels
func (router *Router) Start() {
	for {
		select {
		case client := <-router.Register:
			router.Clients[client] = true
			fmt.Println("Size of Connection Router: ", len(router.Clients))
		case client := <-router.Unregister:
			delete(router.Clients, client)
			fmt.Println("Size of Connection Router: ", len(router.Clients))
			for client := range router.Clients {
				client.Send(Message{Type: 1, Body: "User Disconnected..."})
			}
		case conversation := <-router.CreateConversation:
			fmt.Println("Create Conversation")
			fmt.Println("Conversation Name", conversation.Name)
			fmt.Println("Conversation Participants: ", conversation.Participants)
		case conversation := <-router.GetConversation:
			fmt.Println("Get Conversation:")
			fmt.Println("Client ID: ", conversation.ClientID)
			fmt.Println("Conversation ID: ", conversation.ConversationID)
			fmt.Println("Offset: ", conversation.Offset)
		case client := <-router.GetConversations:
			fmt.Println("Get Conversations:")
			fmt.Println("Client ID", client.ID)
		case message := <-router.SendMessage:
			//fmt.Println("The message is: ", message)
			fmt.Println("Sending message to all clients in Router")
			for client := range router.Clients {
				if err := client.Send(message); err != nil {
					fmt.Println(err)
					continue
				}
			}
		}

	}
}

// Message function
func (r *Router) Message(sender string, msg string, conversation string) {

}

// RetrieveConversations takes in a receiver and grabs conversations w/first message
func (r *Router) RetrieveConversations(receiver string) {

}

// RetrieveMessages returns 100 messages in a conversation from index offset
func (r *Router) RetrieveMessages(receiver string, conversation string, offset int) {

}
