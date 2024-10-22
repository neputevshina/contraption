// Contraption: a simple framework for user interfaces.
//
// A good user interface framework must be an engine for a word processing game.
//
// TODO:
//	- Round shape sizes inside aligners.
//		- Noround modifier
//		? Smooth rounding — don't round if in motion
//		- Probably when implementing this will be the best time to go for
//			geom.Rect for storing xywh
//	? LimitOverride
//	- Grid aligner
//		- wo.Hgrid(cols int) + wo.Halign() — secondary alignment
//		- wo.Vgrid(rows int) + wo.Valign()
//		- Primary alignment won't work — makes no sense
//		- Negative sizes in primary axis won't work.
//		- No instructions for items like in CSS grid.
//		- Negative sizes in secondary axis are distributed.
//	- BufferSequence
//		- MemoBufferSequence
//	- RuneScannerSequence
//		- Needs recording of the whole given RuneScanner, at least with current Sequence implementation.
//	± Imaginary sizes.
//		- -1 + 20i — negative stretch, add 20 scaled by local transform pixels to size on layout step
//		- -1 - 20i  <  -1  <  -1 + 20i
//		+ Use: extending hitbox of elements without changing layout (together with Override)
//		± Works only on noaligner, sequences ignore it yet.
//	- Progress reader for IO operations.
//		- contraption.ProgressReader{ rd io.Reader; bytes, byteswritten int }
//		- var _ io.Reader = &ProgressReader{}
//		- (*ProgressReader).Remaining() float64 -> 0.0–0.1
//		- (*ProgressReader).RemainingBytes() int
//	- Display sinks on F2 view. Make F2 configurable also.
//	- Stylus events
//		- Touch(<50), Touch(>= 10) — threshold for pressure
//	- File drag event: Drag(*.txt)
//		- A companion for Drop — matches when the file is dragged above the area.
//		- Needs changes in GLFW or changing input library. SDL supports this.
//	- Interactive views for very large 1d and 2d data: waveforms, giant Minecraft maps, y-log STFT frames, etc.
//		- Easy insertion of Sorms between the data. See https://www.youtube.com/watch?v=Cz0OvnR_aoY.
//			- Probably very easy to implement by simply slicing the data.
//		- Why? Try to display STFT of a music file using Matplotlib, then rescale the window. Enjoy the delay.
//	- Text area
//		- https://rxi.github.io/textbox_behaviour.html
//	+ Sequence must be a special shape that pastes Sorms inside a compound, not being compound itself
//		- So wo.Text(io.RuneReader) could be Sequence
//		+ Not clear how to reuse memory of pools in this case
// 	- Quadtree for pool
//	- Matching past in regexps and coords change [MAJOR TOPIC]
//		- Easily solved with hitmaps — just draw a hitmap with all component's transformations
//		- Could use per-event UV deltas, but 64x viewport memory overhead is too much
//		- Use VDOM — retain, reconcile and feedback
//		- Save matrix for every shape that looks behind, 64×8×16×[shape count] bytes of overhead
//			- Simpler version: for every shape that has Cond/CondPaint, which is larger but still less than 10% of a tree
//	- Laziness and scrolling [MAJOR TOPIC]
//		- wo.Sequence(seq Sequence) — a window to infinity!
//			- Every returned Sorm is included to the parent
//	- A way to create uniformly sized buttons (as per tonsky)
//		- Just create a special component for this, dumb ass
//		- func Huniform(...Sorm) Sorm
//		- func Vuniform(...Sorm) Sorm
//		- Can use new negative value behavior? Just make needed widths/heights equal negative values.
//		- Integer key to determine which sizes must be equal:
//			- func Eqkey() Eqkey
//			- func Hequal(Eqkey) Sorm
//			- func Vequal(Eqkey) Sorm
//		- Can be used to implement grid layout.
//		- Other proposed names: Hequalize, Vequalize
//		- H2Vfollow, V2Hfollow — stretch as one, lay out as another
//	- Subworlds — layout inside canvases
//	- Modifier to shape position independence
//	- Fix paint interface
//	- Animations
//	- Remove bodges from layout (impossible)
//	+ Drag'n'drop
//	- Vector boolean ops
//		- Intersecton()
//		- Union()
//		- Subtraction()
//		- Difference()
//	- Word layout
//		- Together with stretch creates a flexbox-like system
//		- Together with laziness creates a universal layout framework, capable of word processing
//		- Will be used for text.
//		- func Hwords(perline func() Sorm) Sorm
//		- func Vwords(perline func() Sorm) Sorm
//		- Secondary axis limits influenced by perpendicular Void
//		- wo.Text(io.RuneScanner) []Sorm
//			- Returns Knuth-Plass-ready stream of boxes, Glues and Penalties.
//			? How to insert anything in between symbols?
//				? RuneScanner splitter?
//		- wo.Cap(float64) (can't be negative)
//		- wo.Lsp(float64)
//		- Knuth-Plass
//			? Interpret negative sizes as glue.
// 			- func Hknuth(perline func() Sorm) Sorm
// 			- func Vknuth(perline func() Sorm) Sorm
// 			- func Glue(width, minus, plus float64) Sorm // Analogous to wo.Void() but undirectional.
// 			- func Penalty(replacewith func() Sorm, penalty float64) Sorm
//				- Alt: Penalty as a builder on a target shape
//			- Void(0, y) is already a “strut”
//	- Tiling
//		- Tiled rect
//		- Tiled path
//	- Localization and internationalization guideline
//	- Strict methodology of usage
//	- func Retain(Sorm) (Sorm, struct{})   // Second returned value is needed only to restrict user from pasting
//	- func Deploy(Sorm) (Sorm, struct{})  // it directly to Compound. Because it will break slices.
//	- Other backends: Gio, software
//		- https://rxi.github.io/cached_software_rendering.html
//	~ func Whereis(Sorm) Sorm — prints where object is on overlay for debug
//	- func Target(onScreen *bool) Sorm
//	- Commenting the interface
//	- Rotations
//	- Activator stack
//	~ Move -> Transform
//	+ Scale -> Pretransform
//	± Click and and get every line of code that tried to paint over that pixel.
//		- Now make it less suck.
//	+ Separate update and render loops and add func (wo *World) Simulate(until time.Time)
//	+ Alignment
//	+ Positioning
//	+ Spacing
//	+ Key prop
//	+ Conditional fill and stroke, event handling
//	+ Bounds override
//	+ Draw mark
//	+ Canvas
// 	+ Depth-first layout
//		+ Stretch
//		+ Draw order is not dependent on call order
//		+ Scissor
//	+ Autovoid container -> Between modifier
//	+ Fix scale modifier
//	+ wo.Max() (wo.Limit())
//
//
//	Anti-todo:
//	- Geometric constraints
//	- 3D
//	- Stylesheets
//

