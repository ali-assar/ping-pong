package main

import (
	"fmt"
	"time"

	"github.com/Ali-Assar/Games/game"
)

func main() {
	game.InitScreen()
	game.InitGameState()
	inputChan := game.InitUserInput()

	for !game.IsGameOver() {
		game.HandleUserInput(game.ReadInput(inputChan))
		game.UpdateState()
		game.DrawState()

		time.Sleep(70 * time.Millisecond)
	}

	screenWidth, screenHeight := game.Screen.Size()
	winner := game.GetWinner()
	game.PrintStringCentered(screenHeight/2-1, screenWidth/2, "Game Over")
	game.PrintStringCentered(screenHeight/2, screenWidth/2, fmt.Sprintf("%s wins...", winner))
	game.Screen.Show()
	time.Sleep(3 * time.Second)
	game.Screen.Fini()
}
