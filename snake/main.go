package main

import (
	"image/color"
	"log"
	"time"
	"math"

	"github.com/hajimehoshi/ebiten"
)

type Size struct {
	width  int
	height int
}

type Position struct {
	x int
	y int
}
type Config struct {
	size      Size
	fieldSize int
	fps int
}

var config Config = Config{
	size: Size{
		width:  40,
		height: 50,
	},
	fieldSize: 10,
	fps: 6,
}
type Snake struct {
	body []Position
}

var snake = newSnake()
var lastFrameTime time.Time


func newSnake() Snake {
	middleX := int(math.Ceil(float64(config.size.width) / 2))
	middleY := int(math.Ceil(float64(config.size.height) / 2))
	return Snake{[]Position{Position{middleX, middleY}, {middleX + 1, middleY},{middleX + 2, middleY}}}
}
func prepend(x []Position, y Position) []Position {
	x = append(x, Position{0, 0})
	copy(x[1:], x)
	x[0] = y
	return x
}

type Direction int
const (
	LEFT Direction = iota
	RIGHT
	UP
	DOWN
)
var direction = RIGHT
var lastPressedKey = RIGHT
var pressed = map[Direction]bool{
	LEFT: false,
	RIGHT: false,
	UP: false,
	DOWN: false,
}

func handleKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if pressed[LEFT] == false && direction != RIGHT {
			lastPressedKey = LEFT
		}
		pressed[LEFT] = true
	} else {
		pressed[LEFT] = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight)  {
		if pressed[RIGHT] == false && direction != LEFT {
			lastPressedKey = RIGHT
		}
		pressed[RIGHT] = true
	} else {
		pressed[RIGHT] = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if pressed[UP] == false && direction != DOWN {
			lastPressedKey = UP
		}
		pressed[UP] = true
	} else {
		pressed[UP] = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown)  {
		if pressed[DOWN] == false && direction != UP {
			lastPressedKey = DOWN
		}
		pressed[DOWN] = true
	} else {
		pressed[DOWN] = false
	}
}

func setDirection() {
	direction = lastPressedKey
}

func (snake *Snake) move(direction Direction) {
	if direction == LEFT {
		snake.body = prepend(snake.body, Position{snake.body[0].x - 1, snake.body[0].y})
	} else if direction == RIGHT {
		snake.body = prepend(snake.body, Position{snake.body[0].x + 1, snake.body[0].y})
	} else if direction == UP {
		snake.body = prepend(snake.body, Position{snake.body[0].x, snake.body[0].y - 1})
	} else if direction == DOWN {
		snake.body = prepend(snake.body, Position{snake.body[0].x, snake.body[0].y + 1})
	}

	snake.body = snake.body[:len(snake.body) - 1]

	if snake.body[0].x == -1 {
		snake.body[0].x = config.size.width - 1
	} else if snake.body[0].x == config.size.width {
		snake.body[0].x = 0
	}

	if snake.body[0].y == -1 {
		snake.body[0].y = config.size.height - 1
	} else if snake.body[0].y == config.size.height {
		snake.body[0].y = 0
	}
}

func (snake Snake) render(screen *ebiten.Image) {
	for _, chunk := range snake.body {
		image, _ := ebiten.NewImage(config.fieldSize, config.fieldSize, ebiten.FilterDefault)
		image.Fill(color.RGBA{0xff, 0, 0, 0xff})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(chunk.x*config.fieldSize), float64(chunk.y*config.fieldSize))
		screen.DrawImage(image, op)
	}
}

func update(screen *ebiten.Image) error {
	handleKeyPress()

	if lastFrameTime.IsZero() {
		lastFrameTime = time.Now()
		setDirection()
		snake.move(direction)
	} else if time.Now().Sub(lastFrameTime).Milliseconds() > int64(1000 / config.fps) {
		lastFrameTime = time.Now()
		setDirection()
		snake.move(direction)
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	snake.render(screen)

	return nil
}

func main() {
	widthInPx := config.size.width * config.fieldSize
	heightInPx := config.size.height * config.fieldSize
	if err := ebiten.Run(update, widthInPx, heightInPx, 1, "Ssssnake!"); err != nil {
		log.Fatal(err)
	}
}
