package main

import (
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
	"time"
	"math"
	"math/rand"

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
		width:  20,
		height: 30,
	},
	fieldSize: 20,
	fps: 6,
}
type Snake struct {
	body []Position
}

var images = map[string]*ebiten.Image{}

var snake = newSnake()
var lastFrameTime time.Time

type Target = Position
var target = newTarget()

var isGameOver = false

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

func (snake *Snake) move(direction Direction, isEaten bool) {
	if direction == LEFT {
		snake.body = prepend(snake.body, Position{snake.body[0].x - 1, snake.body[0].y})
	} else if direction == RIGHT {
		snake.body = prepend(snake.body, Position{snake.body[0].x + 1, snake.body[0].y})
	} else if direction == UP {
		snake.body = prepend(snake.body, Position{snake.body[0].x, snake.body[0].y - 1})
	} else if direction == DOWN {
		snake.body = prepend(snake.body, Position{snake.body[0].x, snake.body[0].y + 1})
	}

	if !isEaten {
		snake.body = snake.body[:len(snake.body) - 1]
	}

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

func (snake Snake) contains(pos Position) bool {
	for _, item := range snake.body {
		if item.x == pos.x && item.y == pos.y {
			return true
		}
	}

	return false
}

func (snake Snake) collides() bool {
	head := snake.body[0]
	for i, item := range snake.body {
		if i != 0 && item.x == head.x && item.y == head.y {
			return true
		}
	}

	return false
}

func handleTargetEat() bool {
	if target.x == snake.body[0].x && target.y == snake.body[0].y {
		target = newTarget()
		return true
	}

	return false
}

func gameOver(screen *ebiten.Image) {
	if isGameOver {
		imgBounds := images["game_over"].Bounds()
		scale := float64(config.size.width * config.fieldSize) / float64(imgBounds.Max.X)
		translateY := (float64(config.size.height * config.fieldSize) - float64(imgBounds.Max.Y) * scale) / 2
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(0, translateY)
		screen.DrawImage(images["game_over"], op)
	}
}

func update(screen *ebiten.Image) error {
	screen.Fill(color.RGBA{0, 0, 0xff, 0xff})

	if !isGameOver {
		handleKeyPress()

		if lastFrameTime.IsZero() {
			lastFrameTime = time.Now()
			isEaten := handleTargetEat()
			setDirection()
			snake.move(direction, isEaten)
		} else if time.Now().Sub(lastFrameTime).Milliseconds() > int64(1000/config.fps) {
			lastFrameTime = time.Now()
			isEaten := handleTargetEat()
			setDirection()
			snake.move(direction, isEaten)
			if snake.collides() {
				isGameOver = true
			}
		}
	}
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	target.render(screen)
	snake.render(screen)
	gameOver(screen)

	return nil
}

func newTarget() Target {
	randX := rand.Intn(config.size.width - 1)
	randY := rand.Intn(config.size.height - 1)
	pos := Position{randX, randY}

	if snake.contains(pos) {
		return newTarget()
	}

	return Position{randX, randY}
}

func (target Target) render(screen *ebiten.Image) {
	image, _ := ebiten.NewImage(config.fieldSize, config.fieldSize, ebiten.FilterDefault)
	image.Fill(color.RGBA{0, 0xff, 0, 0xff})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(target.x*config.fieldSize), float64(target.y*config.fieldSize))
	screen.DrawImage(image, op)
}

func init() {
	var err error
	img, _, err := ebitenutil.NewImageFromFile("game_over.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	images["game_over"] = img
}

func main() {
	widthInPx := config.size.width * config.fieldSize
	heightInPx := config.size.height * config.fieldSize
	if err := ebiten.Run(update, widthInPx, heightInPx, 1, "Ssssnake!"); err != nil {
		log.Fatal(err)
	}
}
