package encoder

import (
	"image"
	icolor "image/color"
	"strings"

	"github.com/fatih/color"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
	"github.com/sudosz/qmars/internal/qrcode"
)

var ecLevelMap = map[qrcode.ErrorCorrectionLevel]decoder.ErrorCorrectionLevel{
	qrcode.ErrorCorrectionLevelLow:      decoder.ErrorCorrectionLevel_L,
	qrcode.ErrorCorrectionLevelMedium:   decoder.ErrorCorrectionLevel_M,
	qrcode.ErrorCorrectionLevelQuartile: decoder.ErrorCorrectionLevel_Q,
	qrcode.ErrorCorrectionLevelHigh:     decoder.ErrorCorrectionLevel_H,
}

func ecLevelToGozxingECLevel(l qrcode.ErrorCorrectionLevel) decoder.ErrorCorrectionLevel {
	if ecLevel, exists := ecLevelMap[l]; exists {
		return ecLevel
	}
	return decoder.ErrorCorrectionLevel_L
}

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

func qrCode2SmallString(qr qrcode.QRCode, invert bool) string {
	arr := qr.ToBoolArray()
	h := len(arr)
	w := len(arr[0])
	var sb strings.Builder
	sb.Grow((w + 1) * ((h + 1) / 2))

	for i := 0; i < h-h%2; i += 2 {
		for j := 0; j < w; j++ {
			sb.WriteRune(smallChars[getCharOfBlockBools(invert, arr[i][j], arr[i+1][j])])
		}
		sb.WriteRune('\n')
	}
	if h%2 == 1 {
		for j := 0; j < w; j++ {
			sb.WriteRune(smallChars[getCharOfBlockBools(invert, arr[h-1][j], false)])
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func QRCode2SmallString(qr qrcode.QRCode) string {
	return qrCode2SmallString(qr, false)
}

func QRCode2SmallStringInvert(qr qrcode.QRCode) string {
	return qrCode2SmallString(qr, true)
}

func QRCode2ColoredSmallString(qr qrcode.QRCode, fg, bg icolor.Color) string {
	coloredString := QRCode2SmallString(qr)
	return colorizeString(coloredString, fg, bg)
}

func QRCode2ColoredSmallStringInvert(qr qrcode.QRCode, fg, bg icolor.Color) string {
	coloredString := QRCode2SmallStringInvert(qr)
	return colorizeString(coloredString, fg, bg)
}

func colorizeString(s string, fg, bg icolor.Color) string {
	fR, fG, fB, _ := fg.RGBA()
	bR, bG, bB, _ := bg.RGBA()
	col := color.RGB(int(fR>>8), int(fG>>8), int(fB>>8)).AddBgRGB(int(bR>>8), int(bG>>8), int(bB>>8))

	var b strings.Builder
	b.Grow(len(s) + strings.Count(s, "\n"))

	for _, line := range strings.Split(s, "\n") {
		col.Fprint(&b, line)
		b.WriteByte('\n')
	}

	return b.String()
}

func colorsEqual(c1, c2 icolor.Color) bool {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2
}

func coloredImage(img image.Image, lastFg, fg, bg icolor.Color) image.Image {
	bounds := img.Bounds()
	coloredImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if colorsEqual(img.At(x, y), lastFg) {
				coloredImg.Set(x, y, fg)
			} else {
				coloredImg.Set(x, y, bg)
			}
		}
	}

	return coloredImg
}
