package main

import (
	"image/color"
	"log"

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
}

var config Config = Config{
	size: Size{
		width:  40,
		height: 50,
	},
	fieldSize: 10,
}

var snake = []Position{Position{3, 6}}

func prepend(x []Position, y Position) []Position {
	x = append(x, Position{0, 0})
	copy(x[1:], x)
	x[0] = y
	return x
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	snake = prepend(snake, Position{snake[0].x + 1, snake[0].y})

	for _, chunk := range snake {
		image, _ := ebiten.NewImage(config.fieldSize, config.fieldSize, ebiten.FilterDefault)
		image.Fill(color.RGBA{0xff, 0, 0, 0xff})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(chunk.x*config.fieldSize), float64(chunk.y*config.fieldSize))
		screen.DrawImage(image, op)
	}

	return nil
}

func main() {
	widthInPx := config.size.width * config.fieldSize
	heightInPx := config.size.height * config.fieldSize
	if err := ebiten.Run(update, widthInPx, heightInPx, 1, "Ssssnake!"); err != nil {
		log.Fatal(err)
	}
}
