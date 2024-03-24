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
	ballSymbol   = 0x25CF

	initialBallVelocityX = 1
	initialBallVelocityY = 1
)

type GameObject struct {
	x, y, width, height  int
	yVelocity, xVelocity int
	symbol               rune
}

var screen tcell.Screen
var player1 *GameObject
var player2 *GameObject
var ball *GameObject
var debugLog string

var gameObjects []*GameObject

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

func UpdateState() {
	for i := range gameObjects {
		gameObjects[i].x += gameObjects[i].xVelocity
		gameObjects[i].y += gameObjects[i].yVelocity
	}
}

func DrawState() {
	screen.Clear()
	PrintString(0, 0, debugLog)
	for _, obj := range gameObjects {
		Print(obj.x, obj.y, obj.width, obj.height, obj.symbol)
	}
	screen.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()

	for {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()

		time.Sleep(50 * time.Millisecond)
	}
}

func HandleUserInput(key string) {
	_, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && player1.y > 0 {
		player1.y--

	} else if key == "Rune[s]" && player1.y+player1.height < screenHeight {
		//TODO: create a func called is paddle at boundaries
		player1.y++
	} else if key == "Up" && player2.y > 0 {
		player2.y--

	} else if key == "Down" && player2.y+player2.height < screenHeight {
		player2.y++
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

	player1 = &GameObject{
		x:         0,
		y:         paddleStart,
		width:     paddleWidth,
		height:    paddleHight,
		symbol:    paddleSymbol,
		yVelocity: 0,
		xVelocity: 0,
	}
	player2 = &GameObject{
		x:         width - 1,
		y:         paddleStart,
		width:     paddleWidth,
		height:    paddleHight,
		symbol:    paddleSymbol,
		yVelocity: 0,
		xVelocity: 0,
	}

	ball = &GameObject{
		x:         width / 2,
		y:         height / 2,
		width:     1,
		height:    1,
		xVelocity: initialBallVelocityX,
		yVelocity: initialBallVelocityY,
		symbol:    ballSymbol,
	}
	gameObjects = []*GameObject{
		player1, player2, ball,
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