package contraption

import (
	"encoding/gob"
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/neputevshina/geom"
	"github.com/neputevshina/nanovgo"
	"golang.org/x/exp/slices"
)

const sequenceChunkSize = 10

type flagval uint

const (
	flagSource flagval = 1 << iota
	flagOvrx
	flagOvry
	flagBetweener
	flagScissor
	flagMark
	flagFindme
	flagHshrink
	flagVshrink
	flagSetStrokewidth
	flagSequenceMark
	flagSequenceSaved
	flagBreakIteration
	flagNegativeOvrx
	flagNegativeOvry
)

//go:generate stringer -type=tagkind -trimprefix=tag
type tagkind int

//go:generate stringer -type=alignerkind
type alignerkind int

const (
	tagCompound tagkind = 0
)
const (
	_ tagkind = iota
	tagCircle
	tagRect
	tagRoundrect
	tagVoid
	tagEquation
	tagText
	tagCanvas
	tagVectorText
	tagTopDownText
	tagBottomUpText
	tagSequence
)
const (
	_ tagkind = -iota
	tagHalign
	tagValign
	tagFill
	tagStroke
	tagStrokewidth
	tagIdentity
	tagCond
	tagCondfill
	tagCondstroke
	tagBetween
	tagScroll
	tagSource
	tagSink
)
const (
	_ tagkind = -100 - iota
	tagTransform
	tagPretransform
	tagScissor
	tagHshrink
	tagVshrink
	tagLimit // Limit must be executed after update of a matrix, but before aligners, because it influences size.
	tagVfollow
	tagHfollow
)
const (
	alignerNone alignerkind = iota
	alignerVfollow
	alignerHfollow
)

// NOTE This might actually be a single table, if needed.
var preActions [100]func(wo *World, compound, mod *Sorm)
var alignerActions [100]func(wo *World, compound *Sorm)
var shapeActions [100]func(wo *World, shape *Sorm)
var modActions [100]func(wo *World, compound, mod *Sorm)

// Before troubleshooting Sorm-returning methods, check this:
//   - tagX -> xrun
//   - (*World).X(...) sets tagX to Sorm
//   - `stringer` was called
func init() {
	shapeActions[tagCircle] = circlerun
	shapeActions[tagRect] = rectrun
	shapeActions[tagRoundrect] = roundrectrun
	shapeActions[tagVoid] = voidrun
	shapeActions[tagEquation] = equationrun
	shapeActions[tagCanvas] = canvasrun
	shapeActions[tagVectorText] = vectortextrun
	shapeActions[tagText] = generaltextrun(tagText)
	shapeActions[tagTopDownText] = generaltextrun(tagTopDownText)
	shapeActions[tagBottomUpText] = generaltextrun(tagBottomUpText)
	shapeActions[tagSequence] = sequencerun

	modActions[-tagHalign] = halignrun
	modActions[-tagValign] = valignrun
	modActions[-tagFill] = fillrun
	modActions[-tagStroke] = strokerun
	modActions[-tagStrokewidth] = strokewidthrun
	modActions[-tagIdentity] = identityrun
	modActions[-tagCond] = condrun
	modActions[-tagCondfill] = condfillrun
	modActions[-tagCondstroke] = condstrokerun
	modActions[-tagBetween] = betweenrun
	modActions[-tagSource] = sourcerun
	modActions[-tagSink] = sinkrun

	preActions[-100-tagTransform] = transformrun
	preActions[-100-tagPretransform] = pretransformrun
	preActions[-100-tagScissor] = scissorrun
	preActions[-100-tagVfollow] = vfollowrun
	preActions[-100-tagHfollow] = hfollowrun
	preActions[-100-tagHshrink] = hshrinkrun
	preActions[-100-tagVshrink] = vshrinkrun
	preActions[-100-tagLimit] = limitrun

	alignerActions[alignerNone] = noaligner
	alignerActions[alignerVfollow] = vfollowaligner
	alignerActions[alignerHfollow] = hfollowaligner
}

type Eqn func(pt geom.Point) (dist float64)

type Equation interface {
	Eqn(pt geom.Point) (dist float64)
	Size() geom.Point
}

type World struct {
	*Events

	Window       Window
	Oscilloscope Oscilloscope

	tmp []Sorm

	nextn   int
	pool    []Sorm
	auxn    int
	auxpool []Sorm
	prefix  int

	old    []Sorm
	auxold []Sorm

	Vgo *nanovgo.Context

	Wwin, Hwin float64

	eqnCache map[any][]geom.Point
	eqnLife  map[Equation]int
	eqnWh    map[Equation]geom.Point

	nvgofontids []int
	capmap      map[int]float64

	rend int

	runepool []rune

	keys      map[any]*labelt
	BeforeVgo func()

	drag       any
	dragstart  geom.Point
	sinks      []func(any)
	DragEffect func(interval [2]geom.Point, drag any) Sorm

	showOutlines bool
	f1           bool

	alloc func(n int) (left, right int)
}

