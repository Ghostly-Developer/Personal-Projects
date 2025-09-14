package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Request struct {
	PlayerId string          `json:"Player"`
	Flow     string          `json:"Flow"`
	Action   string          `json:"Action"`
	PayLoad  json.RawMessage `json:"Payload"`
}

type Response struct {
	Status  string      `json:"Status"`
	Message string      `json:"Message"`
	PayLoad interface{} `json:"Payload"`
}

var PlayerConnections = make(map[string]*websocket.Conn)
var connections = make(map[*websocket.Conn]string)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		if r.Header.Get("security-key") == "1234" {
			return true
		}
		return false
	},
}

func main() {
	http.HandleFunc("/games", handleConnections)
	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected:", err)
			break
		}

		// parse message
		var request Request
		if err := json.Unmarshal(msg, &request); err != nil {
			log.Println("Error unmarshalling message:", err)
			SendError(request.PlayerId, "Invalid request format.")
			continue
		}
		if request.PlayerId == "" || request.Flow == "" || request.Action == "" {
			log.Println("Missing required fields in request")
			SendError(request.PlayerId, "Missing required fields.")
			continue
		}

		// skip if player is mismatched
		if !AddPlayer(request.PlayerId, conn) {
			continue
		}

		go handleRequest(request)
	}
}

func handleRequest(req Request) {
	switch req.Flow {
	case "room":
		RoomHandler(req)
	case "message":
		MessageHandler(req)
	case "tictactac":
		TicTacToeHandler(req)
	default:
		{
			log.Println("Invalid flow: ", req.Flow)
			SendError(req.PlayerId, "Invalid flow.")
		}
	}
}

// Response
func SendResponse(playerId string, status string, message string, payload interface{}) {
	if playerId != "" || PlayerConnections[playerId] != nil {
		PlayerConnections[playerId].WriteJSON(Response{
			Status:  status,
			Message: message,
			PayLoad: payload,
		})
	}
}

func SendError(playerId string, message string) {
	if playerId != "" || PlayerConnections[playerId] != nil {
		PlayerConnections[playerId].WriteJSON(Response{
			Status:  "error",
			Message: message,
		})
	}
}

// Player
func AddPlayer(playerId string, conn *websocket.Conn) bool {
	if PlayerConnections[playerId] != nil && connections[conn] != playerId {
		SendError(playerId, "Player Id mistmatch.")
		return false
	} else if PlayerConnections[playerId] == nil {
		PlayerConnections[playerId] = conn
		connections[conn] = playerId
		return true
	}
	return true
}

func DeletePlayer(playerId string) {
	if PlayerConnections[playerId] != nil {
		delete(PlayerConnections, playerId)
		delete(connections, PlayerConnections[playerId])
	}
}
