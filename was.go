package ebitentwfx

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type WasHeader struct {
	Flag   uint16 // 精灵文件标志 SP 0x5053
	Len    uint16 // 文件头的长度 默认为 12
	Group  uint16 // 精灵图片的组数，即方向数
	Frame  uint16 // 每组的图片数，即帧数
	Width  uint16 // 精灵动画的宽度，单位像素
	Height uint16 // 精灵动画的高度，单位像素
	KeyX   uint16 // 精灵动画的关键位X
	KeyY   uint16 //  精灵动画的关键位Y
}

type WasFrame struct {
	KeyX   int32  //  图片的关键位X
	KeyY   int32  // 图片的关键位Y
	Width  uint32 //图片的宽度，单位像素
	Height uint32 // 图片的高度，单位像素
}

type PalRGB struct {
	R uint32
	G uint32
	B uint32
}

type PalProgram struct {
	Color [3]PalRGB
}

type WAS struct {
	WasHeader
	Data      []byte
	Len       int
	Palette16 [256]uint16
	Palette32 [256]uint32
	FrameList []uint32
}

func NewWasFromMem(data []byte, lexn int) (*WAS, error) {
	was := &WAS{
		Data: data,
		Len:  len(data),
	}

	err := was.parse()
	if err != nil {
		return nil, err
	}
	fmt.Println("was:len ", was.Len)
	return was, nil
}

func NewWasFromWasFile(filename string) (*WAS, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	was := &WAS{
		Data: bs,
		Len:  len(bs),
	}
	err = was.parse()
	if err != nil {
		return nil, err
	}
	return was, nil
}

func (w *WAS) parse() error {
	buf := bytes.NewReader(w.Data)
	var header WasHeader
	err := binary.Read(buf, binary.LittleEndian, &header)
	if err != nil {
		fmt.Println("解析失败:", err)
		return err
	}
	w.WasHeader = header

	if header.Flag == 0x5053 {
		tmp := make([]byte, 512)
		buf.Read(tmp)
		for i := 0; i < 256; i++ {
			val := binary.LittleEndian.Uint16(tmp[2*i : 2*i+2])
			w.Palette16[i] = val
		}
		w.RGB565to888()
		w.FrameList = make([]uint32, header.Group*header.Frame)
		tmp = make([]byte, 4*len(w.FrameList))
		buf.Read(tmp)
		for i := 0; i < len(w.FrameList); i++ {
			val := binary.LittleEndian.Uint32(tmp[4*i : 4*i+4])
			w.FrameList[i] = val
		}
	} else {
		return errors.New("format invalid")
	}

	return nil
}

func (w *WAS) RGB565to888() {
	for i := 0; i < 256; i++ {
		r := uint8(w.Palette16[i]>>11) & 0xFF
		g := uint8(w.Palette16[i]>>5) & 0x3F
		b := uint8(w.Palette16[i]) & 0x1F
		w.Palette32[i] = uint32(b<<3|b>>2) | uint32(g<<2|g>>4)<<8 | uint32(r<<3|r>>2)<<16 | 0xFF000000
	}
}

func (w *WAS) GetTexture(id int) (*ebiten.Image, int32, int32) {
	fmt.Println("GetTexture", id)
	if id > len(w.FrameList) {
		tex := ebiten.NewImage(int(w.Width), int(w.Height))
		return tex, 0, 0
	}
	if id < len(w.FrameList) && w.FrameList[id] > 0 {
		offset := w.FrameList[id] + 16
		buf := bytes.NewReader(w.Data)
		buf.Seek(int64(offset), io.SeekStart)
		var frame WasFrame
		binary.Read(buf, binary.LittleEndian, &frame)

		lineOffset := make([]uint32, frame.Height)
		tmp := make([]byte, frame.Height*4)
		buf.Seek(int64(offset+16), io.SeekStart)
		buf.Read(tmp)
		for i := 0; i < int(frame.Height); i++ {
			val := binary.LittleEndian.Uint32(tmp[4*i : 4*i+4])
			lineOffset[i] = val
		}

		if frame.Width > 0 && frame.Height > 0 {
			bitmap := make([]uint32, frame.Width*frame.Height)

			for h := 0; h < int(frame.Height); h++ {
				rp := offset + lineOffset[h]
				pos := h * int(frame.Width)
				posLen := pos + int(frame.Width)
				for w.Data[rp] > 0 {
					style := w.Data[rp] >> 6
					if style == 0 { // 00******
						if w.Data[rp]&0x20 == 0x20 { // 001*****
							level := w.Data[rp] & 0x1F
							rp = rp + 1
							if pos < posLen {
								bitmap[pos] = _SetAlpha(w.Palette32[int(w.Data[rp])], (level<<3 | level>>2))
								pos++
								rp++
							} else {
								break
							}
						} else { // 000*****
							repeat := w.Data[rp] & 0x1F
							rp++
							level := w.Data[rp]
							rp++
							c := _SetAlpha(w.Palette32[w.Data[rp]], (level<<3)|(level>>2))
							for i := 0; i < int(repeat); i++ {
								if pos < posLen {
									bitmap[pos] = c
									pos++
								} else {
									break
								}
							}
							rp++
						}
					} else if style == 1 { //01******
						repeat := w.Data[rp] & 0x3F
						rp++

						for i := 0; i < int(repeat); i++ {
							if pos < posLen {
								bitmap[pos] = w.Palette32[w.Data[rp]]
								pos++
								rp++
							} else {
								break
							}
						}
					} else if style == 2 { // 10******
						repeat := w.Data[rp] & 0x3F
						rp++
						c := w.Palette32[w.Data[rp]]
						for i := 0; i < int(repeat); i++ {
							if pos < posLen {
								bitmap[pos] = c
								pos++
							} else {
								break
							}
						}
						rp++
					} else if style == 3 { // 11******
						repeat := w.Data[rp] & 0x3F
						if repeat != 0 {
							for i := 0; i < int(repeat); i++ {
								if pos < posLen {
									pos++
								} else {
									break
								}
							}
						} else {
							if pos-1 >= 0 {
								if w.Data[rp+1]&0x20 == 0x20 {
									level := w.Data[rp+1] & 0x1F
									level = (level<<3 | level>>2)
									lv2 := bitmap[pos-1] >> 24
									if level > byte(lv2) {
										bitmap[pos-1] = _SetAlpha(w.Palette32[w.Data[rp+2]], level)
									}
								}
							}
							if rp+3 < uint32(w.Len) {
								rp = rp + 2
							}
						}
						rp++
					}
				}
			}
			tex := ebiten.NewImage(int(frame.Width), int(frame.Height))
			SetColorList(tex, bitmap)

			return tex, frame.KeyX, frame.KeyY
		}
	}
	tex := ebiten.NewImage(int(w.Width), int(w.Height))
	return tex, 0, 0
}

func _SetAlpha(color uint32, alpha uint8) uint32 {
	return (color & 0xFFFFFF) | (uint32(alpha) << 24)
}

func SetColorList(dst *ebiten.Image, colors []uint32) {
	fmt.Println(dst.Bounds().Dx(), dst.Bounds().Dy(), len(colors))

	n := 0
	for j := 0; j < dst.Bounds().Dy(); j++ {
		for i := 0; i < dst.Bounds().Dx(); i++ {
			dst.Set(i, j, IntToRGBA(colors[n]))
			n++
		}
	}

}
