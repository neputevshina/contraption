package contraption

import (
	"github.com/golang/freetype/truetype"
	"github.com/neputevshina/geom"
	"github.com/neputevshina/nanovgo"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// Font is a multipurpose wrapper for a loaded to memory TrueType font.
type Font struct {
	Data           []byte
	Parsed         *sfnt.Font
	FreeTypeParsed *truetype.Font
	Vgoid          int
	Name           string

	vgok      float64
	capmap0   int
	capmap    []int
	segcache  map[rune][]Segment
	advcache  map[rune][3]float64
	minxcache map[rune]float64
	buf       sfnt.Buffer
}

func NewFont(vgo *nanovgo.Context, data []byte, name string) (*Font, error) {
	f, err := sfnt.Parse(data)
	if err != nil {
		return nil, err
	}
	f2, err := truetype.Parse(data)
	f3 := &Font{
		Data:           data,
		Parsed:         f,
		FreeTypeParsed: f2,
		Name:           name,
		segcache:       map[rune][]Segment{},
		advcache:       map[rune][3]float64{},
		minxcache:      map[rune]float64{},
	}
	if vgo != nil {
		f3.Vgoid = vgo.CreateFontFromMemory(name, data, 0)
	}
	k := 1.0
	if vgo != nil {
		k = f3.getvgok(vgo)
	}
	f3.vgok = k
	return f3, err
}

func (f *Font) getvgok(vgo *nanovgo.Context) float64 {
	return float64(f.FreeTypeParsed.VMetric(72<<6, f.FreeTypeParsed.Index('H')).TopSideBearing) / (72 << 6)
}

func (f *Font) CaptoemFixed(cap fixed.Int26_6) fixed.Int26_6 {
	// TODO Memoize?
	em0, em1 := 0.0, 72.0
	em := 0.0
	fem := func() fixed.Int26_6 { return fixed.Int26_6(em * 64) }

	for i := 0; i <= 30; i++ {
		em = (em0 + em1) / 2
		r, err := f.Parsed.Metrics(&f.buf, fem(), font.HintingNone)
		if err != nil {
			panic(err)
		}
		if r.CapHeight < cap {
			em0 = em
		} else {
			em1 = em
		}
	}
	return fem()
}

func (f *Font) Captoem(cap float64) float64 {
	return float64(f.CaptoemFixed(fixed.Int26_6(cap*64))) / 64
}

func (f *Font) EmtocapFixed(em fixed.Int26_6) float64 {
	r, err := f.Parsed.Metrics(&f.buf, em, font.HintingNone)
	if err != nil {
		panic(err)
	}
	return float64(r.CapHeight) / 64
}

func (f *Font) Emtocap(em float64) float64 {
	return f.EmtocapFixed(fixed.Int26_6(em * 64))
}

// func (f *Font) Rune(size float64, r rune) *gel.Texture {
// 	s := int(size)
// 	if s <= 6 { // Some fonts are not defined for less than 6 px
// 		return nil
// 	}
// 	if f.Drawable[s].Face == nil {
// 		f.Drawable[s] = PixelFace{
// 			Face: truetype.NewFace(f.FreeTypeParsed, &truetype.Options{
// 				Size:    size,
// 				DPI:     72,
// 				Hinting: font.HintingFull,
// 			}),
// 			Rasterized: make(map[rune]*struct {
// 				*gel.Texture
// 				Seen bool
// 			}),
// 		}
// 	}
// 	return f.Drawable[s].Rune(r)
// }

// func (f *Font) Outline(r rune) *gel.Outline {
// 	o, ok := f.OnGpu[r]
// 	if !ok {
// 		f.OnGpu[r] = gel.UploadOutline(f.Segments(r))
// 	}
// 	return o
// }

func gelSegmentFromSfnt(op int, args [3]fixed.Point26_6) Segment {
	pair := func(i fixed.Point26_6) (f geom.Point) {
		return geom.Pt(float64(i.X)/64.0, float64(i.Y)/64.0).Mul(1.0 / 12)
	}
	return Segment{
		Op: [...]byte{'M', 'L', 'Q', 'C'}[op], // Must be in sync with golang.org/x/image/font/sfnt.SegmentOp

		Args: [3]geom.Point{pair(args[0]), pair(args[1]), pair(args[2])},
	}
}

func (f *Font) Segments(r rune) []Segment {
	o, ok := f.segcache[r]
	if !ok {
		g, _ := f.Parsed.GlyphIndex(&testingBuffer, r)
		oldsegs, err := f.Parsed.LoadGlyph(&testingBuffer, g, 12<<6, nil)
		if err != nil {
			panic(err)
		}
		f.segcache[r] = collect(oldsegs, func(seg sfnt.Segment) Segment { return gelSegmentFromSfnt(int(seg.Op), seg.Args) })
		o = f.segcache[r]
		if len(o) > 0 {
			f.minxcache[r] = o[0].LastComponent().X
			for _, seg := range o[1:] {
				f.minxcache[r] = min(f.minxcache[r], seg.LastComponent().X)
			}
		}
	}
	return o
}

func (f *Font) Advance(r rune) float64 {
	o, ok := f.advcache[r]
	if !ok {
		g, _ := f.Parsed.GlyphIndex(&testingBuffer, r)
		bo, adv, err := f.Parsed.GlyphBounds(&testingBuffer, g, 120<<6, font.HintingNone)
		if err != nil {
			panic(err)
		}
		fadv := float64(adv) / (120 << 6)
		fmaxx := float64(bo.Max.X) / (120 << 6)
		fwidth := fmaxx - float64(bo.Min.X)/(120<<6)
		f.advcache[r] = [...]float64{fadv, fwidth, fmaxx}
		o = f.advcache[r]
	}
	return o[0]
}

func (f *Font) Width(r rune) float64 {
	o, ok := f.advcache[r]
	if !ok {
		_ = f.Advance(r)
		o = f.advcache[r]
	}
	return o[1]
}

func (f *Font) PureAdvance(r rune) float64 {
	o, ok := f.advcache[r]
	if !ok {
		_ = f.Advance(r)
		o = f.advcache[r]
	}
	return o[2]
}

func (f *Font) TrueXBearing(r rune) float64 {
	// _ = f.Segments(r)
	// return max(0, -f.minxcache[r])
	return f.Width(r) - f.PureAdvance(r)
}

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
	return makealine(vg, font, size, runes, false, true)
}

