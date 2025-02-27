package qrcode

import (
	"image"
	"image/color"
)

type StringBlock interface {
	At(x, y int) rune
	Bounds() (w int, h int)
}

type simpleStringBlock struct{}

func (b simpleStringBlock) At(x, y int) rune {
	return '█'
}

func (b simpleStringBlock) Bounds() (int, int) {
	return 2, 1
}

var SimpleStringBlock = simpleStringBlock{}

type simpleImageBlock struct {
	c color.Color
}

func (b simpleImageBlock) ColorModel() color.Model {
	return color.RGBAModel
}

func (b simpleImageBlock) At(x, y int) color.Color {
	return b.c
}

func (b simpleImageBlock) Bounds() image.Rectangle {
	return image.Rect(0, 0, 1, 1)
}

func SimpleImageBlock(c color.Color) image.Image {
	return simpleImageBlock{c: c}
}
