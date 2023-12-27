package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	ebitentwfx "github.com/twfx520/ebiten-twfx"
)

type Game struct {
}

var spr *ebitentwfx.Sprite

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ebitentwfx.IntToRGBA(0xCCCCCCCC))
	spr.Draw(screen, 0, 0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {

	var err error
	was, err := ebitentwfx.NewWasFromWasFile("was/spr.was")
	if err != nil {
		panic(err)
	}
	tex, _, _ := was.GetTexture(0)

	// copy
	fmt.Println(111)
	save := image.NewRGBA(tex.Bounds())
	for i := 0; i < tex.Bounds().Dx(); i++ {
		for j := 0; i < tex.Bounds().Dy(); j++ {
			save.Set(i, j, tex.At(i, j))
		}
	}
	fmt.Println(222)
	pf, err := os.Create("12.png")
	if err != nil {
		panic(err)
	}
	fmt.Println(111)
	err = png.Encode(pf, save)
	if err != nil {
		panic(err)
	}
	pf.Close()

	// spr, err = ebitentwfx.NewSpriteFromImage(tex)
	// if err != nil {
	// 	panic(err)
	// }

	// ebiten.SetWindowSize(800, 600)
	// ebiten.SetWindowTitle("Text (Ebitengine Demo)")
	// if err := ebiten.RunGame(&Game{}); err != nil {
	// 	log.Fatal(err)
	// }
}
