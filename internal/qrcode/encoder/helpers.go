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