type Sorm struct {
	z, z2, i             int
	tag                  tagkind
	flags                flagval
	W, H, wl, hl         float64
	addw, addh           float64
	r, x, y              float64
	m, prem              geom.Geom
	aligner              alignerkind
	known, props, eprops geom.Point
	kidsl, kidsr,
	modsl, modsr,
	presl, presr int

	index *Index

	// TODO Essentially, this is always equal to (0, 0, wl, hl) when unscaled.
	//	Makes sense to remove wl and hl and only use this.
	scissor geom.Rectangle

	fill    nanovgo.Paint
	stroke  nanovgo.Paint
	strokew float32

	fontid  int
	vecfont *Font

	// Some objects use key field for own purposes:
	// 	- Equation stores an Equation object
	// 	- Text stores a io.RuneReader
	// 	- Compound stores an Identity, which also works
	//	  out as Source's dropable object
	// 	- Between stores func() Sorm
	key any

	condfill       func(rect geom.Rectangle) nanovgo.Paint
	condstroke     func(rect geom.Rectangle) nanovgo.Paint
	condfillstroke func(rect geom.Rectangle) (nanovgo.Paint, nanovgo.Paint)
	cond           func(m Matcher)
	canvas         func(vgo *nanovgo.Context, wt geom.Geom, rect geom.Rectangle)

	sinkid int

	callerline int
	callerfile string
}

type Index struct {
	I int
	O float64
}

func (s Sorm) auxkids(wo *World) []Sorm {
	return wo.auxpool[s.kidsl:s.kidsr]
}

func (s Sorm) kids2(wo *World) []Sorm {
	return wo.pool[s.kidsl:s.kidsr]
}

func (wo *World) beginvirtual() (pool []Sorm) {
	if sameslice(wo.pool, wo.auxpool) {
		panic(`contraption: nested Sequences are not allowwed`)
	}
	pool = wo.pool
	wo.pool = wo.auxpool
	wo.nextn, wo.auxn = wo.auxn, wo.nextn
	return
}

func (wo *World) endvirtual(pool []Sorm) {
	wo.auxpool = wo.pool
	wo.pool = pool
	wo.nextn, wo.auxn = wo.auxn, wo.nextn
}

func (s Sorm) kidsiter(wo *World, f func(*Sorm)) {
	// TODO Idea: take a limit in kidsiter, and if it is scissored, stop iteration when over it.
	// It should work since scissored compounds can't stretch kids.
	if s.tag == tagSequence {
		return
	}
	kids := wo.pool[s.kidsl:s.kidsr]
out:
	for i := range kids {
		k := &kids[i]

		q, ok := k.key.(Sequence)
		if ok {
			if k.flags&flagSequenceSaved == 0 {
				args := wo.tmpalloc(q.Length(wo))
				reall := len(wo.auxpool)

				// Treat the aux pool as a main pool and a sequence as a root compound.
				pop := wo.beginvirtual()
				wo.prefix = k.z
				for i := 0; i < q.Length(wo); i++ {
					t := q.Get(wo, i)
					args[i] = t
				}
				wo.prefix = 0
				wo.endvirtual(pop)

				// Copy the elements materialized from sequence to the aux pool,
				// treat them like arguments of (*World).Compound
				l, r := wo.allocaux(q.Length(wo))
				copy(wo.auxpool[l:r], args)
				k.kidsl = reall
				k.kidsr = r
				// Save immediate kids.
				k.presl = l
				k.presr = r
				k.flags |= flagSequenceSaved
			}
			aux := wo.auxpool[k.presl:k.presr]
			for i := range aux {
				k := &aux[i]
				pop := wo.beginvirtual()
				f(k) // (1)
				wo.endvirtual(pop)
				if k.flags&flagBreakIteration > 0 { // (2)
					break out
				}
			}
		} else {
			f(k)                                // (1)
			if k.flags&flagBreakIteration > 0 { // (2)
				break out
			}
		}
	}
}

func (s Sorm) mods(wo *World) []Sorm {
	return wo.pool[s.modsl:s.modsr]
}

func (s Sorm) pres(wo *World) []Sorm {
	return wo.pool[s.presl:s.presr]
}

// allocmain allocates new memory in pool and returns index range for an object.
func (wo *World) allocmain(n int) (left, right int) {
	return alloc(&wo.pool, n)
}

// allocaux is wo.alloc for wo.auxpool.
func (wo *World) allocaux(n int) (left, right int) {
	return alloc(&wo.auxpool, n)
}

func alloc(pool *[]Sorm, n int) (left, right int) {
	if len(*pool)+n > cap(*pool) {
		*pool = append(*pool, make([]Sorm, n)...)
	} else {
		*pool = (*pool)[0 : len(*pool)+n]
	}
	return len(*pool) - n, len(*pool)
}

// This allocator is not correct, because it gets more memory
// if capacity has changed, thus making all previously allocated slices
// be out of pool.
//
// But in this case it don't matter. These objects are not going to the pool.
//
// wo.alloc() which previously looked like this is fixed now by using indices in
// Sorms instead of slices.
func (wo *World) tmpalloc(n int) []Sorm {
	if len(wo.tmp)+n > cap(wo.tmp) {
		wo.tmp = append(wo.tmp, make([]Sorm, n)...)
	} else {
		wo.tmp = wo.tmp[0 : len(wo.tmp)+n]
	}
	return wo.tmp[len(wo.tmp)-n:]
}

func (wo *World) stash(s []Sorm) []Sorm {
	dst := wo.tmpalloc(len(s))
	copy(dst, s)
	return dst
}

