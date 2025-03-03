package qrcode

import (
	"image"
	"image/color"
	"strings"

	"golang.org/x/image/draw"
)

type QRCode struct {
	bArr   [][]bool
	fg, bg color.Color
	invert bool
	margin int
	w, h   int
	BitMatrix
}

type BitMatrix interface {
	GetWidth() int
	GetHeight() int
	Get(x, y int) bool
}

func NewQRCode(b BitMatrix, invert bool, margin int, colors ...color.Color) *QRCode {
	qr := &QRCode{
		BitMatrix: b,
		w:         b.GetWidth(),
		h:         b.GetHeight(),
		margin:    margin,
		invert:    invert,
	}

	qr.fg = DefaultForeground
	qr.bg = DefaultBackground

	if len(colors) > 0 {
		qr.fg = colors[0]
	}

	if len(colors) > 1 {
		qr.bg = colors[1]
	}

	if invert {
		qr.fg, qr.bg = qr.bg, qr.fg
	}

	return qr
}

func (q *QRCode) SetForeground(fg color.Color) { q.fg = fg }
func (q *QRCode) SetBackground(bg color.Color) { q.bg = bg }
func (q *QRCode) SetInvert(i bool)             { q.invert = i }

func (q QRCode) GetForeground() color.Color { return q.fg }
func (q QRCode) GetBackground() color.Color { return q.bg }
func (q QRCode) GetWidth() int              { return q.w }
func (q QRCode) GetHeight() int             { return q.h }
func (q QRCode) GetMarginSize() int         { return q.margin }

func (q QRCode) ToBoolArray() [][]bool {
	q.bArr = make([][]bool, q.h)
	for i := 0; i < q.h; i++ {
		q.bArr[i] = make([]bool, q.w)
		for j := 0; j < q.w; j++ {
			q.bArr[i][j] = q.Get(j, i)
		}
	}
	return q.bArr
}

func (q QRCode) ToSmallString() string {
	var sb strings.Builder
	sb.Grow((q.w + 1) * (q.h / 2))

	for i := 0; i < q.h; i += 2 {
		for j := 0; j < q.w; j++ {
			writeColor(&sb, q.fg, q.bg)
			sb.WriteRune(smallChars[getCharOfBlockBools(q.Get(i, j), q.Get(i+1, j))])
			resetColor(&sb)
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (q QRCode) ToString(b StringBlock) string {
	var sb strings.Builder
	bw, bh := b.Bounds()
	l := 1 + q.w*bw
	sb.Grow(l * q.h * bh)

	for i := 0; i < q.h*bh; i++ {
		for j := 0; j < q.w*bw; j++ {
			if q.Get(j/bw, i/bh) {
				writeColor(&sb, q.fg, q.bg)
				sb.WriteRune(b.At(j%bw, i%bh))
				resetColor(&sb)
			} else {
				writeColor(&sb, q.bg, q.bg)
				sb.WriteRune(' ')
				resetColor(&sb)
			}
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (q QRCode) ToResizedImage(w, h int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dst, dst.Rect, q, q.Bounds(), draw.Over, nil)
	return dst
}

func (q QRCode) ColorModel() color.Model { return color.RGBAModel }

func (q QRCode) Bounds() image.Rectangle {
	return image.Rect(0, 0, q.GetWidth(), q.GetHeight())
}

func (q QRCode) At(x, y int) color.Color {
	if q.Get(x, y) {
		return q.fg
	}
	return q.bg
}

type customBlockQRCode struct {
	b    image.Image
	w, h int
	QRCode
}

func (q customBlockQRCode) ColorModel() color.Model { return color.RGBAModel }

func (q customBlockQRCode) Bounds() image.Rectangle {
	return image.Rect(0, 0, q.w, q.h)
}

func (q customBlockQRCode) At(x, y int) color.Color {
	c := q.bg
	realX, realY := x/q.b.Bounds().Dx(), y/q.b.Bounds().Dy()
	if q.Get(realX, realY) {
		c = q.b.At(x%q.b.Bounds().Dx(), y%q.b.Bounds().Dy())
	}
	if _, _, _, a := c.RGBA(); a == 0 {
		return q.bg
	}
	return c
}

func (q QRCode) ToImageWithBlock(block image.Image) image.Image {
	w, h := block.Bounds().Dx(), block.Bounds().Dy()
	return customBlockQRCode{b: block, QRCode: q, w: q.GetWidth() * w, h: q.GetHeight() * h}
}
