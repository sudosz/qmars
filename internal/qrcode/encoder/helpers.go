package encoder

import (
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