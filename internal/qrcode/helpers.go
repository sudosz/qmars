package qrcode

import (
	"image/color"
	"io"
	"strconv"
)

var (
	smallChars = [4]rune{' ', '▀', '▄', '█'}
)

func getCharOfBlockBools(bools ...bool) uint8 {
	var idx uint8
	for i, b := range bools {
		if b {
			idx |= 1 << i
		}
	}
	return idx
}

const (
	foreground = 38
	background = 48
	reset      = 0
)

func appendColor(b []byte, params ...int64) []byte {
	b = append(b, "\x1b["...)
	for i, param := range params {
		if i > 0 {
			b = append(b, ';')
		}
		b = strconv.AppendInt(b, param, 10)
	}
	return append(b, 'm')
}

func writeColor(w io.Writer, fg, bg color.Color) {
	fR, fG, fB, fA := fg.RGBA()
	bR, bG, bB, bA := bg.RGBA()
	if fA == 0 && bA == 0 {
		return
	}

	var b []byte
	if fA != 0 {
		b = appendColor(b, foreground, 2, int64(fR>>8), int64(fG>>8), int64(fB>>8))
	}
	if bA != 0 {
		b = appendColor(b, background, 2, int64(bR>>8), int64(bG>>8), int64(bB>>8))
	}

	w.Write(b)
}

func resetColor(w io.Writer) {
	w.Write(appendColor(nil, reset))
}