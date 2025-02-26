package encoder

import (
	"image"
	icolor "image/color"

	"github.com/makiuchi-d/gozxing"
	"github.com/sudosz/qmars/internal/qrcode"
)

type qrCode struct {
	bArr   [][]bool
	fg, bg icolor.Color
	*gozxing.BitMatrix
}

func newQRCode(bit *gozxing.BitMatrix, fg, bg icolor.Color) qrcode.QRCode {
	return qrCode{
		BitMatrix: bit,
		fg:        fg,
		bg:        bg,
	}
}

func (q qrCode) Foreground() icolor.Color {
	return q.fg
}

func (q qrCode) Background() icolor.Color {
	return q.bg
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
	if q.fg == qrcode.DefaultForeground && q.bg == qrcode.DefaultBackground {
		return q
	}
	return coloredImage(q, qrcode.DefaultForeground, q.fg, q.bg)
}

func (q qrCode) ToColoredString(set, unset string) string {
	return colorizeString(q.ToString(set, unset), q.fg, q.bg)
}
