package encoder

import (
	"errors"
	icolor "image/color"

	"github.com/makiuchi-d/gozxing"
	gqrcode "github.com/makiuchi-d/gozxing/qrcode"

	"github.com/sudosz/qmars/internal/qrcode"
)

type QRCodeBuilder struct {
	content       Content
	level         qrcode.ErrorCorrectionLevel
	version       qrcode.Version
	disableBorder bool
	invert bool
	width         int
	height        int
	fg, bg        icolor.Color
}

func NewQRCodeBuilder() *QRCodeBuilder {
	return &QRCodeBuilder{
		level:   qrcode.DefaultRecoveryLevel,
		version: qrcode.DefaultVersion,
		width:   qrcode.DefaultWidth,
		height:  qrcode.DefaultHeight,
		fg:      qrcode.DefaultForeground,
		bg:      qrcode.DefaultBackground,
	}
}

func (b *QRCodeBuilder) SetContent(c Content) *QRCodeBuilder {
	b.content = c
	return b
}

func (b *QRCodeBuilder) SetErrorCorrectionLevel(l qrcode.ErrorCorrectionLevel) *QRCodeBuilder {
	b.level = l
	return b
}

func (b *QRCodeBuilder) SetVersion(v qrcode.Version) *QRCodeBuilder {
	b.version = v
	return b
}

func (b *QRCodeBuilder) SetDisableBorder(disableBorder bool) *QRCodeBuilder {
	b.disableBorder = disableBorder
	return b
}

func (b *QRCodeBuilder) SetWidth(w int) *QRCodeBuilder {
	b.width = w
	return b
}

func (b *QRCodeBuilder) SetHeight(h int) *QRCodeBuilder {
	b.height = h
	return b
}

func (b *QRCodeBuilder) SetForeground(fg icolor.Color) *QRCodeBuilder {
	b.fg = fg
	return b
}

func (b *QRCodeBuilder) SetBackground(bg icolor.Color) *QRCodeBuilder {
	b.bg = bg
	return b
}

func (b *QRCodeBuilder) SetInvert(i bool) *QRCodeBuilder {
	b.invert = i
	return b
}

func (b *QRCodeBuilder) Build() (qrcode.QRCode, error) {
	if b.content == nil {
		return nil, errors.New("content is empty")
	}

	data := b.content.Get()

	hints := map[gozxing.EncodeHintType]interface{}{
		gozxing.EncodeHintType_ERROR_CORRECTION: ecLevelToGozxingECLevel(b.level),
	}
	if b.disableBorder {
		hints[gozxing.EncodeHintType_MARGIN] = 0
	}
	if b.version > 0 {
		hints[gozxing.EncodeHintType_QR_VERSION] = b.version
	}

	bit, err := gqrcode.NewQRCodeWriter().Encode(data, gozxing.BarcodeFormat_QR_CODE, b.width, b.height, hints)
	if err != nil {
		return nil, err
	}
	
	return qrcode.NewQRCode(bit, b.invert, b.fg, b.bg), nil
}