func (s Sorm) String() string {
	type Paint struct {
		xform      nanovgo.TransformMatrix
		extent     [2]float32
		radius     float32
		feather    float32
		innerColor nanovgo.Color
		outerColor nanovgo.Color
		image      int
	}
	p := (*Paint)(unsafe.Pointer(&s.fill))
	color := func(c nanovgo.Color) string {
		return fmt.Sprintf("#%2.2x%2.2x%2.2x%2.2x", int(c.R*255), int(c.G*255), int(c.B*255), int(c.A*255))
	}

	f := func(f float64) string {
		return fmt.Sprintf("%g", math.Ceil(f*100)/100)
	}
	var vals string
	if s.tag < 0 {
		vals = fmt.Sprint(f(s.W), ", ", f(s.H), ", ", f(s.x), ", ", f(s.y), ", _", f(s.wl), ", _", f(s.hl))
	} else {
		vals = fmt.Sprint(f(s.W), "×", f(s.H), ", ", f(s.x), "y", f(s.y), ", ↓", f(s.wl), "×", f(s.hl), " ", color(p.innerColor))
	}
	key := cond(s.key != nil, fmt.Sprint(" [", s.key, "]"), "")
	ovrx := cond(s.flags&flagOvrx > 0, "↑X", "  ")
	ovry := cond(s.flags&flagOvry > 0, "↑Y", "  ")
	btw := cond(s.flags&flagBetweener > 0, "↔", " ")
	seq := " "
	z := ""
	if s.flags&flagSequenceMark > 0 {
		seq = "S"
		z = sprint([2]int{s.z, s.z2})
	} else {
		z = sprint(s.z)
	}

	d := func(i int) int {
		return int(math.Floor(math.Log10(float64(i))))
	}
	digits := d(s.z) + max(d(s.z2), 0)
	if seq == "S" {
		digits += 2 // "[" and "]"
	}
	scissor := cond(s.scissor.Dx() > 0 && s.scissor.Dy() > 0, fmt.Sprint(s.scissor), "{/}")

	return fmt.Sprint(z, strings.Repeat(" ", max(0, 10-digits)), seq, ovrx, ovry, btw, " ", s.tag.String(), " ", vals, ` `, s.props, ` `, scissor, key, " ", s.callerfile, ":", s.callerline)
}

func (s Sorm) decimate() Sorm {
	s.x = math.Floor(s.x)
	s.y = math.Floor(s.y)
	s.W = math.Ceil(s.W)
	s.H = math.Ceil(s.H)
	return s
}

func (s Sorm) Fill(p nanovgo.Paint) Sorm {
	s.fill = p
	return s
}

func (s Sorm) Stroke(p nanovgo.Paint) Sorm {
	s.stroke = p
	return s
}

func (s Sorm) Strokewidth(w float32) Sorm {
	s.strokew = w
	s.flags |= flagSetStrokewidth
	return s
}

func (s Sorm) FillStroke(p nanovgo.Paint) Sorm {
	s.stroke = p
	s.fill = p
	return s
}

func (s Sorm) CondFill(f func(rect geom.Rectangle) nanovgo.Paint) Sorm {
	s.condfill = f
	return s
}

func (s Sorm) CondStroke(f func(rect geom.Rectangle) nanovgo.Paint) Sorm {
	s.condstroke = f
	return s
}

func (s Sorm) CondFillStroke(f func(rect geom.Rectangle) (nanovgo.Paint, nanovgo.Paint)) Sorm {
	s.condfillstroke = f
	return s
}

func (s Sorm) Cond(f func(Matcher)) Sorm {
	s.cond = f
	return s
}

func (s Sorm) Lmb(wo *World, f func()) Sorm {
	s.cond = func(m Matcher) {
		if m.Nochoke().Match(`Click(1):in`) {
			f()
		}
	}
	return s
}

func (s Sorm) Override() Sorm {
	s.flags |= flagOvrx | flagOvry
	return s
}

func (s Sorm) Hoverride() Sorm {
	s.flags |= flagOvrx
	return s
}

func (s Sorm) Voverride() Sorm {
	s.flags |= flagOvry
	return s
}

func (s Sorm) Betweener() Sorm {
	s.flags |= flagBetweener
	return s
}

func (s Sorm) Size(x, y float64) Sorm {
	s.W = x
	s.H = y
	return s
}

func (s Sorm) Rectangle() geom.Rectangle {
	return geom.Rect(s.x, s.y, s.x+s.W, s.y+s.H)
}

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

// BaseWorld returns itself.
// This method allows to access base World class from user worlds.
func (wo *World) BaseWorld() *World {
	return wo
}

func (wo *World) newSorm() (s Sorm) {
	wo.nextn++

	s = Sorm{
		z:    wo.nextn,
		m:    geom.Identity2d(),
		prem: geom.Identity2d(),
	}
	if wo.f1 {
		_, s.callerfile, s.callerline, _ = runtime.Caller(2)
	}
	virtual := wo.prefix > 0
	if virtual {
		s.z = wo.prefix
		s.z2 = wo.nextn
		s.flags = flagSequenceMark
	}
	return
}

func (wo *World) Prevkey(key any) Sorm {
	for _, z := range wo.old {
		if z.tag == 0 && z.key == key {
			return z
		}
	}
	return Sorm{}
}

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
	s.W = clamp(0, amt, 1)
	return
}
func halignrun(wo *World, c, m *Sorm) {
	x := 0.0
	c.kidsiter(wo, func(k *Sorm) {
		x = max(x, k.W)
	})
	c.W = max(c.W, x)
	c.kidsiter(wo, func(k *Sorm) {
		k.x += (x - k.W) * m.W
		// c.h = max(c.h, k.h)
	})
}

func (wo *World) Valign(amt float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagValign
	s.W = clamp(0, amt, 1)
	return
}
func valignrun(wo *World, c, m *Sorm) {
	y := 0.0
	c.kidsiter(wo, func(k *Sorm) {
		y = max(y, k.H)
	})
	c.H = max(c.H, y)
	c.kidsiter(wo, func(k *Sorm) {
		k.y += (y - k.H) * m.W
	})
}

