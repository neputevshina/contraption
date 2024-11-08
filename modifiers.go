package contraption

import (
	"github.com/neputevshina/geom"
	"github.com/neputevshina/nanovgo"
)

func (wo *World) Vfollow() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagVfollow
	return
}
func vfollowrun(wo *World, c, m *Sorm) {
	c.aligner = alignerVfollow
}

func (wo *World) Hfollow() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagHfollow
	return
}
func hfollowrun(wo *World, c, m *Sorm) {
	c.aligner = alignerHfollow
}

func (wo *World) Halign(amt float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagHalign
	s.Size.X = clamp(0, amt, 1)
	return
}
func halignrun(wo *World, c, m *Sorm) {
	x := 0.0
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		x = max(x, k.Size.X)
	})
	c.Size.X = max(c.Size.X, x)
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		k.x += (x - k.Size.X) * m.Size.X
		k.ialign.X = m.Size.X
	})
}

func (wo *World) Valign(amt float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagValign
	s.Size.X = clamp(0, amt, 1)
	return
}
func valignrun(wo *World, c, m *Sorm) {
	y := 0.0
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		y = max(y, k.Size.Y)
	})
	c.Size.Y = max(c.Size.Y, y)
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		k.y += (y - k.Size.Y) * m.Size.X
		k.ialign.Y = m.Size.X
	})
}

func (wo *World) Fill(p nanovgo.Paint) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagFill
	s.fill = p
	return
}
func fillrun(wo *World, s, m *Sorm) {
	s.kidsiter(wo, kiargs{}, func(k *Sorm) {
		if k.tag >= 0 {
			k.fill = m.fill
		}
	})
}

func (wo *World) Stroke(p nanovgo.Paint) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagStroke
	s.stroke = p
	return
}
func strokerun(wo *World, s, m *Sorm) {
	s.kidsiter(wo, kiargs{}, func(k *Sorm) {
		if k.tag >= 0 {
			k.stroke = m.stroke
		}
	})
}

func (wo *World) Strokewidth(w float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagStrokewidth
	s.strokew = float32(w)
	return
}
func strokewidthrun(wo *World, s, m *Sorm) {
	s.flags |= flagSetStrokewidth
	s.kidsiter(wo, kiargs{}, func(k *Sorm) {
		if k.tag >= 0 {
			k.strokew = m.strokew
		}
	})
}

func (wo *World) Identity(key any) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagIdentity
	s.key = key
	return
}
func identityrun(wo *World, s, m *Sorm) {
	s.key = m.key
	m.key = nil
}

// Cond adds an event callback to a Compound.
func (wo *World) Cond(f func(m Matcher)) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagCond
	s.cond = f
	return
}
func condrun(wo *World, s, m *Sorm) {
	s.cond = m.cond
}

// CondFill adds an event callback to a Compound.
func (wo *World) CondFill(f func(geom.Rectangle) nanovgo.Paint) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagCondfill
	s.condfill = f
	return
}
func condfillrun(wo *World, s, m *Sorm) {
	s.condfill = m.condfill
}

func (wo *World) CondStroke(f func(geom.Rectangle) nanovgo.Paint) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagCondstroke
	s.condstroke = f
	return
}
func condstrokerun(wo *World, s, m *Sorm) {
	s.condstroke = m.condstroke
}

// Between adds a Sorm from given constructor between every other shape in a compound.
func (wo *World) Between(f func() Sorm) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagBetween
	s.key = f
	return
}
func betweenrun(wo *World, s, m *Sorm) {
	// Between is the special case and handled in Compound()
}

// BetweenVoid adds a Void between every other shape of a compound.
func (wo *World) BetweenVoid(w, h complex128) (s Sorm) {
	return wo.Between(func() Sorm { return wo.Void(w, h) })
}

// Source marks area of current compound as a drag source.
// It uses compound's identity (set with Identity modifier) as a drag value.
func (wo *World) Source() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagSource
	return
}
func sourcerun(wo *World, s *Sorm, m *Sorm) {
	s.flags |= flagSource
}

