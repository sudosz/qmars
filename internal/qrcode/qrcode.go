package qrcode

import (
	"fmt"
	"image"
	icolor "image/color"
	"strings"

	"golang.org/x/image/draw"
)

type QRCode struct {
	bArr   [][]bool
	fg, bg icolor.Color
	invert bool
	w, h   int
	BitMatrix
}

type BitMatrix interface {
	GetWidth() int
	GetHeight() int
	Get(x, y int) bool
}

func NewQRCode(b BitMatrix, invert bool, colors ...icolor.Color) *QRCode {
	qr := &QRCode{
		BitMatrix: b,
		w:         b.GetWidth(),
		h:         b.GetHeight(),
		invert:    invert,
	}

	if len(colors) > 0 {
		qr.fg = colors[0]
	} else {
		qr.fg = DefaultForeground
	}

	if len(colors) > 1 {
		qr.bg = colors[1]
	} else {
		qr.bg = DefaultBackground
	}

	return qr
}

func (q *QRCode) SetForeground(fg icolor.Color) { q.fg = fg }
func (q *QRCode) SetBackground(bg icolor.Color) { q.bg = bg }
func (q *QRCode) SetWidth(w int)                { q.w = w }
func (q *QRCode) SetHeight(h int)               { q.h = h }
func (q *QRCode) SetInvert(i bool)              { q.invert = i }

func (q QRCode) GetForeground() icolor.Color { return q.fg }
func (q QRCode) GetBackground() icolor.Color { return q.bg }
func (q QRCode) GetWidth() int               { return q.w }
func (q QRCode) GetHeight() int              { return q.h }

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

	printer := createColorPrinter(q.fg, q.bg)

	fmt.Println(q.invert)
	for i := 0; i < q.h; i += 2 {
		printer.SetWriter(&sb)
		for j := 0; j < q.w; j++ {
			sb.WriteRune(smallChars[getCharOfBlockBools(q.invert, q.Get(i, j), q.Get(i+1, j))])
		}
		printer.UnsetWriter(&sb)
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (q QRCode) ToString(set, unset string) string {
	var sb strings.Builder
	l := 1 + q.w*max(len(set), len(unset))
	sb.Grow(l * q.h)

	printer := createColorPrinter(q.fg, q.bg)

	if q.invert {
		set, unset = unset, set
	}

	for i := 0; i < q.h; i++ {
		printer.SetWriter(&sb)
		for j := 0; j < q.w; j++ {
			if q.Get(j, i) {
				sb.WriteString(set)
			} else {
				sb.WriteString(unset)
			}
		}
		printer.UnsetWriter(&sb)
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (q QRCode) ToResizedImage(w, h int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dst, dst.Rect, q, q.Bounds(), draw.Over, nil)
	return dst
}

func (q QRCode) ColorModel() icolor.Model { return icolor.RGBAModel }

func (q QRCode) Bounds() image.Rectangle {
	return image.Rect(0, 0, q.GetWidth(), q.GetHeight())
}

func (q QRCode) At(x, y int) icolor.Color {
	c := q.bg
	if q.Get(x, y) {
		c = q.fg
	}
	return c
}
