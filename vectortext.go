package contraption

import (
	"github.com/neputevshina/nanovgo"
)

func Replay(vg *nanovgo.Context, segs []Segment) {
	for _, s := range segs {
		a, b, c := s.Args[0], s.Args[1], s.Args[2]
		ax, ay := float32(a.X), float32(a.Y)
		bx, by := float32(b.X), float32(b.Y)
		cx, cy := float32(c.X), float32(c.Y)
		switch s.Op {
		case 'M':
			vg.MoveTo(ax, ay)
		case 'L':
			vg.LineTo(ax, ay)
		case 'Q':
			vg.QuadTo(ax, ay, bx, by)
		case 'C':
			vg.BezierTo(ax, ay, bx, by, cx, cy)
		}
	}
}

func Makealine(vg *nanovgo.Context, font *Font, size float64, runes []rune) float64 {
	return makealine(vg, font, size, runes, false)
}

func MakealineStrict(vg *nanovgo.Context, font *Font, size float64, runes []rune) float64 {
	return makealine(vg, font, size, runes, true)
}

func makealine(vg *nanovgo.Context, font *Font, size float64, runes []rune, extendedBox bool) float64 {
	x := -font.Advance(runes[0]) + font.Width(runes[0])
	if extendedBox {
		// x = 0
		m := 1000.0
		for _, s := range font.Segments(runes[0]) {
			m = min(m, s.Args[0].X)
			m = min(m, s.Args[1].X)
			m = min(m, s.Args[2].X)
		}
		x += m
	}
	vg.Save()
	vg.Scale(float32(font.Emtocap(size)), float32(font.Emtocap(size)))
	for _, r := range runes {
		vg.Save()

		vg.Translate(float32(x), 0)
		Replay(vg, font.Segments(r))
		vg.PathWinding(nanovgo.Hole)
		x += font.Advance(r)

		vg.Restore()
	}
	vg.Restore()
	return x * font.Emtocap(size)
}

func (font *Font) Measure(size float64, runes []rune) float64 {
	size = font.Emtocap(size)
	x := 0.0
	for _, r := range runes {
		x += font.Advance(r) * size
	}
	corr := -font.Advance(runes[len(runes)-1]) + font.Width(runes[len(runes)-1])
	// corr := -font.Advance(runes[len(runes)-1]) + font.Width(runes[len(runes)-1]) -
	// 	(font.PureAdvance(runes[0]) - font.Width(runes[0]))
	return x + corr*size
}
