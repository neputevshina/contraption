package contraption

import (
	"image"
	"io"
	"math/rand"
	"strconv"

	"github.com/neputevshina/contraption/nanovgo"
	"github.com/neputevshina/geom"
)

func (s Sorm) paint(wo *World, f func()) {
	if s.fill != (nanovgo.Paint{}) || s.stroke != (nanovgo.Paint{}) {
		wo.Vgo.BeginPath()
	}
	wo.Vgo.ResetTransform()
	if s.fill != (nanovgo.Paint{}) {
		wo.Vgo.SetFillPaint(s.fill)
	}
	if s.stroke != (nanovgo.Paint{}) {
		wo.Vgo.SetStrokePaint(s.stroke)
	}
	f()
	wo.Vgo.ClosePath()
	if s.fill != (nanovgo.Paint{}) {
		wo.Vgo.Fill()
	}
	if s.strokew != 0 {
		wo.Vgo.SetStrokeWidth(s.strokew)
	}
	if s.stroke != (nanovgo.Paint{}) {
		if s.fill != (nanovgo.Paint{}) {
			wo.Vgo.Also()
		}
		wo.Vgo.Stroke()
	}
}

func (wo *World) NewText(font []byte) func(size float64, str []rune) Sorm {
	return wo.generalNewText(font, tagText)
}

func (wo *World) NewTopDownText(font []byte) func(size float64, str []rune) Sorm {
	return wo.generalNewText(font, tagTopDownText)
}

func (wo *World) NewBottomUpText(font []byte) func(size float64, str []rune) Sorm {
	return wo.generalNewText(font, tagBottomUpText)
}

func (wo *World) generalNewText(font []byte, kind tagkind) func(size float64, str []rune) Sorm {
	name := strconv.FormatUint(rand.Uint64(), 36)
	f, err := NewFont(wo.Vgo, font, name)
	if err != nil {
		panic(err)
	}
	if wo.capmap == nil {
		wo.capmap = map[int]float64{}
	}
	wo.nvgofontids = append(wo.nvgofontids, f.Vgoid)

	return func(size float64, str []rune) Sorm {
		s := wo.beginsorm()
		s.tag = kind

		s.Size.Y = size

		s.vecfont = f
		s.fontid = f.Vgoid
		// wo.Vgo.SetFontFaceID(s.fontid)
		// wo.Vgo.SetFontSize(size)
		// _, abcd := wo.Vgo.TextBounds(0, 0, str)
		// _, space := wo.Vgo.TextBounds(0, 0, " ")
		// s.Size.X = float64(abcd[2]-abcd[0]) - float64(space[2]-space[0])
		// if s.Size.X < 0 {
		// 	s.Size.X = 0
		// }
		if kind == tagTopDownText || kind == tagBottomUpText {
			s.Size.X, s.Size.Y = s.Size.Y, s.Size.X
		}
		s.key = str
		return s
	}
}
func generaltextrun(kind tagkind) func(wo *World, s *Sorm) {
	return func(wo *World, s *Sorm) {
		// TODO use io.RuneReader
		if s.fill == (nanovgo.Paint{}) {
			return
		}

		horizontal := kind == tagText

		wo.Vgo.ResetTransform()
		if horizontal {
			// Adjust baseline so 0y0 is top left.
			wo.Vgo.SetTransform(nanovgo.TranslateMatrix(float32(s.p.X), float32(s.p.Y+s.Size.Y)))
			wo.Vgo.SetFontSize(s.Size.Y)
		} else {
			wo.Vgo.SetTransform(nanovgo.TranslateMatrix(float32(s.p.X), float32(s.p.Y)))
			if kind == tagTopDownText {
				wo.Vgo.SetTransform(nanovgo.RotateMatrix(nanovgo.PI / 2))
			} else if kind == tagBottomUpText {
				wo.Vgo.SetTransform(nanovgo.RotateMatrix(-nanovgo.PI / 2))
			}
			wo.Vgo.SetFontSize(s.Size.X)
		}
		wo.Vgo.SetFontFaceID(s.fontid)
		wo.Vgo.SetFillPaint(s.fill)
		// println(s.key.(string), wo.vgo.CurrentTransform(), s.pos.X, s.pos.Y, s.Size.X, s.Size.Y)
		wo.Vgo.TextRune(0, 0, s.key.([]rune))
	}
}

// TODO This is the main and preferred method to do vector text.
// TODO Pool of (Runes)
func (wo *World) NewVectorTextReader(font []byte) func(size float64, rd io.RuneScanner) Sorm {
	name := strconv.FormatUint(rand.Uint64(), 36)
	id, err := NewFont(wo.Vgo, font, name)
	if err != nil {
		panic(err)
	}
	return func(size float64, rd io.RuneScanner) Sorm {
		s := wo.beginsorm()
		s.tag = tagVectorText
		s.Size.Y = size
		s.vecfont = id
		s.Size.X = s.vecfont.MeasureReader(size, rd)
		s.r = 0
		s.key = rd
		return s
	}
}

