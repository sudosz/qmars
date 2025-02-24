package encoder

import (
	"image/color"

	gqrcode "github.com/skip2/go-qrcode"
	"github.com/sudosz/qmars/internal/qrcode"
)

type QRCodeBuilder struct {
	content       Content
	level         qrcode.RecoveryLevel
	version       qrcode.Version
	fgColor       color.Color
	bgColor       color.Color
	disableBorder bool
}

func NewQRCodeBuilder() *QRCodeBuilder {
	return &QRCodeBuilder{
		level:   qrcode.DefaultRecoveryLevel,
		version: qrcode.DefaultVersion,
		fgColor: color.Black,
		bgColor: color.White,
	}
}

func (b *QRCodeBuilder) SetContent(c Content) *QRCodeBuilder {
	b.content = c
	return b
}

func (b *QRCodeBuilder) SetRecoveryLevel(l qrcode.RecoveryLevel) *QRCodeBuilder {
	b.level = l
	return b
}

func (b *QRCodeBuilder) SetVersion(v qrcode.Version) *QRCodeBuilder {
	b.version = v
	return b
}

func (b *QRCodeBuilder) SetForegroundColor(fg color.Color) *QRCodeBuilder {
	b.fgColor = fg
	return b
}

func (b *QRCodeBuilder) SetBackgroundColor(bg color.Color) *QRCodeBuilder {
	b.bgColor = bg
	return b
}

func (b *QRCodeBuilder) SetDisableBorder(disableBorder bool) *QRCodeBuilder {
	b.disableBorder = disableBorder
	return b
}

func (b *QRCodeBuilder) Build() (_ qrcode.QRCode, err error) {
	data := b.content.Get()
	var qr *gqrcode.QRCode
	if b.version > 0 {
		qr, err = gqrcode.NewWithForcedVersion(data, int(b.version), gqrcode.RecoveryLevel(b.level))
	} else {
		qr, err = gqrcode.New(data, gqrcode.RecoveryLevel(b.level))
	}
	qr.ForegroundColor = b.fgColor
	qr.BackgroundColor = b.bgColor
	qr.DisableBorder = b.disableBorder

	return qr, err
}
