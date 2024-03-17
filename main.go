package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

const (
	paddleHight  = 4
	paddleWidth  = 1
	paddleSymbol = 0x2588
)

type Paddle struct {
	x, y, width, height int
}

var screen tcell.Screen
var player1 *Paddle
var player2 *Paddle

func PrintString(x, y int, str string) {
	for _, c := range str {
		screen.SetContent(x, y, c, nil, tcell.StyleDefault)
		x += +1
	}
}

func Print(x, y, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(x+c, y+r, ch, nil, tcell.StyleDefault)
		}
	}
}

func DrawState() {
	screen.Clear()
	InitGameState()

	Print(player1.x, player1.y, player1.width, paddleHight, paddleSymbol)
	Print(player2.x, player2.y, player2.width, paddleHight, paddleSymbol)

	screen.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	InitScreen()
	InitGameState()
	DrawState()

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEnter {
				screen.Fini()
				os.Exit(0)
			}
		}
	}
}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - paddleHight/2

	player1 = &Paddle{
		x:      0,
		y:      paddleStart,
		width:  paddleWidth,
		height: paddleHight,
	}
	player2 = &Paddle{
		x:      width - 1,
		y:      paddleStart,
		width:  paddleWidth,
		height: paddleHight,
	}
}
