package main

import (
	"encoding/json"
	"log"
)

type message struct {
	RoomId     string `json:"RoomId"`
	PlayerName string `json:"PlayerName"`
	Message    string `json:"Message"`
}

func MessageHandler(req Request) {

	var messageReq message
	if err := json.Unmarshal(req.PayLoad, &messageReq); err != nil {
		log.Println("Error unmarshalling message payload:", err)
		SendError(req.PlayerId, "Invalid message payload format.")
		return
	}

	if messageReq.PlayerName == "" || messageReq.RoomId == "" || messageReq.Message == "" {
		log.Println("Missing required field in payload.")
		SendError(req.PlayerId, "Missing required fields.")
		return
	}

	switch req.Action {
	case "room":
		roomChatMessage(messageReq)
	case "global":
		globalChatMessage(messageReq)
	default:
		log.Println("Unknown Message action: ", req.Action)
		SendError(req.PlayerId, "Unknown action.")
	}
}

func roomChatMessage(messageReq message) {
	for _, playerId := range Rooms[messageReq.RoomId].Players {
		SendResponse(playerId, "info", "Room Message", messageReq)
	}
}

func globalChatMessage(messageReq message) {
	for playerId := range PlayerConnections {
		SendResponse(playerId, "info", "Global Message", messageReq)
	}
}