func (wo *World) NewVectorText(font []byte) func(size float64, str []rune) Sorm {
	name := strconv.FormatUint(rand.Uint64(), 36)
	id, err := NewFont(wo.Vgo, font, name)
	if err != nil {
		panic(err)
	}
	return func(size float64, str []rune) Sorm {
		s := wo.beginsorm()
		s.tag = tagVectorText
		s.Size.Y = size
		s.vecfont = id
		rr := Runes(str)
		s.Size.X = s.vecfont.MeasureReader(size, rr)
		s.r = 0
		s.key = rr
		return s
	}
}
func vectortextrun(wo *World, s *Sorm) {
	wo.Vgo.ResetTransform()
	fail := false
	if s.fill != (nanovgo.Paint{}) {
		wo.Vgo.SetFillPaint(s.fill)
	} else {
		fail = true
	}
	if s.stroke != (nanovgo.Paint{}) {
		wo.Vgo.SetStrokePaint(s.stroke)
	} else if fail {
		return
	}
	if s.fill != (nanovgo.Paint{}) || s.stroke != (nanovgo.Paint{}) {
		wo.Vgo.BeginPath()
	}
	wo.Vgo.SetTransform(nanovgo.TranslateMatrix(float32(s.p.X), float32(s.p.Y+s.Size.Y))) // Top left

	MakealineReader(wo.Vgo, s.vecfont, s.Size.Y, s.key.(io.RuneScanner))

	if s.fill != (nanovgo.Paint{}) {
		wo.Vgo.Fill()
	}
	if s.stroke != (nanovgo.Paint{}) {
		if s.fill != (nanovgo.Paint{}) {
			wo.Vgo.Also()
		}
		wo.Vgo.Stroke()
	}
	// if wo.Events.Match(`Press(F4)`) {
	// 	wo.vgo.DebugDumpPathCache()
	// }

}

// Circle is a circle shape.
func (wo *World) Circle(d float64) (s Sorm) {
	s = wo.beginsorm()
	s.tag = tagCircle
	s.Size.X = d
	s.Size.Y = d
	return
}
func circlerun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		r := s.Size.X / 2
		wo.Vgo.Circle(s.p.X+r, s.p.Y+r, r)
	})
}

// Rectangle is a rectangle shape.
func (wo *World) Rectangle(w, h complex128) (s Sorm) {
	s = wo.beginsorm()
	s.tag = tagRect
	s.Size.X = real(w)
	s.Size.Y = real(h)
	s.add.X = imag(w)
	s.add.Y = imag(h)
	wo.endsorm(s)
	return
}
func rectrun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		wo.Vgo.Rect(s.p.X, s.p.Y, s.Size.X, s.Size.Y)
	})
}

// Roundrect is a rounded rectangle shape.
func (wo *World) Roundrect(w, h complex128, r float64) (s Sorm) {
	s = wo.beginsorm()
	s.tag = tagRoundrect
	s.Size.X = real(w)
	s.Size.Y = real(h)
	s.add.X = imag(w)
	s.add.Y = imag(h)
	s.r = r
	wo.endsorm(s)
	return
}
func roundrectrun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		wo.Vgo.RoundedRect(s.p.X, s.p.Y, s.Size.X, s.Size.Y, s.r)
	})
}

func (wo *World) Void(w, h complex128) (s Sorm) {
	s = wo.beginsorm()
	s.tag = tagVoid
	s.Size.X = real(w)
	s.Size.Y = real(h)
	s.add.X = imag(w)
	s.add.Y = imag(h)
	wo.endsorm(s)
	return
}
func voidrun(wo *World, s *Sorm) {}

func (wo *World) Equation(eqn Equation) (s Sorm) {
	s = wo.beginsorm()
	s.tag = tagEquation

	pt := eqn.Size()
	s.Size.X, s.Size.Y = pt.X, pt.Y
	s.key = eqn

	if wo.eqnCache == nil {
		wo.eqnCache = map[any][]geom.Point{}
	}
	if wo.eqnLife == nil {
		wo.eqnLife = map[Equation]int{}
	}
	if wo.eqnWh == nil {
		wo.eqnWh = map[Equation]geom.Point{}
	}
	_, ok := wo.eqnCache[eqn]
	if !ok {
		wo.eqnCache[eqn] = impMarch(nil, eqn.Eqn, s.Size.X, s.Size.Y)
	}
	wo.eqnLife[eqn] = 2
	wo.endsorm(s)
	return
}
func equationrun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		a := wo.eqnCache[s.key]
		wo.Vgo.ResetTransform()
		wo.Vgo.SetTransform(nanovgo.TranslateMatrix(float32(s.p.X), float32(s.p.Y)))
		wo.Vgo.MoveTo(a[0].X, a[0].Y)
		for i := range a {
			if i == 0 {
				continue
			}
			wo.Vgo.LineTo(a[i].X, a[i].Y)
		}
	})
}

