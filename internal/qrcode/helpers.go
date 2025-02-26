package qrcode

import (
	icolor "image/color"

	"github.com/fatih/color"
)

var (
	smallChars = [4]rune{' ', '▀', '▄', '█'}
)

func getCharOfBlockBools(invert bool, bools ...bool) uint8 {
	var idx uint8
	for i, b := range bools {
		if b {
			idx |= 1 << i
		}
	}
	if invert {
		return ^idx & 0b11
	}
	return idx
}

func createColorPrinter(fg, bg icolor.Color) *color.Color {
	fR, fG, fB, fA := fg.RGBA()
	bR, bG, bB, bA := bg.RGBA()
	c := color.New()
	if fA != 0 {
		c = c.AddRGB(int(fR>>8), int(fG>>8), int(fB>>8))
	} 
	if bA != 0 {
		c = c.AddBgRGB(int(bR>>8), int(bG>>8), int(bB>>8))
	}
	return c
}
