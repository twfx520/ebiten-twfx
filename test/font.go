package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	ebitentwfx "github.com/twfx520/ebiten-twfx"
)

type Game struct {
}

var i int
var text string
var font *ebitentwfx.Font

func (g *Game) Update() error {

	if i > 30 {

		i = 0
	}

	i++

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ebitentwfx.IntToRGBA(0xCCCCCCCC))
	tps := ebiten.ActualTPS()
	font.Draw(screen, 10, 10, fmt.Sprintf("%f", tps))
	font.Draw(screen, 10, 30, "欢迎来到德莱联盟!")
	font.DrawWidthOption(screen, 10, 50, "欢迎来到德莱联盟!", &ebitentwfx.FontOption{
		Color: 0xFF00FF00,
		Size:  18,
	})
	font.Draw(screen, 10, 70, fmt.Sprintf("%f", tps))
	font.Draw(screen, 10, 90, "欢迎来到德莱联盟!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {

	var err error
	font, err = ebitentwfx.NewFont("fonts/hkyt.ttf", 16, 0xFFFF0000)
	if err != nil {
		panic(err)
	}
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Text (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
