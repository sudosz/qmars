package qrcode

import (
	icolor "image/color"
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
	foreground = "\x1b[38;2;"
	background = "\x1b[48;2;"
	reset      = "\x1b[0m"
	ender      = "m"
)

func writeColor(w io.StringWriter, prefix string, r, g, b uint32) {
	w.WriteString(prefix)
	w.WriteString(strconv.Itoa(int(r >> 8)))
	w.WriteString(";")
	w.WriteString(strconv.Itoa(int(g >> 8)))
	w.WriteString(";")
	w.WriteString(strconv.Itoa(int(b >> 8)))
	w.WriteString(ender)
}

func writeColoredBlock(w io.StringWriter, block string, fg, bg icolor.Color) {
	fR, fG, fB, fA := fg.RGBA()
	bR, bG, bB, bA := bg.RGBA()
	if fA != 0 {
		writeColor(w, foreground, fR, fG, fB)
	}
	if bA != 0 {
		writeColor(w, background, bR, bG, bB)
	}
	w.WriteString(block)
	w.WriteString(reset)
}
