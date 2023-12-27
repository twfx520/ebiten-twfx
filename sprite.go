package ebitentwfx

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Sprite struct {
	*ebiten.Image
	scaleX float64
	scaleY float64

	vflip  int
	hflip  int
	rotate float64

	width    int
	height   int
	x        int
	y        int
	showRect bool
}

func decode(filename string) (image.Image, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()
	return image.Decode(bufio.NewReader(f))
}

func NewSpriteFromFile(filename string) (*Sprite, error) {
	img, _, err := decode(filename)
	if err != nil {
		return nil, err
	}
	obj := ebiten.NewImageFromImage(img)

	return &Sprite{
		Image: obj,
	}, nil
}

func NewSpriteFromImage(img image.Image) (*Sprite, error) {
	if img == nil {
		return nil, errors.New("img empty")
	}
	obj := ebiten.NewImageFromImage(img)
	return &Sprite{
		Image: obj,
	}, nil
}

func NewSpriteFromMemData(data []byte) (*Sprite, error) {
	if len(data) == 0 {
		return nil, errors.New("data empty")
	}
	r := bytes.NewReader(data)
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	obj := ebiten.NewImageFromImage(img)
	return &Sprite{
		Image: obj,
	}, nil
}

func NewEmptySpirite(w, h int) (*Sprite, error) {
	obj := ebiten.NewImage(w, h)
	return &Sprite{
		Image: obj,
	}, nil
}

func (s *Sprite) Draw(screen *ebiten.Image, x, y int) {
	s.x = x
	s.y = y

	var geo ebiten.GeoM
	if s.hflip != 0 {
		geo.Scale(-1, 1)
		geo.Translate(float64(s.Bounds().Dx()), 0)
	}
	if s.vflip != 0 {
		geo.Scale(1, -1)
		geo.Translate(0, float64(s.Bounds().Dy()))
	}
	if s.rotate != 0 {
		geo.Rotate(s.rotate)
	}
	geo.Translate(float64(x), float64(y))
	if s.scaleX != 0 && s.scaleY != 0 {
		geo.Scale(s.scaleX, s.scaleY)
	}

	op := &ebiten.DrawImageOptions{
		Blend: ebiten.Blend{
			BlendFactorSourceAlpha:      1,
			BlendFactorDestinationAlpha: 1,
		},
	}
	op.GeoM = geo
	screen.DrawImage(s.Image, op)
	if s.showRect {
		vector.StrokeRect(screen, float32(x), float32(y), float32(s.Bounds().Dx()), float32(s.Bounds().Dy()), 1, IntToRGBA(0xFFFF0000), true)
	}
}

func (s *Sprite) DrawRect(screen *ebiten.Image) {
	vector.StrokeRect(screen, float32(s.x), float32(s.y), float32(s.Bounds().Dx()), float32(s.Bounds().Dy()), 1, IntToRGBA(0xFFFF0000), true)
}

func (s *Sprite) SetTexture(tex *ebiten.Image) {
	s.Image = tex
}

func (s *Sprite) GetTexture() *ebiten.Image {
	return s.Image
}

func (s *Sprite) SetScale(x, y float64) {
	s.scaleX = x
	s.scaleY = y
}

func (s *Sprite) SetFlip(h, v int) {
	s.hflip = h
	s.vflip = v
}

func (s *Sprite) SetRotate(v float64) {
	s.rotate = v
}