// Canvas gives a direct access to Nanovgo for painting a vector image.
func (wo *World) Canvas(w, h complex128, run func(vgo *Context, wt geom.Geom, rect geom.Rectangle)) (s Sorm) {
	s = wo.beginsorm()
	s.tag = tagCanvas
	s.canvas = run
	s.Size.X = real(w)
	s.Size.Y = real(h)
	s.add.X = imag(w)
	s.add.Y = imag(h)
	wo.endsorm(s)
	return
}
func canvasrun(wo *World, s *Sorm) {
	vgo := wo.Vgo

	vgo.Save()
	vgo.SetTransform(geom2nanovgo(s.m.Translate(s.p.X, s.p.Y)))
	if s.fill != (nanovgo.Paint{}) {
		vgo.SetFillPaint(s.fill)
	}
	if s.stroke != (nanovgo.Paint{}) {
		vgo.SetStrokePaint(s.stroke)
	}
	vgo.SetStrokeWidth(s.strokew)
	s.canvas(vgo, s.m, geom.Rect(s.p.X, s.p.Y, s.p.X+s.Size.X, s.p.Y+s.Size.Y))
	vgo.Restore()
}

// Sequence transforms external data to stream of shapes.
//
// Modifiers in Sequence right now are ignored, but can trigger a panic in future versions.
func (wo *World) Sequence(q Sequence) (s Sorm) {
	s = wo.beginsorm()
	s.tag = tagSequence
	s.key = q
	wo.endsorm(s)
	return s
}
func sequencerun(wo *World, s *Sorm) {
	// Special form, see (*Sorm).kidsiter and (*World).Develop
}

// Illustration is a static, unchangeable image.
// Src is decoded by image.Decoder.
// You must register codecs for every type of image you want to use with this shape.
//
// If w or h are zero, the corresponding axis of the resulting shape will be sized as the original picture.
//
// Negative sizes will result in a stretched image, whose size is fully controlled by the parent compound.
//
// mode is determining how the image will be resized if the proportions of the parent compound and the proportions of
// the image are different.
//
// Available modes are
//   - "stretch": the image will be stretched without preserving proportions.
//   - "zoom": the smallest dimension of the image will be resized to the largest available dimension.
//   - "pad": the smallest dimension of the image will be resized to the smallest available dimension.
//
// The image won't be cropped in any of those cases. Use Scissor in the parent compound to limit the size of an image.
//
// The texture from the image may be interpolated to a smaller size.
// Deallocation of the texture is the subject to the two-frame policy like any other resource in Contraption.
func (wo *World) Illustration(w, h complex128, mode string, src io.Reader) (s Sorm) {
	s = wo.beginsorm()
	s.tag = tagIllustration
	s.key = src
	u, ok := wo.images[src]
	if !ok {
		img, _, err := image.Decode(src)
		if err != nil {
			panic(err)
		}
		u.texid = wo.Vgo.CreateImageFromGoImage(0, img)
		u.origsiz = geom.PromotePt(img.Bounds().Size())
		wo.images[src] = u
	}

	s.Size = geom.Pt(real(w), real(h))
	s.add = geom.Pt(imag(w), imag(h))
	if w == 0 {
		s.Size.X = u.origsiz.X
	}
	if h == 0 {
		s.Size.Y = u.origsiz.Y
	}

	u.gen = wo.gen

	switch mode {
	case "stretch":
		s.fontid = 1
	case "zoom":
		s.fontid = 2
	case "pad":
		s.fontid = 3
	default:
		panic(`contraption: picture stretch mode is not one of "stretch", "zoom" or "pad"`)
	}

	wo.endsorm(s)
	return s
}
func illustrationrun(wo *World, s *Sorm) {
	vgo := wo.Vgo

	vgo.Save()

	vgo.SetTransform(geom2nanovgo(geom.Translate2d(s.p.X, s.p.Y)))

	u, ok := wo.images[s.key.(io.Reader)]
	if !ok {
		panic(`unreachable`)
	}

	frame := s.Size
	oz := u.origsiz
	origprop := oz.X / oz.Y

	// TODO Can be more pretty
	switch s.fontid {
	case 1:
		// frame == s.Size
	case 2:
		if origprop > frame.X/frame.Y {
			frame.X = frame.Y * origprop
		} else {
			frame.Y = frame.X / origprop
		}
	case 3:
		if origprop < frame.X/frame.Y {
			frame.X = frame.Y * origprop
		} else {
			frame.Y = frame.X / origprop
		}
	default:
		panic(`unreachable`)
	}

	o := s.Size.Sub(frame)
	o.X *= s.ialign.X
	o.Y *= s.ialign.Y

	vgo.BeginPath()
	vgo.SetFillPaint(nanovgo.ImagePattern(float32(o.X), float32(o.Y), float32(frame.X), float32(frame.Y), 0, u.texid, 1))
	vgo.Rect(o.X, o.Y, frame.X, frame.Y)
	vgo.Fill()

	vgo.Restore()
}

// Framebuffer is a raw image, possibly frequently updated.
//
// Every sizing rule of Illustration applies to Framebuffer.
// func (wo *World) Framebuffer(fb *Framebuffer, w, h complex128, mode string) (s Sorm) {
// 	s = wo.newSorm()
// 	s.tag = tagFramebuffer
// 	s.key = &src[0]
// 	return s
// }