// Sink marks area of current compound as a drag sink.
// When program receives Release(1) event with mouse cursor inside a sink,
// it calls given function with a drag value.
func (wo *World) Sink(f func(drop any)) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagSink
	s.sinkid = len(wo.sinks)
	wo.sinks = append(wo.sinks, f)
	return
}
func sinkrun(wo *World, s *Sorm, m *Sorm) {
	s.sinkid = m.sinkid
	m.sinkid = 0
}

// Hshrink is a modifier that makes negative horizontal values inside a compound without Hfollow or Vfollow to
// be set not to horizontal limit, but to maximum horizontal value of objects with known size.
// TODO Make description more readable.
func (wo *World) Hshrink() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagHshrink
	return
}
func hshrinkrun(wo *World, s *Sorm, m *Sorm) {
	s.flags |= flagHshrink
}

// Vshrink is a modifier that works exactly like Hshrink, but for vertical sizes.
func (wo *World) Vshrink() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagVshrink
	return
}
func vshrinkrun(wo *World, s *Sorm, m *Sorm) {
	s.flags |= flagVshrink
}

// Scissor limits the painting area of a Compound to its Limit.
func (wo *World) Scissor() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagScissor
	return
}
func scissorrun(wo *World, s *Sorm, m *Sorm) {
	s.flags |= flagScissor
}

// Limit limits the maximum compound size to specified limits.
// If a given size is negative, it limits the corresponding size of a compound by
// the rules of negative units for shapes.
// TODO Imaginary limits.
func (wo *World) Limit(w, h float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagLimit
	s.Size.X, s.Size.Y = w, h
	return
}
func limitrun(wo *World, s, m *Sorm) {
	p := s.m.ApplyPt(geom.Pt(m.Size.X, m.Size.Y))
	if m.Size.X > 0 {
		m.Size.X = p.X
		s.l.X = min(s.l.X, m.Size.X)
	} else if m.Size.X < 0 {
		s.eprops.X = -m.Size.X
	}
	if m.Size.Y > 0 {
		m.Size.Y = p.Y
		s.l.Y = min(s.l.Y, m.Size.Y)
	} else if m.Size.Y < 0 {
		s.eprops.Y = -m.Size.Y
	}
}

// Posttransform applies transformation that only affects objects visually.
// It doesn't affect object sizes for layout.
func (wo *World) Posttransform(x, y float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagPosttransform
	s.Size.X = x
	s.Size.Y = y
	return
}
func posttransformrun(wo *World, c, m *Sorm) {
	c.x += m.Size.X
	c.y += m.Size.Y
	// Because moves are inherited in a separate pass
}

// Transform applies transformation that affects objects sizes for layout.
func (wo *World) Transform(m geom.Geom) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagTransform
	s.m = m
	return
}
func transformrun(wo *World, c, m *Sorm) {
	c.m = c.m.Mul(m.m)
	// Matrices are cascaded in (*World).resolvepremods
}

type labelt struct {
	value   any
	counter int
}

func (wo *World) Whereis(s Sorm) Sorm {
	s.flags |= flagFindme
	return s
}

// Key is a temporary key-value storage for on-screen state.
// The value is deleted if it had been not accessed for two frames.
func (wo *World) Key(k any) (v *any) {
	mv, ok := wo.keys[k]
	if !ok {
		wo.keys[k] = &labelt{}
		mv = wo.keys[k]
	}
	mv.counter = 2
	v = &mv.value
	return
}

func (wo *World) DontDecimate() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagDontDecimate
	return
}
func dontdecimaterun(wo *World, c, m *Sorm) {
	c.flags |= flagDontDecimate
}

func (wo *World) Decimate() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagDecimate
	return
}
func decimaterun(wo *World, c, m *Sorm) {
	c.flags |= flagDecimate
}