func MakealineStrict(vg *nanovgo.Context, font *Font, size float64, runes []rune) float64 {
	return makealine(vg, font, size, runes, true, true)
}

func makealine(vg *nanovgo.Context, font *Font, cap float64, runes []rune, extendedBox bool, draw bool) float64 {
	em := font.Captoem(cap)
	if len(runes) == 0 {
		return 0
	}

	x := font.TrueXBearing(runes[0])
	if extendedBox {
		m := 1000.0
		for _, s := range font.Segments(runes[0]) {
			m = min(m, s.Args[0].X)
			m = min(m, s.Args[1].X)
			m = min(m, s.Args[2].X)
		}
		x += m
	}
	if draw {
		vg.Save()
		vg.Scale(float32(em), float32(em))
	}
	for i, r := range runes {
		if draw {
			vg.Save()
			vg.Translate(float32(x), 0)
			Replay(vg, font.Segments(r))
			vg.PathWinding(nanovgo.Hole)
			vg.Restore()
		}

		if i < len(runes)-1 {
			x += font.Advance(r)
		} else {
			x += font.Advance(r)
		}
	}

	if draw {
		vg.Restore()
	}
	return x * em
}

func (font *Font) Measure(cap float64, runes []rune) float64 {
	return makealine(nil, font, cap, runes, false, false)
}
