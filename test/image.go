package main

import (
	"bufio"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func decode(filename string) (image.Image, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return image.Decode(bufio.NewReader(f))
}

type Game struct {
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {

	img, _, err := decode("images/spr.png")
	if err != nil {
		panic(err)
	}
	obj := ebiten.NewImageFromImage(img)

	pf, err := os.Create("12.png")
	if err != nil {
		panic(err)
	}
	// fmt.Println(111)
	err = png.Encode(pf, obj)
	if err != nil {
		panic(err)
	}
	pf.Close()

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Text (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
