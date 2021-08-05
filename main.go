package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"trivia-backend/questions"
	"trivia-backend/stringgen"
	"trivia-backend/websocket"
)

type secretsJSON struct {
	ClientID string `json:clientID`
	SecretID string `json:secretID`
}

var (
	clientID = ""
	secretID = ""
	lobbies  = make(map[string]*websocket.Lobby)
)

func serveWs(lobby *websocket.Lobby, w http.ResponseWriter, r *http.Request, name string) {
	fmt.Println("Endpoint Hit: WebSocket")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	clientPublicInfo := &websocket.ClientPublicInfo{
		Name:   name,
		Ready:  false,
		Score:  0,
		Answer: 0,
	}

	client := &websocket.Client{
		ID:         stringgen.String(10),
		PublicInfo: clientPublicInfo,
		Conn:       conn,
		Lobby:      lobby,
	}

	lobby.Register <- client
	client.Read()
}

func setupRoutes() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/getcategories", func(w http.ResponseWriter, r *http.Request) {
		// enable CORS to allow browser to make call to API
		enableCors(&w)

		type CategoryInfo struct {
			Name        string
			Description string
		}

		var categoryNames []CategoryInfo

		categories := questions.GetCategories()

		for _, category := range categories.Categories {
			categoryInfo := CategoryInfo{Name: category.Name, Description: category.Description}
			categoryNames = append(categoryNames, categoryInfo)
		}

		categoriesBytes, _ := json.Marshal(categoryNames)

		fmt.Fprintf(w, string(categoriesBytes))
	})
	rtr.HandleFunc("/gameexists/{id}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		vars := mux.Vars(r)
		varID := vars["id"]

		if _, ok := lobbies[varID]; ok {
			fmt.Fprintf(w, "{\"Exists\":true}")
		} else {
			fmt.Fprintf(w, "{\"Exists\":false}")
		}
	})
	rtr.HandleFunc("/joingame/{id}/{name}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		vars := mux.Vars(r)
		varID := vars["id"]
		varName := vars["name"]

		if lobby, ok := lobbies[varID]; ok {
			serveWs(lobby, w, r, varName)
		} else {
			fmt.Fprintf(w, "A game with that ID does not exist")
		}
	})
	rtr.HandleFunc("/creategame/{name}/{category}", func(w http.ResponseWriter, r *http.Request) {
		// enable CORS to allow browser to make call to API
		enableCors(&w)

		vars := mux.Vars(r)
		varName := vars["name"]
		varCategory := vars["category"]

		lobbyID := stringgen.String(5)

		for _, ok := lobbies[lobbyID]; ok; {
			lobbyID = stringgen.String(5)
		}

		if lobby, ok := lobbies[lobbyID]; ok {
			fmt.Fprintf(w, "A game with that ID already exists")
		} else {
			lobby = websocket.NewLobby(lobbyID, clientID, secretID, varCategory)
			fmt.Println(clientID, secretID)
			lobbies[lobbyID] = lobby

			go lobby.Start()

			serveWs(lobby, w, r, varName)
		}
	})

	http.Handle("/", rtr)
}

func main() {
	secretsFile, _ := ioutil.ReadFile("./secrets.json")

	secrets := secretsJSON{}

	_ = json.Unmarshal(secretsFile, &secrets)
	clientID = secrets.ClientID
	secretID = secrets.SecretID

	setupRoutes()

	fmt.Println("Running Go Backend on port 8080")

	http.ListenAndServe(":8080", nil)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
