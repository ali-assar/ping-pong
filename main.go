package main

import (
	"fmt"
	"os"
	"time"

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
var debugLog string

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
	PrintString(0, 0, debugLog)
	Print(player1.x, player1.y, player1.width, paddleHight, paddleSymbol)
	Print(player2.x, player2.y, player2.width, paddleHight, paddleSymbol)
	screen.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()

	for {
		DrawState()
		time.Sleep(50 * time.Millisecond)

		key := ReadInput(inputChan)
		if key == "Rune[q]" {
			screen.Fini()
			os.Exit(0)
		} else if key == "Rune[w]" {
			player1.y--

		} else if key == "Rune[s]" {
			player1.y++
		} else if key == "Up" {
			player2.y--

		} else if key == "Down" {
			player2.y++
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

func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan
}

func ReadInput(inputChan chan string) string {
	var key string
	select {
	case key = <-inputChan:
	default:
		key = ""
	}
	return key
}