func (wo *World) Fill(p nanovgo.Paint) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagFill
	s.fill = p
	return
}
func fillrun(wo *World, s, m *Sorm) {
	s.kidsiter(wo, func(k *Sorm) {
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
	s.kidsiter(wo, func(k *Sorm) {
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
	s.kidsiter(wo, func(k *Sorm) {
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
	s.W, s.H = w, h
	return
}
func limitrun(wo *World, s, m *Sorm) {
	p := s.m.ApplyPt(geom.Pt(m.W, m.H))
	if m.W > 0 {
		m.W = p.X
		s.wl = min(s.wl, m.W)
	} else if m.W < 0 {
		s.eprops.X = -m.W
	}
	if m.H > 0 {
		m.H = p.Y
		s.hl = min(s.hl, m.H)
	} else if m.H < 0 {
		s.eprops.Y = -m.H
	}
}

// Transform applies transformation that only affects objects visually.
// It doesn't affect object sizes for layout.
func (wo *World) Transform(x, y float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagTransform
	s.W = x
	s.H = y
	return
}
func transformrun(wo *World, c, m *Sorm) {
	c.x += m.W
	c.y += m.H
	// Because moves are inherited in a separate pass
}

// Pretransform applies transformation that affects objects sizes for layout.
func (wo *World) Pretransform(m geom.Geom) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagPretransform
	s.m = m
	return
}
func pretransformrun(wo *World, c, m *Sorm) {
	c.m = c.m.Mul(m.m)
	// c.sx *= m.W
	// c.sy *= m.H
}

// If a Compound has no aligner set (“stack” layout)
func noaligner(wo *World, c *Sorm) {
	minp := Point{}
	c.kidsiter(wo, func(k *Sorm) {
		c.W = max(c.W, k.W)
		c.H = max(c.H, k.H)
		if k.W < 0 {
			minp.X = min(minp.X, k.W)
		}
		if k.H < 0 {
			minp.Y = min(minp.Y, k.H)
		}
	})
	if c.flags&flagHshrink > 0 {
		c.wl = c.W
	}
	if c.flags&flagVshrink > 0 {
		c.hl = c.H
	}
	c.kidsiter(wo, func(k *Sorm) {
		stretch := false
		if k.W < 0 {
			k.W = c.wl * k.W / minp.X
			k.W += k.addw
			stretch = true
		}
		if k.H < 0 {
			k.H = c.hl * k.H / minp.Y
			k.H += k.addh
			stretch = true
		}
		if stretch {
			k.hl = k.H
			k.wl = k.W
			wo.apply(c, k)
		}
	})
	c.kidsiter(wo, func(k *Sorm) {
		c.W = max(c.W, k.W)
		c.H = max(c.H, k.H)
	})
}

func vfollowaligner(wo *World, c *Sorm) {
	sequencealigner(wo, c, false)
}

func hfollowaligner(wo *World, c *Sorm) {
	sequencealigner(wo, c, true)
}

func axis(h bool) (begin, end func(k *Sorm)) {
	swap := func(k *Sorm) {
		if h {
			// Y is main axis, X is secondary.
			k.W, k.H = k.H, k.W
			k.x, k.y = k.y, k.x
			k.wl, k.hl = k.hl, k.wl
			k.props.X, k.props.Y = k.props.Y, k.props.X
			k.eprops.X, k.eprops.Y = k.eprops.Y, k.eprops.X
			k.known.X, k.known.Y = k.known.Y, k.known.X
		}
	}
	return swap, swap
}

func trypos(p, a geom.Point) geom.Point {
	if p.X == 0 {
		p.X = a.X
	}
	if p.Y == 0 {
		p.Y = a.Y
	}
	return p
}

func sequencedivider(wo *World, c *Sorm, h bool) {
	beginaxis, endaxis := axis(h)
	beginaxis(c)
	c.kidsiter(wo, func(k *Sorm) {
		beginaxis(k)
		if k.tag == 0 {
			c.props.X = max(c.props.X, k.eprops.X)
			c.props.Y += k.props.Y
		} else {
			k.eprops = geom.Pt(-min(0, k.W), -min(0, k.H))
			if k.W >= 0 {
				c.known.X = max(c.known.X, k.W)
			} else {
				c.props.X = max(c.props.X, -k.W)
			}
			if k.H >= 0 {
				c.known.Y += k.H
			} else {
				c.props.Y += -k.H
			}
		}
		endaxis(k)
	})
	// Don't set if overriden by Limit with negative size.
	c.eprops = trypos(c.eprops, c.props)
	endaxis(c)
}

func sequencealigner(wo *World, c *Sorm, h bool) {
	c.W, c.H = 0, 0
	beginaxis, endaxis := axis(h)

	// Calculate unknowns and apply kids which sizes were unknown,
	// lay out the sequence then.
	y := 0.0
	beginaxis(c)
	c.kidsiter(wo, func(k *Sorm) {
		beginaxis(k)
		stretch := false
		if k.eprops.Y > 0 {
			// max() means we don't stretch if we're out of limit.
			k.H = max(0, (c.hl-c.known.Y)/c.props.Y*k.eprops.Y)
			k.hl = k.H
			stretch = true
		}
		if k.eprops.X > 0 {
			k.W = c.wl / c.props.X * k.eprops.X
			k.wl = k.W
			stretch = true
		}
		endaxis(k)

		if stretch {
			wo.apply(c, k)
		}

		beginaxis(k)
		k.y = y
		y += k.H
		c.W = max(c.W, k.W)
		endaxis(k)
	})
	c.H = y
	endaxis(c)
}

func (wo *World) apply(_ *Sorm, c *Sorm) {
	// Apply is called on every shape, so ignore anything that is not Compound.
	if c.tag != 0 {
		return
	}

	// apply presumes that premodifiers are already sorted.
	for _, m := range c.pres(wo) {
		// A loop in (*World).resolvealigners is idemponent to this one
		// in the case of tagVfollow and tagHfollow.
		if m.tag == tagLimit && (m.W < 0 || m.H < 0) {
			continue
		}
		preActions[-100-m.tag](wo, c, &m)
	}

	// Set scissor to limit if needed.
	if c.flags&flagScissor > 0 {
		c.scissor = geom.Rect(0, 0, c.wl, c.hl)
	}
	c.kidsiter(wo, func(k *Sorm) {
		// NOTE Aligner is called after these assignments.
		// 	So this can't influence limits at later stages.
		k.wl = c.wl
		k.hl = c.hl
		// Inherit scissors and apply scale to them.
		k.scissor = c.scissor

		// Apply scale.
		ns := k.m.ApplyPt(geom.Pt(k.W, k.H))
		ims := k.m.ApplyPt(geom.Pt(k.addw, k.addh))
		// Don't scale stretch coefficients.
		k.W = cond(k.W >= 0, ns.X, k.W)
		k.H = cond(k.H >= 0, ns.Y, k.H)
		// But scale imaginaries.
		k.addw = ims.X
		k.addh = ims.Y

		// Process only kids which sizes are known first.
		if k.W >= 0 && k.H >= 0 {
			wo.apply(c, k)
		}
	})

	alignerActions[c.aligner](wo, c)

	// Apply postorder/anyorder modifiers.
	for _, m := range c.mods(wo) {
		modActions[-m.tag](wo, c, &m)
	}

	c.kidsiter(wo, func(k *Sorm) {
		// Apply size override to c if it has one.
		if k.flags&flagOvrx > 0 {
			c.W = k.W
			c.x = min(c.x, -k.x)
		}
		if k.flags&flagOvry > 0 {
			c.H = k.H
			c.y = min(c.y, -k.y)
		}
	})

	if c.flags&flagScissor > 0 {
		c.W = min(c.W, c.wl)
		c.H = min(c.H, c.hl)
	}
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
// The value is deleted if it has been not accessed for two frames.
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

// Cat returns a Compound from two Sorm sources.
// Its intended usage is for defining user's own abstractions.
func (wo *World) Cat(a, b []Sorm) (s Sorm) {
	tmp := wo.tmpalloc(len(a) + len(b))
	copy(tmp[:len(a)], a)
	copy(tmp[len(a):], b)
	return wo.compound(wo.newSorm(), tmp...)
}

// Compound is a shape container.
// It combines multiple Sorms into a single shape.
// Every other container is just a rebranded Compound.
func (wo *World) Compound(args ...Sorm) (s Sorm) {
	return wo.compound(wo.newSorm(), args...)
}

func (wo *World) realroot(args ...Sorm) Sorm {
	return wo.compound2(wo.newSorm(), true, args...)
}

func (wo *World) compound(s Sorm, args ...Sorm) Sorm {
	return wo.compound2(s, false, args...)
}

func (wo *World) compound2(s Sorm, realroot bool, args ...Sorm) Sorm {
	var kidc, modc, prec, btwc int
	var void func() Sorm

	if realroot {
		s.wl = wo.Wwin
		s.hl = wo.Hwin
	}

	for _, a := range args {
		if a.tag == tagBetween {
			void = a.key.(func() Sorm)
		}
		switch {
		case a.tag >= 0:
			kidc++
			if a.flags&flagBetweener > 0 {
				btwc++
			}
		case a.tag < 0 && a.tag >= -100:
			modc++
		case a.tag < -100:
			prec++
		}
	}
	// TODO Allocate only what needed by index and limit.
	// Between-shape allocation pass
	tmpn := len(wo.tmp)
	voidc := max(kidc-btwc-1, 0)
	if void != nil {
		kidc = kidc + voidc
		// Place new voids into temporary storage so their allocations won't
		// break breadth-first order of the pool.
		for i := 0; i < voidc; i++ {
			wo.tmp = append(wo.tmp, void())
		}
	}

	// s = wo.newCompound()
	s.tag = 0

	// Shape pass
	s.kidsl, s.kidsr = wo.alloc(kidc)
	i := 0
	for _, a := range args {
		if a.tag >= 0 {
			s.kids2(wo)[i] = a
			i++
			if void != nil && voidc > 0 && !(a.flags&flagBetweener > 0) {
				s.kids2(wo)[i] = wo.tmp[tmpn]
				voidc--
				tmpn++
				i++
			}
		}
	}

	// Modifier pass
	s.presl, s.presr = wo.alloc(prec)
	s.modsl, s.modsr = wo.alloc(modc)
	i, j := 0, 0
	for _, a := range args {
		switch {
		case a.tag < -100:
			s.pres(wo)[i] = a
			i++
		case a.tag < 0:
			s.mods(wo)[j] = a
			j++
		}
	}

	return s
}

func (wo *World) Root(s ...Sorm) {
	wo.Compound(
		wo.Void(complex(wo.Wwin, 0), complex(wo.Hwin, 0)),
		func() Sorm {
			if wo.DragEffect != nil && wo.drag != nil {
				ps := [2]geom.Point{wo.dragstart, wo.Trace[0].Pt}
				return wo.DragEffect(ps, wo.drag)
			}
			return Sorm{}
		}(),
		wo.displayOscilloscope(),
		wo.realroot(s...))
	wo.rend = len(wo.pool)
}

// Next prepares the Contraption for rendering the next frame.
// See package description for preferred use of Contraption.
func (wo *World) Next() bool {
	wo.MatchCount = 0

	window := wo.Window
	if window.ShouldClose() {
		return false
	}

	// TODO Decouple backend and put more stuff in that .next().
	wo.Events.next()

	w, h := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))

	w, h = window.GetSize()
	wo.Wwin, wo.Hwin = float64(w), float64(h)
	wo.Events.Viewport = geom.Pt(float64(w), float64(h))

	cl := hex(`#ffffff`)
	gl.ClearColor(cl.R, cl.G, cl.B, cl.A)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)

	sc, _ := window.GetContentScale()
	wo.Vgo.BeginFrame(w, h, float32(math.Ceil(float64(sc))))

	wo.Vgo.SetFontFace(`go`)
	wo.Vgo.SetFontSize(17) // FIXME Cap size of Go 17 is 11

	return true
}

func (wo *World) resolvealigners(_ *Sorm, c *Sorm) {
	if c.tag != 0 {
		return
	}
	// Premodifiers have order of execution.
	slices.SortFunc(c.pres(wo), func(a, b Sorm) int {
		return int(b.tag - a.tag)
	})
	// A premodifier loop in (*World).apply is idemponent to this one.
	for _, m := range c.pres(wo) {
		if m.tag == tagLimit && (m.W < 0 || m.H < 0) {
			preActions[-100-m.tag](wo, c, &m)
		}
		if oneof(m.tag, tagVfollow, tagHfollow) {
			preActions[-100-m.tag](wo, c, &m)
		}
	}
	c.kidsiter(wo, func(k *Sorm) {
		wo.resolvealigners(c, k)
	})
}

// Develop applies the layout and renders the next frame.
// See package description for preferred use of Contraption.
func (wo *World) Develop() {
	vgo := wo.Vgo

	vgo.ResetTransform()
	root := &wo.pool[len(wo.pool)-1]
	root.wl = wo.Wwin
	root.hl = wo.Hwin

	pool := wo.pool[0:wo.rend]

	for i := range pool {
		pool[i].i = i
	}

	// Resolve aligners and stack negative sizes.
	wo.resolvealigners(nil, last(pool))
	for i := 0; i < len(pool); i++ {
		s := &pool[i]
		if s.tag == tagSequence {
			p := wo.beginvirtual()
			pool := wo.auxpool[s.kidsl:s.kidsr]
			for i := len(pool) - 1; i >= 0; i-- {
				s := &pool[i]
				switch s.aligner {
				case alignerVfollow:
					sequencedivider(wo, s, false)
				case alignerHfollow:
					sequencedivider(wo, s, true)
				}
			}
			wo.endvirtual(p)
		} else {
			switch s.aligner {
			case alignerVfollow:
				sequencedivider(wo, s, false)
			case alignerHfollow:
				sequencedivider(wo, s, true)
			}
		}
	}

	if wo.Events.Match(`Press(F4)`) {
		println("@ Tree before applying stretches")
		sormp(wo, *last(pool), 0)
		println()
	}

	// Do the layout.
	wo.apply(nil, last(pool))

	// Inherit moves and paints.
	// TODO Maybe wo.apply should do it? Kind of makes more sense.
	inh := func(c, efc *Sorm) {
		c.kidsiter(wo, func(k *Sorm) {
			k.x += efc.x
			k.y += efc.y
			k.scissor = k.scissor.Add(geom.Pt(efc.x, efc.y))
			if k.fill == (nanovgo.Paint{}) {
				k.fill = efc.fill
			}
			if k.stroke == (nanovgo.Paint{}) {
				k.stroke = efc.stroke
			}
			if k.flags&flagSetStrokewidth == 0 {
				k.strokew = efc.strokew
			}
		})
	}
	for i := len(pool) - 1; i >= 0; i-- {
		s := &pool[i]
		if s.tag == tagSequence {
			p := wo.beginvirtual()
			pool := wo.auxpool[s.kidsl:s.kidsr]
			for i := len(pool) - 1; i >= 0; i-- {
				s := &pool[i]
				inh(s, s)
			}
			wo.endvirtual(p)
		} else {
			inh(s, s)
		}
	}

	// Print tree for debug. Do it before sorting.
	if wo.f1 {
		println("@ Tree after layout")
		sormp(wo, *last(pool), 0)
		println()
		wo.f1 = false
	}
	if wo.Events.Match(`Press(F1)`) {
		wo.f1 = true
	}

	if wo.Events.Match(`Press(F6)`) {
		println("@ Auxpool")
		for i := range wo.auxpool {
			println(i, wo.auxpool[i])
		}
		println()
	}

	// Sort in draw order.
	slices.SortFunc(pool, func(a, b Sorm) int {
		return a.z - b.z
	})

	// Apply conditional paints.
	for i := len(pool) - 1; i >= 0; i-- {
		s := &pool[i]
		r := geom.Rect(s.x, s.y, s.x+s.W, s.y+s.H)
		if s.condfillstroke != nil {
			s.fill, s.stroke = s.condfillstroke(r)
		}
		if s.condfill != nil {
			s.fill = s.condfill(r)
		}
		if s.condstroke != nil {
			s.stroke = s.condstroke(r)
		}
		m := wo.Events.In(r).WithZ(i + 1)
		if s.flags&flagSource > 0 {
			if m.Match(`Click(1):in`) {
				wo.drag = s.key
				wo.dragstart = wo.Last.FirstTouch
			}
		}
		if s.cond != nil {
			s.cond(m)
		}
		if s.sinkid > 0 {
			if m.Match(`Unclick(1):in`) && wo.drag != nil {
				wo.sinks[s.sinkid](wo.drag)
				wo.drag = nil
			}
		}
	}

	// Final drag match — resets drag if it was dropped in nowhere.
	if wo.Anywhere().Match(`!Click(1)* Unclick(1)`) {
		wo.drag = nil
	}

	// Draw.
	draw := func(s Sorm) {
		s = s.decimate()
		// Set scissor up.
		vgo.ResetScissor()
		if s.scissor.Dx() > 0 && s.scissor.Dy() > 0 {
			x, y, w, h := rect2nvgxywh(s.scissor)
			x = float32(math.Floor(float64(x)))
			y = float32(math.Floor(float64(y)))
			w = float32(math.Ceil(float64(w)))
			h = float32(math.Ceil(float64(h)))
			x -= 0.5
			y -= 0.5
			vgo.Scissor(x, y, w, h)
		}

		// Positioning bodges
		switch s.tag {
		default:
			// This fix is needed by every vector drawing library i know.
			// And i only know Nanovg and Love2d.
			s.x -= 0.5
			s.y -= 0.5
		case tagTopDownText:
			fallthrough
		case tagBottomUpText:
			fallthrough
		case tagText:
			// s.x -= 1
			// s.y -= 1

		case tagVectorText:
		case tagCanvas:
		}
		shapeActions[s.tag](wo, &s)
		vgo.ResetTransform()
	}
	for _, s := range pool {
		if s.tag <= 0 {
			continue
		}
		if s.tag == tagSequence {
			// At this moment every shape inside a sequence is cached,
			// just read it directly.
			// TODO This approach can be extended with recursion and a list of auxpools to support nested sequences.
			pool := wo.auxpool[s.kidsl:s.kidsr]
			for _, s := range pool {
				if s.tag <= 0 {
					continue
				}
				draw(s)
			}
		} else {
			draw(s)
		}
	}

	wo.Vgo.Reset()

	if wo.Match(`Press(F2)`) {
		wo.showOutlines = !wo.showOutlines
	}
	if wo.showOutlines {
		vgo.SetStrokeWidth(1)
		vgo.SetStrokePaint(hexpaint(`#00000020`))
		vgo.BeginPath()
		for _, s := range pool {
			if s.H < .5 && s.W < .5 {
				continue
			}
			s := s.decimate()
			s.x -= 0.5
			s.y -= 0.5
			if s.tag >= 0 {
				wo.Vgo.Rect(float32(s.x), float32(s.y), float32(s.W), float32(s.H))
			}
		}
		vgo.ClosePath()
		vgo.Stroke()
	}

	wo.recorder()

	if wo.Events.Match(`!Release(Ctrl)* Press(Ctrl)`) {
		if wo.Events.Match(`!Release(Shift)* Press(Shift)`) {
			if wo.Events.Match(`Press(I)`) {
				wo.Oscilloscope.on = !wo.Oscilloscope.on
			}
		}
	}

	// Note that after sorting pool by order it can't be used to
	// correctrly determine relationships.
	// Hence it is applicable to the old pool, the only reason it is
	// saved is to preserve keys.
	wo.old, wo.pool = wo.pool, wo.old
	wo.auxold, wo.auxpool = wo.auxpool, wo.auxold
	zeroandclear := func(pool *[]Sorm) {
		for i := range *pool {
			(*pool)[i] = Sorm{}
		}
		(*pool) = (*pool)[:0]
	}

	zeroandclear(&wo.pool)
	zeroandclear(&wo.auxpool)
	zeroandclear(&wo.tmp)
	wo.nextn = 0
	wo.auxn = 0

	wo.sinks = wo.sinks[:1]

	wo.Events.develop()
	wo.windowDevelop()

	for k, v := range wo.keys {
		v.counter--
		if v.counter <= 0 {
			delete(wo.keys, k)
		}
	}
}

func (wo *World) recorder() {
	vgo := wo.Vgo

	if wo.Events.Match(`Press(F5)`) {
		if wo.RecordPath == `` {
			panic(`can't record events, (*World).Events.RecordPath is empty`)
		}
		wo.Events.rec = 1 - wo.Events.rec
		if wo.Events.rec == 0 {
			f, err := os.Create(wo.RecordPath)
			if err != nil {
				panic(err)
			}
			err = gob.NewEncoder(f).Encode(wo.Events.records)
			if err != nil {
				panic(err)
			}
			err = f.Close()
			if err != nil {
				panic(err)
			}
		}
	}

	switch wo.Events.rec {
	case 1:
		vgo.SetFillColor(hex(`#ff0000`))
		vgo.BeginPath()
		vgo.Circle(12, 12, 9)
		vgo.Fill()
	case 2:
		vgo.SetFillColor(hex(`#ff0000`))
		vgo.BeginPath()
		vgo.MoveTo(6, 6)
		vgo.LineTo(6+16, 6+8)
		vgo.LineTo(6, 6+16)
		vgo.ClosePath()
		vgo.Fill()

		vgo.BeginPath()
		t := wo.Events.Trace[0].Pt
		vgo.Circle(float32(t.X), float32(t.Y), 2)
		vgo.Fill()

		fid := wo.nvgofontids[0]
		vgo.SetFontFaceID(fid)
		vgo.SetFillColor(hex(`#000000`))
		vgo.SetFontSize(float32(12 * wo.capmap[fid]))
		vgo.Text(16+8, 16+5, wo.Events.Trace[0].valuestring())
		adv, _ := vgo.TextBounds(16+8, 16+5, wo.Events.Trace[0].valuestring())
		vgo.SetFillColor(hex(`#00000050`))
		vgo.Text(16+8+adv+8, 16+5, wo.future.valuestring())

	case 3:
		vgo.SetFillColor(hex(`#ff0000`))
		vgo.BeginPath()
		vgo.Rect(6, 6, 16, 16)
		vgo.ClosePath()
		vgo.Fill()
	}
}

func sormp(wo *World, s Sorm, tab int) {
	println(fmt.Sprint(strings.Repeat("| ", tab), s))
	for _, s := range s.pres(wo) {
		sormp(wo, s, tab+1)
	}
	for _, s := range s.mods(wo) {
		sormp(wo, s, tab+1)
	}
	s.kidsiter(wo, func(k *Sorm) {
		sormp(wo, *k, tab+1)
	})
}

// Activator — это любой объект, на котором может быть сконцентрирован
// фокус ввода.
// Этот объект может обрабатывать свои события внутри метода Activate.
// В системе есть какое-то глобальное поле, которое обозначает
// текущий сфокусированный объект.
//
// Activate возвращает Silence, когда объект не среагировал на событие,
// Ack, когда среагировал и Deactivate, когда объект сбросил с
// себя фокус.
//
// Activator временный и умирает после смены фокуса.
type Activator interface {
	Activate(events *Events) (action Symbol)
}

// ActivatorPainter — это Activator, который может рисовать себя
// во время Drag.
type ActivatorPainter interface {
	Activator
	Paint(wo *World) Sorm
}

type Config struct {
	// Default window frame.
	// Effective if not zero.
	WindowRect image.Rectangle

	// If not nil, events are not received from user, but replayed from this reader
	// as recorded after pressing F5.
	ReplayReader io.Reader
}

type Window struct {
	FPS int
	window
}

func (wi *Window) Rect() image.Rectangle {
	return wi.rect()
}

func New(config Config) (wo *World) {
	if config.WindowRect.Dx() == 0 {
		config.WindowRect.Max.X = config.WindowRect.Min.X + 1024
	}
	if config.WindowRect.Dy() == 0 {
		config.WindowRect.Max.Y = config.WindowRect.Min.Y + 768
	}

	runtime.LockOSThread()

	wo = &World{}
	concretenew(config, wo)

	wo.Events = NewEventTracer(wo.Window.window, config.ReplayReader)
	wo.sinks = make([]func(any), 1)
	wo.keys = map[any]*labelt{}
	wo.alloc = wo.allocmain
	return wo
}
