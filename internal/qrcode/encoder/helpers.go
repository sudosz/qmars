package encoder

import (
	"strings"

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
	smallChars = map[int]rune{
		0:     ' ',
		1:     '▀',
		2:     '▄',
		1 | 2: '█',
	}
	smallCharsInvert = map[int]rune{
		0:     '█',
		1:     '▄',
		2:     '▀',
		1 | 2: ' ',
	}
)

func getCharOfBlockBools(bools ...bool) int {
	idx := 0
	for i, b := range bools {
		if b {
			idx |= 1 << i
		}
	}
	return idx
}

func getQRCodeSmallestString(qr qrcode.QRCode, charset map[int]rune) string {
	b := &strings.Builder{}
	arr := qr.ToBoolArray()
	h := len(arr)
	hodd := h%2 == 1
	for i := 0; i < h-h%2; i += 2 {
		w := len(arr[i])
		for j := 0; j < w; j++ {
			b.WriteRune(charset[getCharOfBlockBools(arr[i][j], arr[i+1][j])])
		}
		b.WriteRune('\n')
	}
	if hodd {
		for j := 0; j < len(arr[h-1]); j++ {
			b.WriteRune(charset[getCharOfBlockBools(arr[h-1][j], false)])
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func GetQRCodeSmallestString(qr qrcode.QRCode) string {
	return getQRCodeSmallestString(qr, smallChars)
}

func GetQRCodeSmallestStringInvert(qr qrcode.QRCode) string {
	return getQRCodeSmallestString(qr, smallCharsInvert)
}
