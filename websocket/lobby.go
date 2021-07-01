package websocket

import (
	"encoding/json"
	"fmt"
	"log"
)

/**
* Lobby
*  - Register: Channel for clients to register to the Lobby
*  - Unregister: Channel for clients to unregister from the Lobby
*  - Clients: A map of clients connected to the Lobby
*  - Broadcast: Channel for messaging all clients in Lobby
*  - Neighbors: Potentially useless :)
 */
type Lobby struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Host       *Client
	HostName   string
	Category   string
	ID         string
	ClientID   string
	SecretID   string
}

/*
* NewLobby
* @return a generated Lobby
 */
func NewLobby(ID string, ClientID string, SecretID string, Category string) *Lobby {
	return &Lobby{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Category:   Category,
		ID:         ID,
		ClientID:   ClientID,
		SecretID:   SecretID,
	}
}

func (lobby *Lobby) updateClientsStatus() {
	// First generate an array of clients with their status
	var clientArr []*ClientPublicInfo

	for client := range lobby.Clients {
		clientArr = append(clientArr, client.PublicInfo)
	}

	jsonClients, err := json.Marshal(clientArr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Now update all of the clients
	for client := range lobby.Clients {
		client.Send(Message{Type: 1, Body: string(jsonClients)})
	}
}

func (lobby *Lobby) Start() {

	for {
		select {
		case client := <-lobby.Register:
			if len(lobby.Clients) == 0 {
				lobby.Host = client
			}

			lobby.Clients[client] = true
			client.Send(Message{Type: 0, Body: lobby.ID})
			client.Send(Message{Type: 2, Body: client.ID})

			// Now send a list of all the players to everyone
			lobby.updateClientsStatus()

			fmt.Println("Size of Connection Lobby: ", len(lobby.Clients))
			fmt.Println("Lobby ID:", lobby.ID)
		case client := <-lobby.Unregister:
			if client == lobby.Host {
				log.Println("Host unregister")

				for client := range lobby.Clients {
					client.Send(Message{Type: 1, Body: "Session Ended"})
					delete(lobby.Clients, client)
				}
			} else {
				delete(lobby.Clients, client)
				fmt.Println("Size of Connection Lobby: ", len(lobby.Clients))
			}
			// Now send a list of all the players to everyone
			lobby.updateClientsStatus()
		}

	}
}
