package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type roomRequest struct {
	RoomId     string `json:"RoomId"`
	PlayerName string `json:"PlayerName"`
	Game       string `json:"Game"`
}

type roomResponse struct {
	RoomId string         `json:"RoomId"`
	Score  map[string]int `json:"Score"`
}

type room struct {
	Players []string
	Score   map[string]int
}

var Rooms = make(map[string]*room)
var GameMaxPlayers = map[string]int{
	"tic_tac_toe": 2,
}

func RoomHandler(req Request) {

	var roomReq roomRequest
	if err := json.Unmarshal(req.PayLoad, &roomReq); err != nil {
		log.Println("Error unmarshalling room payload:", err)
		SendError(req.PlayerId, "Invalid room payload format.")
		return
	}

	if roomReq.Game == "" {
		log.Println("Missing Game in payload.")
		SendError(req.PlayerId, "Missing required fields.")
		return
	}

	switch req.Action {
	case "addPlayer":
		addPlayerToRoom(req.PlayerId, roomReq)
	case "deleteRoom":
		deteleRoom(roomReq.RoomId)
	case "roomScore":
		roomScore(req.PlayerId, roomReq.RoomId)
	default:
		log.Println("Unknown Room action: ", req.Action)
		SendError(req.PlayerId, "Unknown action.")
	}
}

// Player
func addPlayerToRoom(playerId string, payload roomRequest) {
	if payload.RoomId == "" || Rooms[payload.RoomId] == nil {
		newRoomId := uuid.New().String()[0:8]
		Rooms[newRoomId] = &room{Players: []string{playerId}, Score: make(map[string]int)}
		SendResponse(playerId, "success", "New room created", roomResponse{RoomId: newRoomId})
		Rooms[newRoomId].Score[playerId] = 0
	} else if len(Rooms[payload.RoomId].Players) < GameMaxPlayers[payload.Game] {
		Rooms[payload.RoomId].Players = append(Rooms[payload.RoomId].Players, playerId)
		SendResponse(playerId, "success", "Joined room", roomResponse{RoomId: payload.RoomId})
		Rooms[payload.RoomId].Score[playerId] = 0
	} else if len(Rooms[payload.RoomId].Players) >= GameMaxPlayers[payload.Game] {
		log.Println("Room full: ", payload.RoomId)
		SendError(playerId, "Room is full.")
	}
}

func deteleRoom(roomId string) {
	if Rooms[roomId] != nil {
		for _, playerId := range Rooms[roomId].Players {
			SendResponse(playerId, "success", "Room deleted", roomResponse{RoomId: roomId, Score: Rooms[roomId].Score})
			DeletePlayer(playerId)
		}
		delete(Rooms, roomId)
	}
}

// Score
func UpdateScore(playerId string, roomId string) {
	if Rooms[roomId] != nil {
		for _, player := range Rooms[roomId].Players {
			Rooms[roomId].Score[player] += 1
		}
	}
}

func roomScore(playerId string, roomId string) {
	SendResponse(playerId, "success", "Room scores", Rooms[roomId].Score)
}
