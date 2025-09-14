package main

import (
	"log"
)

type game struct {
	Board  [9]string
	Player string
	Turn   int
}

// Handler
func TicTacToeHandler(req Request) {

	switch req.Action {
	case "newGame":
		newGame()
	case "playerMove":
		playerMove()
	case "checkWin":
		checkWin()
	default:
		log.Println("Unknown TicTacToe action: ", req.Action)
		SendError(req.PlayerId, "Unknown action.")
	}
}

// Game
func newGame() {

}

func playerMove() {

}

func checkWin() {

}
