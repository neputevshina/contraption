package contraption

import (
	"regexp"

	"github.com/golang/freetype/truetype"
	"github.com/neputevshina/geom"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type bang = struct{}

// Font is a multipurpose wrapper for a loaded to memory TrueType font.
type Font struct {
	Data           []byte
	Parsed         *sfnt.Font
	FreeTypeParsed *truetype.Font
	// Drawable       map[int]PixelFace
	// OnGpu          map[rune]*gel.Outline
	Capk     fixed.Int26_6
	Name     string
	segcache map[rune][]Segment
	advcache map[rune][3]float64
}

func NewFont(data []byte, name string) (*Font, error) {
	f, err := sfnt.Parse(data)
	if err != nil {
		return nil, err
	}
	f2, err := truetype.Parse(data)
	return &Font{
		Data:           data,
		Parsed:         f,
		FreeTypeParsed: f2,
		// Drawable:       make(map[int]PixelFace),
		// OnGpu:          make(map[rune]*gel.Outline),
		Capk:     Capk(f),
		Name:     name,
		segcache: map[rune][]Segment{},
		advcache: map[rune][3]float64{},
	}, err
}

func (f *Font) CaptoemFixed(cap float64) fixed.Int26_6 {
	return f.Capk.Mul(fixed.Int26_6(cap * 64))
}

func (f *Font) Captoem(cap float64) float64 {
	return float64(f.CaptoemFixed(cap)) / 64
}

func (f *Font) EmtocapFixed(em fixed.Int26_6) float64 {
	return float64(em.Mul(f.Capk)) / 64
}

func (f *Font) Emtocap(em float64) float64 {
	return em * float64(f.Capk) / 64
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
	}
	return o
}

func (f *Font) Advance(r rune) float64 {
	o, ok := f.advcache[r]
	if !ok {
		g, _ := f.Parsed.GlyphIndex(&testingBuffer, r)
		bo, adv, err := f.Parsed.GlyphBounds(&testingBuffer, g, 12<<6, font.HintingNone)
		if err != nil {
			panic(err)
		}
		fadv := float64(adv) / (12 << 6)
		fmaxx := float64(bo.Max.X) / (12 << 6)
		fwidth := fmaxx - float64(bo.Min.X)/(12<<6)
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

// type PixelFace struct {
// 	font.Face
// 	// Cache values are preserved only between two consecutive frames.
// 	Rasterized map[rune]*struct {
// 		*gel.Texture
// 		Seen bool
// 	}
// }

// // Next is the new frame handler for PixelFace.
// func (pf PixelFace) Next() {
// 	for r, t := range pf.Rasterized {
// 		if t.Seen == false {
// 			delete(pf.Rasterized, r)
// 		}
// 	}
// }

// func (pf PixelFace) Rune(r rune) *gel.Texture {
// 	t, ok := pf.Rasterized[r]
// 	if !ok {
// 		// Almost copy-paste from unto/text.go
// 		dr, mask, mp, _, ok := pf.Face.Glyph(fixed.P(0, 0), r)
// 		if !ok {
// 			println("MISSING GLYPH", r)
// 		}
// 		w := int32(dr.Dx())
// 		h := int32(dr.Dy())
// 		// spaces and non-printable characters
// 		if w == 0 || h == 0 {
// 			return nil
// 		}
// 		su := image.NewRGBA(image.Rectangle{Max: dr.Size()})

// 		co := premultiply(hex(`#000000`))
// 		draw.DrawMask(su, su.Bounds(), &image.Uniform{co}, image.Point{}, mask, mp, draw.Over)
// 		te := gel.UploadUnfilteredTexture(su)
// 		t = &struct {
// 			*gel.Texture
// 			Seen bool
// 		}{Texture: te, Seen: true}
// 		pf.Rasterized[r] = t
// 	}
// 	t.Seen = true
// 	return t.Texture
// }

// func premultiply(straight unto.Color) (premultiplied color.RGBA) {
// 	premultiplied.R = (byte)(int(straight.R) * int(straight.A) / 255)
// 	premultiplied.G = (byte)(int(straight.G) * int(straight.A) / 255)
// 	premultiplied.B = (byte)(int(straight.B) * int(straight.A) / 255)
// 	premultiplied.A = straight.A
// 	return
// }

// // image/color.Color takes premultiplied color, so the image/draw.
// // formula from https://microsoft.github.io/Win2D/WinUI3/html/PremultipliedAlpha.htm
// var premulBlendMode = sdl.ComposeCustomBlendMode(
// 	sdl.BLENDFACTOR_ONE,
// 	sdl.BLENDFACTOR_ONE_MINUS_SRC_ALPHA,
// 	sdl.BLENDOPERATION_ADD,
// 	sdl.BLENDFACTOR_ZERO,
// 	sdl.BLENDFACTOR_ONE,
// 	sdl.BLENDOPERATION_ADD,
// )

var _ = func() {
}

var sizeRegexp = regexp.MustCompile(`[\pZ]*(\d+\.?\d*|\.\d+)(mm|pt|cm)[\pZ]*`)

func rectBox(r geom.Rectangle, thickness float64) (rs [4]geom.Rectangle) {
	for i := range rs {
		rs[i] = r
	}

	rs[0].Max.Y = rs[0].Min.Y
	rs[0].Min.Y -= thickness
	rs[0].Max.X += thickness

	rs[1].Min.X = rs[1].Max.X
	rs[1].Max.X += thickness
	rs[1].Max.Y += thickness

	rs[2].Min.Y = rs[2].Max.Y
	rs[2].Max.Y += thickness
	rs[2].Max.X += thickness

	rs[3].Max.X = rs[3].Min.X
	rs[3].Min.Y -= thickness
	rs[3].Min.X -= thickness

	return
}

func repeat[T any](count int, t T) (sl []T) {
	for i := 0; i < count; i++ {
		sl = append(sl, t)
	}
	return
}

type Sormer[T interface{ BaseWorld() *World }] interface {
	Sorm(T) Sorm
}
