package encoder

import (
	"errors"
	"image"

	"github.com/makiuchi-d/gozxing"
	gqrcode "github.com/makiuchi-d/gozxing/qrcode"

	"github.com/sudosz/qmars/internal/qrcode"
)

type qrCode struct {
	bArr [][]bool
	*gozxing.BitMatrix
}

func newQRCode(bit *gozxing.BitMatrix) qrcode.QRCode {
	return qrCode{
		BitMatrix: bit,
	}
}

func (q qrCode) GetBitMatrix() *gozxing.BitMatrix {
	return q.BitMatrix
}

func (q qrCode) ToBoolArray() [][]bool {
	if q.bArr == nil {
		q.bArr = make([][]bool, q.GetHeight())
		for i := 0; i < q.GetHeight(); i++ {
			q.bArr[i] = make([]bool, q.GetWidth())
			for j := 0; j < q.GetWidth(); j++ {
				q.bArr[i][j] = q.Get(j, i)
			}
		}
	}
	return q.bArr
}

func (q qrCode) ToImage() image.Image {
	return q
}

type QRCodeBuilder struct {
	content       Content
	level         qrcode.ErrorCorrectionLevel
	version       qrcode.Version
	disableBorder bool
	width         int
	height        int
}

func NewQRCodeBuilder() *QRCodeBuilder {
	return &QRCodeBuilder{
		level:   qrcode.DefaultRecoveryLevel,
		version: qrcode.DefaultVersion,
		width:   qrcode.DefaultWidth,
		height:  qrcode.DefaultHeight,
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

	enc := gqrcode.NewQRCodeWriter()
	bit, err := enc.Encode(data, gozxing.BarcodeFormat_QR_CODE, b.width, b.height, hints)
	if err != nil {
		return nil, err
	}

	return newQRCode(bit), nil
}
