package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"../questions"
)

/**
* Lobby
*  - Register: Channel for clients to register to the Lobby
*  - Unregister: Channel for clients to unregister from the Lobby
*  - Clients: A map of clients connected to the Lobby
*  - Broadcast: Channel for messaging all clients in Lobby
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
 * ScoreUpdate
 *  - Name
 *  - Score
 */
type ScoreUpdate struct {
	Name  string
	Score int
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

/*
 * countdownToStart
 *
 */
func (lobby *Lobby) countdownToStart() {
	// Countdown lasts 5 seconds
	for i := 5; i > 0; i-- {
		allPlayersReady := true

		// Check to make sure all of the players are still ready
		for player := range lobby.Clients {
			if !player.PublicInfo.Ready {
				allPlayersReady = false
			}
		}

		// If all of the players are ready update the amount of time left and if not reset to waiting and exit
		if allPlayersReady {
			for player := range lobby.Clients {
				player.Send(Message{Type: 3, Body: strconv.Itoa(i) + "..."})
			}
		} else {
			for player := range lobby.Clients {
				player.Send(Message{Type: 3, Body: "Waiting"})
			}
			return
		}

		// Now wait for a second
		time.Sleep(time.Second)
	}

	// Start the game
	lobby.runGame()
}

/*
 * countDown
 *
 */
func (lobby *Lobby) countDown(countdown int) {
	for i := countdown; i > 0; i-- {
		for player := range lobby.Clients {
			player.Send(Message{Type: 3, Body: strconv.Itoa(i)})
		}
		time.Sleep(time.Second)
	}
}

/*
 * runGame
 *
 */
func (lobby *Lobby) runGame() {
	// There are 5 rounds in a single game
	for rounds := 5; rounds > 0; rounds-- {
		// First grab a random question from the category
		question := questions.GetRandomQuestionInCategory(lobby.Category)
		questionString, _ := json.Marshal(question)

		// Then send everyone the question
		for player := range lobby.Clients {
			player.Send(Message{Type: 4, Body: string(questionString)})
		}

		// Countdown
		lobby.countDown(10)

		// Next everyone views the points
		var playerScores []*ScoreUpdate
		for player := range lobby.Clients {
			var tmpScore ScoreUpdate

			tmpScore.Name = player.PublicInfo.Name
			tmpScore.Score = player.PublicInfo.Points
			playerScores = append(playerScores, &tmpScore)
		}

		jsonPlayerScores, _ := json.Marshal(playerScores)
		for player := range lobby.Clients {
			player.Send(Message{Type: 5, Body: string(jsonPlayerScores)})
		}

		// Countdown till next question
		lobby.countDown(8)
	}

	// Finally the game ends
	for player := range lobby.Clients {
		player.Send(Message{Type: 6, Body: "Game Over"})
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
