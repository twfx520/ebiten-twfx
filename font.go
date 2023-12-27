package ebitentwfx

import (
	"bytes"
	_ "embed"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed FiraSans-Regular.ttf
var firaSansRegular []byte

type Font struct {
	*text.GoTextFaceSource
	fontSize int
	color    uint32
}

func NewFont(filename string, fontSize int, color uint32) (*Font, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	face, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &Font{
		GoTextFaceSource: face,
		fontSize:         fontSize,
		color:            color,
	}, nil
}

func (f *Font) Draw(screen *ebiten.Image, x, y int, str string) {

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(IntToRGBA(f.color))
	op.GeoM.Translate(float64(x), float64(y))
	text.Draw(screen, str, &text.GoTextFace{
		Source: f.GoTextFaceSource,
		Size:   float64(f.fontSize),
	}, op)
}

type FontOption struct {
	Color uint32
	Size  uint32
}

func (f *Font) DrawWidthOption(screen *ebiten.Image, x, y int, str string, opt *FontOption) {
	op := &text.DrawOptions{}
	if opt.Color != 0 {
		op.ColorScale.ScaleWithColor(IntToRGBA(opt.Color))
	} else {
		op.ColorScale.ScaleWithColor(IntToRGBA(f.color))
	}
	size := f.fontSize
	if opt.Size != 0 {
		size = int(opt.Size)
	}

	op.GeoM.Translate(float64(x), float64(y))
	text.Draw(screen, str, &text.GoTextFace{
		Source: f.GoTextFaceSource,
		Size:   float64(size),
	}, op)
}
