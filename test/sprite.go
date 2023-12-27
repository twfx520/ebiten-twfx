package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	ebitentwfx "github.com/twfx520/ebiten-twfx"
)

type Game struct {
}

var i int

func (g *Game) Update() error {

	if i > 30 {

		i = 0
	}

	i++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ebitentwfx.IntToRGBA(0xCCCCCCCC))
	if spr != nil {
		spr.Draw(screen, 100, 100)
		spr.DrawRect(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

var spr *ebitentwfx.Sprite

func main() {
	var err error
	spr, err = ebitentwfx.NewSpriteFromFile("images/spr.png")
	if err != nil {
		panic(err)
	}

	spr.SetRotate(float64(math.Pi / 3))
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Text (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
