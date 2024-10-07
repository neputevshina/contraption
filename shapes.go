package contraption

import (
	"io"
	"math/rand"
	"strconv"

	"github.com/neputevshina/geom"
	"github.com/neputevshina/nanovgo"
)

func (wo *World) NewText(font []byte, capk float64) func(size float64, str string) Sorm {
	return wo.generalNewText(font, capk, tagText)
}

func (wo *World) NewTopDownText(font []byte, capk float64) func(size float64, str string) Sorm {
	return wo.generalNewText(font, capk, tagTopDownText)
}

func (wo *World) NewBottomUpText(font []byte, capk float64) func(size float64, str string) Sorm {
	return wo.generalNewText(font, capk, tagBottomUpText)
}

func (wo *World) generalNewText(font []byte, capk float64, kind tagkind) func(size float64, str string) Sorm {
	name := strconv.FormatUint(rand.Uint64(), 36)
	f, err := NewFont(wo.Vgo, font, name)
	if err != nil {
		panic(err)
	}
	if wo.capmap == nil {
		wo.capmap = map[int]float64{}
	}
	wo.nvgofontids = append(wo.nvgofontids, f.Vgoid)

	return func(size float64, str string) Sorm {
		s := wo.newSorm()
		s.tag = kind
		s.H = size
		s.r = f.vgok * 1.42
		s.vecfont = f
		s.fontid = f.Vgoid
		wo.Vgo.SetFontFaceID(s.fontid)
		wo.Vgo.SetFontSize(float32(s.vecfont.Captoem(size) * s.r))
		_, abcd := wo.Vgo.TextBounds(0, 0, str)
		_, space := wo.Vgo.TextBounds(0, 0, " ")
		s.W = float64(abcd[2]-abcd[0]) - float64(space[2]-space[0])
		if kind == tagTopDownText || kind == tagBottomUpText {
			s.W, s.H = s.H, s.W
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
			wo.Vgo.SetTransform(nanovgo.TranslateMatrix(float32(s.x), float32(s.y)+float32(s.H)))
			wo.Vgo.SetFontSize(float32(s.vecfont.Captoem(s.H) * s.r))
		} else {
			wo.Vgo.SetTransform(nanovgo.TranslateMatrix(float32(s.x), float32(s.y)))
			if kind == tagTopDownText {
				wo.Vgo.SetTransform(nanovgo.RotateMatrix(nanovgo.PI / 2))
			} else if kind == tagBottomUpText {
				wo.Vgo.SetTransform(nanovgo.RotateMatrix(-nanovgo.PI / 2))
			}
			wo.Vgo.SetFontSize(float32(s.vecfont.Captoem(s.W) * s.r))
		}
		wo.Vgo.SetFontFaceID(s.fontid)
		wo.Vgo.SetFillPaint(s.fill)
		// println(s.key.(string), wo.vgo.CurrentTransform(), s.x, s.y, s.W, s.H)
		wo.Vgo.Text(0, 0, s.key.(string))
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
		s := wo.newSorm()
		s.tag = tagVectorText
		s.H = size
		s.vecfont = id
		s.W = s.vecfont.MeasureReader(size, rd)
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
		s := wo.newSorm()
		s.tag = tagVectorText
		s.H = size
		s.vecfont = id
		rr := Runes(str)
		s.W = s.vecfont.MeasureReader(size, rr)
		s.r = 0
		s.key = rr
		return s
	}
}
func vectortextrun(wo *World, s *Sorm) {
	// TODO use io.RuneReader
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
	wo.Vgo.SetTransform(nanovgo.TranslateMatrix(float32(s.x), float32(s.y+s.H))) // Top left

	// Makealine(wo.Vgo, s.vecfont, s.H, s.key.([]rune))
	MakealineReader(wo.Vgo, s.vecfont, s.H, s.key.(io.RuneScanner))

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

func (wo *World) Circle(d float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagCircle
	s.W = d
	s.H = d
	return
}
func circlerun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		r := s.W / 2
		wo.Vgo.Circle(float32(s.x+r), float32(s.y+r), float32(r))
	})
}

