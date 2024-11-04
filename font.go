package contraption

import (
	"io"
	"math"
	"unicode/utf8"

	"github.com/golang/freetype/truetype"
	"github.com/neputevshina/geom"
	"github.com/neputevshina/nanovgo"
	"golang.org/x/exp/constraints"
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

	capmap0    int
	capmap     []int
	segcache   map[rune][]Segment
	advcache   map[rune][3]float64
	minxcache  map[rune]float64
	readmem    []rune
	readmemptr io.RuneReader
	buf        sfnt.Buffer
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
	return f3, err
}

func (f *Font) CaptoemFixed(cap fixed.Int26_6) fixed.Int26_6 {
	const topthresh = 72.0

	// TODO Memoize?
	em0, em1 := 0.0, topthresh
	em := 0.0
	fem := func() fixed.Int26_6 { return fixed.Int26_6(em * 64) }

	// Use plain coefficient to get cap if on this em glyphs certainly can't be hinted.
	if cap > f.EmtocapFixed(topthresh*64) {
		em = (float64(cap) / 64) / (f.Emtocap(topthresh) / topthresh)
		return fem()
	}

	for i := 0; i <= 15; i++ {
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
	return fixedToFloat(f.CaptoemFixed)(cap)
}

func (f *Font) EmtocapFixed(em fixed.Int26_6) fixed.Int26_6 {
	r, err := f.Parsed.Metrics(&f.buf, em, font.HintingNone)
	if err != nil {
		panic(err)
	}
	return r.CapHeight
}

func (f *Font) Emtocap(em float64) float64 {
	return fixedToFloat(f.EmtocapFixed)(em)
}

func fixedToFloat(f func(fixed.Int26_6) fixed.Int26_6) func(float64) float64 {
	return func(x float64) float64 {
		return float64(f(fixed.Int26_6(x*64))) / 64
	}
}

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
	return makealinerd(vg, font, size, Runes(runes), false, true)
}

func MakealineReader(vg *nanovgo.Context, font *Font, size float64, rd io.RuneScanner) float64 {
	return makealinerd(vg, font, size, rd, false, true)
}

func MakealineStrict(vg *nanovgo.Context, font *Font, size float64, runes []rune) float64 {
	return makealinerd(vg, font, size, Runes(runes), true, true)
}

func (font *Font) Measure(cap float64, runes []rune) float64 {
	return makealinerd(nil, font, cap, Runes(runes), false, false)
}

func makealinerd(vg *nanovgo.Context, font *Font, cap float64, runes io.RuneScanner, extendedBox bool, draw bool) float64 {
	em := font.Captoem(cap)

	r, _, err := runes.ReadRune()
	c := 0
	if err == io.EOF {
		return 0
	} else if err != nil {
		panic(err)
	}

	x := font.TrueXBearing(r)
	if extendedBox {
		m := 1000.0
		for _, s := range font.Segments(r) {
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
	for err == nil {
		if draw {
			vg.Save()
			vg.Translate(float32(x), 0)
			Replay(vg, font.Segments(r))
			vg.PathWinding(nanovgo.Hole)
			vg.Restore()
		}
		x += font.Advance(r)
		r, _, err = runes.ReadRune()
		c++
	}
	if err != io.EOF {
		panic(err)
	}

	if draw {
		vg.Restore()
	}
	for ; c > 0; c-- {
		err := runes.UnreadRune()
		if err != nil {
			panic(err)
		}
	}

	return x * em
}

func (font *Font) MeasureReader(cap float64, rd io.RuneScanner) float64 {
	return makealinerd(nil, font, cap, rd, false, false)
}

type runeSlice struct {
	r []rune
	n int
}

func (s *runeSlice) ReadRune() (r rune, size int, err error) {
	if s.n == len(s.r) {
		return 0, 0, io.EOF
	}
	r = s.r[s.n]
	s.n++
	size = utf8.RuneLen(r)
	return
}

func (s *runeSlice) UnreadRune() (err error) {
	s.n--
	if s.n < 0 {
		panic(`underflow`)
	}
	return nil
}

func Runes(rs []rune) io.RuneScanner {
	return &runeSlice{r: rs}
}

type stringSlice struct {
	r string
	n int
}

func (s *stringSlice) ReadRune() (r rune, size int, err error) {
	if s.n == len(s.r) {
		return 0, 0, io.EOF
	}
	r, size = utf8.DecodeRuneInString(s.r)
	s.n += size
	return
}

func (s *stringSlice) UnreadRune() (err error) {
	_, size := utf8.DecodeLastRuneInString(s.r)
	s.n -= size
	if s.n < 0 {
		panic(`underflow`)
	}
	return nil
}

func String(s string) io.RuneScanner {
	return &stringSlice{r: s}
}

type byteSlice struct {
	r []byte
	n int
}

func (s *byteSlice) ReadRune() (r rune, size int, err error) {
	if s.n == len(s.r) {
		return 0, 0, io.EOF
	}
	r, size = utf8.DecodeRune(s.r)
	s.n += size
	return
}

func (s *byteSlice) UnreadRune() (err error) {
	_, size := utf8.DecodeLastRune(s.r)
	s.n -= size
	if s.n < 0 {
		panic(`underflow`)
	}
	return nil
}

func Bytes(bs []byte) io.RuneScanner {
	return &byteSlice{r: bs}
}

type intSlice[T constraints.Integer] struct {
	r, n, m T
}

func (s *intSlice[T]) ReadRune() (r rune, size int, err error) {
	size = 1
	n := s.n
	if s.r < 0 {
		if s.n > 1 {
			n--
		} else {
			r = '-'
			return
		}
	}
	m := s.r / exp(10, s.m-n)
	if s.m-n < 0 && s.r != 0 {
		return 0, 0, io.EOF
	}
	r = '0' + rune(max(m, -m)%10)
	s.n++
	return
}

func exp[T constraints.Integer](a, b T) (c T) {
	c = 1
	for b > T(0) {
		c *= a
		b--
	}
	return
}

func (s *intSlice[T]) UnreadRune() (err error) {
	s.n--
	if s.n < 0 {
		panic(`underflow`)
	}
	return nil
}

func Int[T constraints.Integer](i T) io.RuneScanner {
	return &intSlice[T]{r: i, m: T(math.Floor(math.Log10(float64(i))))}
}
