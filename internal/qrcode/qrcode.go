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

func (q *QRCode) SetForeground(fg color.Color) *QRCode {
	q.fg = fg
	return q
}

func (q *QRCode) SetBackground(bg color.Color) *QRCode {
	q.bg = bg
	return q
}

func (q *QRCode) SetInvert(i bool) *QRCode {
	q.invert = i
	return q
}

func (q QRCode) GetForeground() color.Color { return q.fg }
func (q QRCode) GetBackground() color.Color { return q.bg }
func (q QRCode) GetWidth() int              { return q.w }
func (q QRCode) GetHeight() int             { return q.h }
func (q QRCode) GetMarginSize() int         { return q.margin }

func (q QRCode) ToBoolArray() [][]bool {
	if q.bArr == nil {
		// Pre-allocate the entire 2D slice at once
		q.bArr = make([][]bool, q.h)
		data := make([]bool, q.w*q.h)
		for i := range q.bArr {
			q.bArr[i] = data[i*q.w : (i+1)*q.w]
			for j := 0; j < q.w; j++ {
				q.bArr[i][j] = q.Get(j, i)
			}
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
		sb.WriteByte('\n')
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
				sb.WriteByte(' ')
				resetColor(&sb)
			}
		}
		sb.WriteByte('\n')
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
	blockW, blockH := q.b.Bounds().Dx(), q.b.Bounds().Dy()
	realX, realY := x/blockW, y/blockH
	ms := q.GetMarginSize()

	// Early return for margin area
	if realX < ms || realY < ms || realX >= q.GetWidth()-ms || realY >= q.GetHeight()-ms {
		return q.bg
	}

	adjustedX, adjustedY := realX-ms, realY-ms
	if color, isPattern := q.checkInPatterns(q.GetWidth()-2*ms, q.GetHeight()-2*ms, adjustedX, adjustedY); isPattern {
		return color
	}

	if q.Get(realX, realY) {
		c := q.b.At(x%blockW, y%blockH)
		if _, _, _, a := c.RGBA(); a > 0 {
			return c
		}
	}

	return q.bg
}

func (q customBlockQRCode) checkInPatterns(w, h, x, y int) (color.Color, bool) {
	// Check timing patterns
	if (x == w-9 || x == w-5) && (h-10 < y && y < h-4) ||
		(y == h-9 || y == h-5) && (w-10 < x && x < w-4) ||
		(x == w-7 && y == h-7) {
		return q.GetForeground(), true
	}

	pt := image.Pt(x, y)
	if pt.In(image.Rect(w-7, h-7, w, h)) {
		return nil, false
	}

	// Check finder patterns
	if (x < 7 || x >= w-7) && (y < 7 || y >= h-7) {
		if x == 0 || x == 6 || x == w-7 || x == w-1 ||
			y == 0 || y == 6 || y == h-7 || y == h-1 {
			return q.GetForeground(), true
		}
	}

	// Check alignment patterns
	if pt.In(image.Rect(2, 2, 5, 5)) ||
		pt.In(image.Rect(w-5, 2, w-2, 5)) ||
		pt.In(image.Rect(2, h-5, 5, h-2)) {
		return q.GetForeground(), true
	}

	return nil, false
}

func (q QRCode) ToImageWithBlock(block image.Image) image.Image {
	w, h := block.Bounds().Dx(), block.Bounds().Dy()
	return customBlockQRCode{b: block, QRCode: q, w: q.GetWidth() * w, h: q.GetHeight() * h}
}