func (wo *World) Rectangle(w, h complex128) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagRect
	s.W = real(w)
	s.H = real(h)
	s.addw = imag(w)
	s.addh = imag(h)
	return
}
func rectrun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		wo.Vgo.Rect(float32(s.x), float32(s.y), float32(s.W), float32(s.H))
	})
}

func (wo *World) Roundrect(w, h complex128, r float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagRoundrect
	s.W = real(w)
	s.H = real(h)
	s.addw = imag(w)
	s.addh = imag(h)
	s.r = r
	return
}
func roundrectrun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		wo.Vgo.RoundedRect(float32(s.x), float32(s.y), float32(s.W), float32(s.H), float32(s.r))
	})
}

func (wo *World) Void(w, h complex128) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagVoid
	s.W = real(w)
	s.H = real(h)
	s.addw = imag(w)
	s.addh = imag(h)
	return
}
func voidrun(wo *World, s *Sorm) {}

func (wo *World) Equation(eqn Equation) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagEquation

	pt := eqn.Size()
	s.W, s.H = pt.X, pt.Y
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
		wo.eqnCache[eqn] = impMarch(nil, eqn.Eqn, s.W, s.H)
	}
	wo.eqnLife[eqn] = 2
	return
}
func equationrun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		a := wo.eqnCache[s.key]
		wo.Vgo.ResetTransform()
		wo.Vgo.SetTransform(nanovgo.TranslateMatrix(float32(s.x), float32(s.y)))
		wo.Vgo.MoveTo(float32(a[0].X), float32(a[0].Y))
		for i := range a {
			if i == 0 {
				continue
			}
			wo.Vgo.LineTo(float32(a[i].X), float32(a[i].Y))
		}
	})
}

func (wo *World) Canvas(w, h complex128, run func(vgo *nanovgo.Context, rect geom.Rectangle)) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagCanvas
	s.canvas = run
	s.W = real(w)
	s.H = real(h)
	s.addw = imag(w)
	s.addh = imag(h)
	return
}
func (wo *World) Canvas2(w, h complex128, run func(vgo *nanovgo.Context, rect geom.Rectangle)) (s Sorm) {
	// FIXME Must be standard behavior of Canvas: scale by transform.
	s = wo.newSorm()
	s.tag = tagCanvas
	s.r = 1
	s.canvas = run
	s.W = real(w)
	s.H = real(h)
	s.addw = imag(w)
	s.addh = imag(h)
	return
}
func canvasrun(wo *World, s *Sorm) {
	vgo := wo.Vgo

	vgo.Save()
	vgo.Reset()
	if s.r > 0 {
		vgo.SetTransform(geom2nanovgo(s.m.Translate(s.x, s.y)))
	}
	if s.fill != (nanovgo.Paint{}) {
		vgo.SetFillPaint(s.fill)
	}
	if s.stroke != (nanovgo.Paint{}) {
		vgo.SetStrokePaint(s.stroke)
	}
	vgo.SetStrokeWidth(s.strokew)
	s.canvas(vgo, geom.Rect(s.x, s.y, s.x+s.W, s.y+s.H))
	vgo.Restore()
}

// Sequence transforms external data to stream of Sorms.
func (wo *World) Sequence(q Sequence, plus ...Sorm) (s Sorm) {
	// TODO Visibility check
	// FIXME Double allocation, implement sequencesequencealigner.

	tmp := wo.tmpalloc(q.Length() + len(plus))
	j := 0
	for i := range tmp[:q.Length()] {
		s := q.Get(i)
		// Skip wo.Void(0,0)
		// if s.tag == tagVoid && s.W == 0 && s.H == 0 {
		// 	continue
		// }
		// Skip Sorm{}
		if zero(s) {
			continue
		}
		tmp[j] = s
		j++
	}
	copy(tmp[j:], plus)
	return wo.compound(wo.newSorm(), tmp...)
}
