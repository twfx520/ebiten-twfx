package ebitentwfx

import "image/color"

func IntToRGBA(c uint32) color.RGBA {
	return color.RGBA{
		A: uint8(c >> 24),
		R: uint8(c >> 16),
		G: uint8(c >> 8),
		B: uint8(c),
	}
}
