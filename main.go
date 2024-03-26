package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	paddleHight  = 8
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

func PrintStringCentered(x, y int, str string) {
	x = x - len(str)/2
	PrintString(x, y, str)
}

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

	if wallCollision(ball) {
		ball.yVelocity = -ball.yVelocity
	}

	if paddleCollision(ball, player1) || paddleCollision(ball, player2) {
		ball.xVelocity = -ball.xVelocity
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

func wallCollision(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return !(obj.y+obj.yVelocity >= 0 && obj.y+obj.yVelocity < screenHeight)
}

func paddleCollision(ball *GameObject, paddle *GameObject) bool {
	return ball.x+ball.xVelocity == paddle.x &&
		ball.y >= paddle.y &&
		ball.y < paddle.y+paddle.height
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()

	for !IsGameOver() {
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()

		time.Sleep(50 * time.Millisecond)
	}

	screenWidth, screenHeight := screen.Size()
	winner := GetWinner()
	PrintStringCentered(screenHeight/2-1, screenWidth/2, "Game Over")
	PrintStringCentered(screenHeight/2, screenWidth/2, fmt.Sprintf("%s wins...", winner))
	screen.Show()
	time.Sleep(3 * time.Second)
	screen.Fini()
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

func IsGameOver() bool {
	return GetWinner() != ""
}

func GetWinner() string {
	screenWidth, _ := screen.Size()
	if ball.x < 0 {
		return "Player 1"
	} else if ball.x >= screenWidth {
		return "Player 2"
	} else {
		return ""
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
