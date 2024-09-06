// Contraption: a simple framework for user interfaces.
//
// A good user interface framework must be an engine for a word processing game.
//
// TODO:
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
//		- Needs changes in GLFW or changing input library.
//	- Interactive views for very large 1d and 2d data: waveforms, giant Minecraft maps, y-log STFT frames, etc.
//		- Why? Try to display STFT of a music file using Matplotlib, then rescale the window. Enjoy the wait.
//	- Text area
//		- https://rxi.github.io/textbox_behaviour.html
//	- Sequence must be a special shape that pastes Sorms inside a compound, not being compound itself
//		- So wo.Text(io.RuneReader) could be Sequence
//		- Not clear how to reuse memory of pools in this case
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
//		- wo.Text(io.RuneReader) []Sorm
//			- Returns Knuth-Plass-ready stream of boxes, Glues and Penalties.
//			? How to insert anything in between symbols?
//		- wo.Cap(float64) (can't be negative)
//		- wo.Lsp(float64)
//		- Knuth-Plass
//			? Interpret negative sizes as glue.
// 			- func Hknuth(perline func() Sorm) Sorm
// 			- func Vknuth(perline func() Sorm) Sorm
// 			- func Glue(width, minus, plus float64) Sorm // Analogous to wo.Void() but undirectional.
// 			- func Penalty(replacewith func() Sorm, penalty float64) Sorm
//			- Void(0, y) is already a “strut”
//	- Grid layout
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
//		- Scissor
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
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/neputevshina/geom"
	"github.com/neputevshina/nanovgo"
	"golang.org/x/exp/slices"
)

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
)

//go:generate stringer -type=tagkind -trimprefix=tag
type tagkind int

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
	tagVfollow
	tagHfollow
	tagHshrink
	tagVshrink
	tagLimit
)

// NOTE This can actually be a single table, if needed.
var modActions [100]func(wo *World, compound, mod *Sorm)
var preActions [100]func(wo *World, compound, mod *Sorm)
var shapeActions [100]func(wo *World, shape *Sorm)
var alignerActions [100]func(wo *World, compound *Sorm)

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
	preActions[-100-tagLimit] = limitrun
	preActions[-100-tagScissor] = scissorrun
	preActions[-100-tagVfollow] = vfollowrun
	preActions[-100-tagHfollow] = hfollowrun
	preActions[-100-tagHshrink] = hshrinkrun
	preActions[-100-tagVshrink] = vshrinkrun

	alignerActions[0] = noaligner
	alignerActions[1] = vfollowaligner
	alignerActions[2] = hfollowaligner
}

// Sequence is the thing that can generate elements for a scroll-enabled compound.
//
// All the logic to differentiate scroll's elements must be in returned Sorms.
type Sequence interface {
	Get(i int) Sorm
	Length() int
}

type adhocSequence struct {
	get    func(i int) Sorm
	length func() int
}

func (s *adhocSequence) Get(i int) Sorm {
	return s.get(i)
}

func (s *adhocSequence) Length() int {
	return s.length()
}

func AdhocSequence(get func(i int) Sorm, length func() int) Sequence {
	return &adhocSequence{get: get, length: length}
}

func SliceSequence[T any](sl []T, produce func(T) Sorm) Sequence {
	return AdhocSequence(func(i int) Sorm { return produce(sl[i]) }, func() int { return len(sl) })
}

func SliceSequence2[T any](sl []T, produce func(int) Sorm) Sequence {
	return AdhocSequence(func(i int) Sorm { return produce(i) }, func() int { return len(sl) })
}

type Scrollptr struct {
	Index  int
	Offset float64
	y      float64

	Dirty bool
}

type Eqn func(pt geom.Point) (dist float64)

type Equation interface {
	Eqn(pt geom.Point) (dist float64)
	Size() geom.Point
}

type Sorm struct {
	z, i            int
	tag             tagkind
	flags           flagval
	W, H, r, wl, hl float64
	x, y            float64
	m, prem         geom.Geom
	aligner         int
	kidsl, kidsr,
	modsl, modsr,
	presl, presr int

	fill    nanovgo.Paint
	stroke  nanovgo.Paint
	strokew float32 // TODO

	fontid  int
	vecfont *Font

	// Some objects use key field for own purposes:
	// 	- Equation stores an Equation object
	// 	- Text stores a io.RuneReader
	// 	- Compound stores an Identity, which also works
	//	  out as Source's dropable object
	// 	- Between stores func() Sorm
	key any // TODO eqn and key must be one field, key can only be in compounds

	condfill       func(rect geom.Rectangle) nanovgo.Paint
	condstroke     func(rect geom.Rectangle) nanovgo.Paint
	condfillstroke func(rect geom.Rectangle) (nanovgo.Paint, nanovgo.Paint)
	cond           func(m Matcher)
	canvas         func(vg *nanovgo.Context, rect geom.Rectangle)

	sinkid int

	eyl float64

	callerline int
	callerfile string
}

func (s Sorm) kids(wo *World) []Sorm {
	return wo.pool[s.kidsl:s.kidsr]
}

func (s Sorm) mods(wo *World) []Sorm {
	return wo.pool[s.modsl:s.modsr]
}

func (s Sorm) pres(wo *World) []Sorm {
	return wo.pool[s.presl:s.presr]
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

	pad := func(s string) string {
		return s
		// return fmt.Sprintf("%.10s", s)
	}
	f := func(f float64) string {
		return fmt.Sprintf("%g", math.Ceil(f*100)/100)
	}
	var vals string
	if s.tag < 0 {
		vals = fmt.Sprint(f(s.W), ", ", f(s.H), ", ", f(s.x), ", ", f(s.y), ", _", f(s.wl), ", _", f(s.hl))
	} else {
		vals = fmt.Sprint(f(s.W), "×", f(s.H), ", ", f(s.x), "y", f(s.y), ", ↓", f(s.wl), "×", f(s.hl), " ", color(p.innerColor))
		// _ = color
		// vals = fmt.Sprint(p.xform)
	}
	key := cond(s.key != nil, fmt.Sprint(" [", s.key, "]"), "")
	ovrx := cond(s.flags&flagOvrx > 0, "↑X", "  ")
	ovry := cond(s.flags&flagOvry > 0, "↑Y", "  ")
	btw := cond(s.flags&flagBetweener > 0, "↔", " ")
	return fmt.Sprint(s.z, " ", ovrx, ovry, btw, " ", pad(s.tag.String()), " ", vals, key, " <", s.callerfile, ":", s.callerline, ">")
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

// // Get implements Sequence.
// func (s *Sorm) Get(i int) Sorm {
// 	if s.tag != tagCompound {
// 		return *s
// 	}
// 	return s.kids(s.wo)[i]
// }

// // Length implements Sequence.
// func (s *Sorm) Length() int {
// 	if s.tag != tagCompound {
// 		return 1
// 	}
// 	return len(s.kids(s.wo))
// }

type World struct {
	*Events
	Window       Window
	Oscilloscope Oscilloscope

	nextn int
	pool  []Sorm
	tmp   []Sorm
	old   []Sorm
	Vgo   *nanovgo.Context

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
}

// BaseWorld returns itself.
// This method allows to access base World class from user worlds.
func (wo *World) BaseWorld() *World {
	return wo
}

// alloc allocates new memory in pool and returns index range for an object.
func (wo *World) alloc(n int) (left, right int) {
	if len(wo.pool)+n > cap(wo.pool) {
		wo.pool = append(wo.pool, make([]Sorm, n)...)
	} else {
		wo.pool = wo.pool[0 : len(wo.pool)+n]
	}
	return len(wo.pool) - n, len(wo.pool)
}

// This allocator is not correct, because it gets more memory
// if capacity has changed, thus making all previously allocated slices
// be out of pool.
//
// But in this case it don't matter. These objects are not going in the pool.
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

func (wo *World) newSorm() Sorm {
	wo.nextn++
	_, file, line, _ := runtime.Caller(2)
	return Sorm{
		// wo:         wo,
		z:          wo.nextn,
		m:          geom.Identity2d(),
		prem:       geom.Identity2d(),
		callerfile: file,
		callerline: line,
	}
}

func (wo *World) Prevkey(key any) Sorm {
	for _, z := range wo.old {
		if z.tag == 0 && z.key == key {
			return z
		}
	}
	return Sorm{}
}

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
		s.W = s.vecfont.Measure(size, str)
		s.r = 0
		s.key = str
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

	Makealine(wo.Vgo, s.vecfont, s.H, s.key.([]rune))
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

func (wo *World) Rectangle(w, h float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagRect
	s.W = w
	s.H = h
	return
}
func rectrun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		wo.Vgo.Rect(float32(s.x), float32(s.y), float32(s.W), float32(s.H))
	})
}

func (wo *World) Roundrect(w, h, r float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagRoundrect
	s.W = w
	s.H = h
	s.r = r
	return
}
func roundrectrun(wo *World, s *Sorm) {
	s.paint(wo, func() {
		wo.Vgo.RoundedRect(float32(s.x), float32(s.y), float32(s.W), float32(s.H), float32(s.r))
	})
}

func (wo *World) Void(w, h float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagVoid
	s.W = w
	s.H = h
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

func (wo *World) Canvas(w, h float64, run func(vgo *nanovgo.Context, rect geom.Rectangle)) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagCanvas
	s.canvas = run
	s.W = w
	s.H = h
	return
}
func (wo *World) Canvas2(w, h float64, run func(vgo *nanovgo.Context, rect geom.Rectangle)) (s Sorm) {
	// FIXME Must be standard behavior of Canvas: scale by transform.
	s = wo.newSorm()
	s.tag = tagCanvas
	s.r = 1
	s.canvas = run
	s.W = w
	s.H = h
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

func (wo *World) Vfollow() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagVfollow
	return
}
func vfollowrun(wo *World, c, m *Sorm) {
	c.aligner = 1
}

func (wo *World) Hfollow() (s Sorm) {
	s = wo.newSorm()
	s.tag = tagHfollow
	return
}
func hfollowrun(wo *World, c, m *Sorm) {
	c.aligner = 2
}

func (wo *World) Halign(amt float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagHalign
	s.W = clamp(0, amt, 1)
	return
}
func halignrun(wo *World, c, m *Sorm) {
	x := 0.0
	for i := range c.kids(wo) {
		x = max(x, c.kids(wo)[i].W)
	}
	c.W = max(c.W, x)
	for i := range c.kids(wo) {
		k := &c.kids(wo)[i]
		k.x += (x - k.W) * m.W
		// c.h = max(c.h, k.h)
	}
}

func (wo *World) Valign(amt float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagValign
	s.W = clamp(0, amt, 1)
	return
}
func valignrun(wo *World, c, m *Sorm) {
	y := 0.0
	for i := range c.kids(wo) {
		y = max(y, c.kids(wo)[i].H)
	}
	c.H = max(c.H, y)
	for i := range c.kids(wo) {
		k := &c.kids(wo)[i]
		k.y += (y - k.H) * m.W
		// c.w = max(c.w, k.w)
	}
}

func (wo *World) Fill(p nanovgo.Paint) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagFill
	s.fill = p
	return
}
func fillrun(wo *World, s, m *Sorm) {
	for i := range s.kids(wo) {
		if s.kids(wo)[i].tag >= 0 {
			s.kids(wo)[i].fill = m.fill
		}
	}
}

func (wo *World) Stroke(p nanovgo.Paint) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagStroke
	s.stroke = p
	return
}
func strokerun(wo *World, s, m *Sorm) {
	for i := range s.kids(wo) {
		if s.kids(wo)[i].tag > 0 {
			s.kids(wo)[i].stroke = m.stroke
		}
	}
}

func (wo *World) Strokewidth(w float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagStrokewidth
	s.strokew = float32(w)
	return
}
func strokewidthrun(wo *World, s, m *Sorm) {
	s.flags |= flagSetStrokewidth
	for i := range s.kids(wo) {
		if s.kids(wo)[i].tag > 0 {
			s.kids(wo)[i].strokew = m.strokew
		}
	}
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
func (wo *World) BetweenVoid(w, h float64) (s Sorm) {
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

// Scissor limits the painting of a given compound to specified limits.
// TODO.
func (wo *World) Scissor(w, h float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagScissor
	s.W = w
	s.H = h
	return
}
func scissorrun(wo *World, s *Sorm, m *Sorm) {}

// Limit limits the maximum compound size to specified limits.
// If a given size is negative, it limits the corresponding size of a compound by
// the rules of negative units for shapes.
func (wo *World) Limit(w, h float64) (s Sorm) {
	s = wo.newSorm()
	s.tag = tagLimit
	s.W, s.H = w, h
	return
}
func limitrun(wo *World, s, m *Sorm) {
	p := s.m.ApplyPt(geom.Pt(m.W, m.H))
	m.W, m.H = p.X, p.Y
	if m.W != 0 {
		s.wl = m.W
	}
	if m.H != 0 {
		s.hl = m.H
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
	for i := range c.kids(wo) {
		k := &c.kids(wo)[i]
		c.W = max(c.W, k.W)
		c.H = max(c.H, k.H)
		if k.W < 0 {
			minp.X = min(minp.X, k.W)
		}
		if k.H < 0 {
			minp.Y = min(minp.Y, k.H)
		}
	}
	if c.flags&flagHshrink > 0 {
		c.wl = c.W
	}
	if c.flags&flagVshrink > 0 {
		c.hl = c.H
	}
	for i := range c.kids(wo) {
		k := &c.kids(wo)[i]
		stretch := false
		if k.W < 0 {
			k.W = c.wl * k.W / minp.X
			stretch = true
		}
		if k.H < 0 {
			k.H = c.hl * k.H / minp.Y
			stretch = true
		}
		if stretch {
			k.hl = k.H
			k.wl = k.W
			wo.apply(c, k)
		}
	}
	for i := range c.kids(wo) {
		k := &c.kids(wo)[i]
		c.W = max(c.W, k.W)
		c.H = max(c.H, k.H)
	}
}

func vfollowaligner(wo *World, c *Sorm) {
	sequencealigner(wo, c, false)
}

func hfollowaligner(wo *World, c *Sorm) {
	sequencealigner(wo, c, true)
}

func sequencealigner(wo *World, c *Sorm, h bool) {
	known := geom.Pt(0, 0)
	props := geom.Pt(c.W, c.H)
	c.W, c.H = 0, 0
	beginaxis := func(k *Sorm) {
		if h {
			// Y is main axis, X is secondary.
			k.W, k.H = k.H, k.W
			k.x, k.y = k.y, k.x
			k.wl, k.hl = k.hl, k.wl
		}
	}
	endaxis := beginaxis

	// Get total known sizes for each axis.
	for i := range c.kids(wo) {
		k := &c.kids(wo)[i]

		beginaxis(k)
		if k.W > 0 {
			known.X = max(known.X, k.W)
		} else {
			props.X = min(props.X, k.W)
		}
		if k.H >= 0 {
			known.Y += k.H
		} else {
			props.Y -= k.H
		}
		endaxis(k)
	}

	// If there was no limit set for the secondary axis, let it be
	// the biggest known size measured by it.
	// if known.X > 0 && c.W > 0 {
	// 	if c.wl == 0 {
	// 		c.wl = known.X
	// 	}
	// }

	// Calculate unknowns and apply kids which sizes were unknown,
	// lay out the sequence then.
	y := 0.0
	beginaxis(c)
	for i := range c.kids(wo) {
		k := &c.kids(wo)[i]

		beginaxis(k)
		stretch := false
		if k.H < 0 {
			k.H = max(0, (c.hl-known.Y)/props.Y*-k.H) // Don't stretch if we're out of limit.
			stretch = true
		}
		if k.W < 0 {
			k.W = c.wl / props.X * k.W
			stretch = true
		}
		endaxis(k)

		k.hl = k.H
		k.wl = k.W
		if stretch {
			wo.apply(c, k)
		}

		beginaxis(k)
		k.y += y
		y += k.H
		c.W = max(c.W, k.W)
		endaxis(k)
	}
	c.H = y
	endaxis(c)
}

func (wo *World) apply(p *Sorm, c *Sorm) {
	// Apply is called on every shape, so ignore anything that is not Compound.
	if c.tag != 0 {
		return
	}

	// Limit must be executed after update of a matrix, because it influences size.
	slices.SortFunc(c.pres(wo), func(a, b Sorm) int {
		return int(b.tag - a.tag)
	})
	for _, m := range c.pres(wo) {
		preActions[-100-m.tag](wo, c, &m)
	}

	// nl := p.m.ApplyPt(geom.Pt(c.wl, c.hl))
	// c.wl = cond(c.wl >= 0, nl.X, c.wl)
	// c.hl = cond(c.wl >= 0, nl.Y, c.hl)
	for i := range c.kids(wo) {
		k := &c.kids(wo)[i]
		k.wl = c.wl
		k.hl = c.hl

		// Apply scale.
		k.m = k.m.Mul(c.m)
		ns := k.m.ApplyPt(geom.Pt(k.W, k.H))
		// Don't scale stretch coefficients.
		k.W = cond(k.W >= 0, ns.X, k.W)
		k.H = cond(k.H >= 0, ns.Y, k.H)

		// Process only kids which sizes are known first.
		if k.W >= 0 && k.H >= 0 {
			wo.apply(c, k)
		}
	}

	alignerActions[c.aligner](wo, c)

	// Apply postorder/anyorder modifiers.
	for _, m := range c.mods(wo) {
		modActions[-m.tag](wo, c, &m)
	}

	for i := range c.kids(wo) {
		k := &c.kids(wo)[i]
		// Apply size override to c if it has one.
		if k.flags&flagOvrx > 0 {
			c.W = k.W
			c.x = min(c.x, -k.x)
		}
		if k.flags&flagOvry > 0 {
			c.H = k.H
			c.y = min(c.y, -k.y)
		}
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
	return wo.compound(wo.newSorm(), false, nil, tmp...)
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
		if s.tag == tagVoid && s.W == 0 && s.H == 0 {
			continue
		}
		// Skip Sorm{}
		if zero(s) {
			continue
		}
		tmp[j] = s
		j++
	}
	copy(tmp[j:], plus)
	return wo.compound(wo.newSorm(), false, nil, tmp...)
}

// Compound is a shape container.
// It combines multiple Sorms into a single shape.
// Every other container is just a rebranded Compound.
func (wo *World) Compound(args ...Sorm) (s Sorm) {
	return wo.compound(wo.newSorm(), false, nil, args...)
}

func (wo *World) compound(s Sorm, isvoid bool, void func() Sorm, args ...Sorm) Sorm {
	var kidc, modc, prec, btwc int
	for _, a := range args {
		if a.tag == tagBetween {
			isvoid = true
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
	// Void allocation pass
	tmpn := len(wo.tmp)
	voidc := max(kidc-btwc-1, 0)
	if isvoid {
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
			s.kids(wo)[i] = a
			i++
			if isvoid && voidc > 0 && !(a.flags&flagBetweener > 0) {
				s.kids(wo)[i] = wo.tmp[tmpn]
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

// TODO For use with Scroll()
// func sequencesequencealigner(wo *World, c *Sorm, seq Sequence, h bool) {
// 	known := geom.Pt(0, 0)
// 	props := geom.Pt(c.W, c.H)
// 	c.W, c.H = 0, 0
// 	beginaxis := func(k *Sorm) {
// 		if h {
// 			// Y is main axis, X is secondary.
// 			k.W, k.H = k.H, k.W
// 			k.x, k.y = k.y, k.x
// 			k.wl, k.hl = k.hl, k.wl
// 		}
// 	}
// 	endaxis := beginaxis

// 	// Get total known sizes for each axis.
// 	// How-to:
// 	//      - if total size by primary axis is greater than limit, escape
// 	// 	- get an element from sequence
// 	//	- materialize it and add size to primary axis
// 	//	- repeat
// 	//
// 	for i := range c.kids() {
// 		k := &c.kids()[i]

// 		beginaxis(k)
// 		if k.W > 0 {
// 			known.X = max(known.X, k.W)
// 		} else {
// 			props.X = min(props.X, k.W)
// 		}
// 		if k.H >= 0 {
// 			known.Y += k.H
// 		} else {
// 			props.Y -= k.H
// 		}
// 		endaxis(k)
// 	}

// 	// If there was no limit set for the secondary axis, let it be
// 	// the biggest known size measured by it.
// 	// if known.X > 0 && c.W > 0 {
// 	// 	c.wl = known.X
// 	// }

// 	// Calculate unknowns and apply kids which sizes were unknown,
// 	// lay out the sequence then.
// 	y := 0.0
// 	beginaxis(c)
// 	for i := range c.kids() {
// 		k := &c.kids()[i]

// 		beginaxis(k)
// 		stretch := false
// 		if k.H < 0 {
// 			k.H = max(0, (c.hl-known.Y)/props.Y*-k.H) // Don't stretch if we're out of limit.
// 			stretch = true
// 		}
// 		if k.W < 0 {
// 			k.W = c.wl / props.X * k.W
// 			stretch = true
// 		}
// 		endaxis(k)

// 		k.hl = k.H
// 		k.wl = k.W
// 		if stretch {
// 			wo.apply(c, k)
// 		}

// 		beginaxis(k)
// 		k.y += y
// 		y += k.H
// 		c.W = max(c.W, k.W)
// 		endaxis(k)
// 	}
// 	c.H = y
// 	endaxis(c)
// }

func (wo *World) Vscroll(ptr *Scrollptr, scrollpx float64, ylimit float64, seq Sequence) (s Sorm) {
	// FIXME Needs one frame latency to repaint and handle new events.
	prev := wo.Prevkey(ptr)
	if prev.tag == 0 {
		return
	}

	switch {
	case wo.Events.MatchIn(`Scroll(-1)`, prev.Rectangle()):
		scrollpx = -scrollpx
	case wo.Events.MatchIn(`Scroll(1)`, prev.Rectangle()):
		// No change
	default:
		scrollpx = 0
	}

	if ptr.Offset-scrollpx < 0 {
		// Request previous shape
		ptr.Index--
	}
	for i := ptr.Index; i <= seq.Length(); i++ {
		// position wo.apply(seq.Get(i)) to pool like overflow

		// do y summation

		// y += last(wo.pool).H
		// if y >= limit {
		// 	y = limit
		// 	break
		// }
	}
	panic(`unimplemented`)
}

func (wo *World) Root(s ...Sorm) {
	wo.Compound(
		wo.Void(wo.Wwin, wo.Hwin),
		func() Sorm {
			if wo.DragEffect != nil && wo.drag != nil {
				ps := [2]geom.Point{wo.dragstart, wo.Trace[0].Pt}
				return wo.DragEffect(ps, wo.drag)
			}
			return Sorm{}
		}(),
		wo.displayOscilloscope(),
		wo.Compound(s...))
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

	cl := hex(`#fefefe`)
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
	// setFontSize(vg, 11)

	rand.Seed(405)
	return true
}

func reachCheck(wo *World, pool []Sorm) {
	for i := range pool {
		k := &pool[i]
		k.flags |= flagMark
		reachCheck(wo, k.kids(wo))
	}
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

	reachCheck(wo, pool)
	for i := range pool {
		pool[i].i = i
		if !(pool[i].flags&flagMark > 0) {
			pool[i].tag = tagVoid
			pool[i].W = 0
			pool[i].H = 0
		}
	}

	// Inherit stretch.
	for i := 0; i < len(pool)-1; i++ {
		s := &pool[i]
		// Stop the chain of inheritance if we have a Limit.
		var stopw, stoph bool
		for j := range s.pres(wo) {
			k := &s.pres(wo)[j]
			if k.tag == tagLimit {
				if k.W != 0 {
					stopw = true
				}
				if k.H != 0 {
					stoph = true
				}
			}
		}
		for j := range s.kids(wo) {
			k := &s.kids(wo)[j]
			if k.W < 0 && !stopw {
				// s.W += k.W
			}
			if k.H < 0 && !stoph {
				// s.H += k.H
			}
		}
	}

	if wo.Events.Match(`Press(F4)`) {
		println("@ Tree before applying stretches")
		sormp(wo, *last(pool), 0)
		println()
	}

	wo.apply(nil, last(pool))

	// Inherit moves and paints.
	for i := len(pool) - 1; i >= 0; i-- {
		s := &pool[i]
		for j := range s.kids(wo) {
			kid := &s.kids(wo)[j]
			kid.x += s.x
			kid.y += s.y
			if kid.fill == (nanovgo.Paint{}) {
				kid.fill = s.fill
			}
			if kid.stroke == (nanovgo.Paint{}) {
				kid.stroke = s.stroke
			}
			if kid.flags&flagSetStrokewidth == 0 {
				kid.strokew = s.strokew
			}
		}
	}

	// Print tree for debug. Do it before sorting.
	if wo.Events.Match(`Press(F1)`) {
		println("@ Tree after layout")
		sormp(wo, *last(pool), 0)
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
	for _, s := range pool {
		if s.tag > 0 {
			// Positioning bodges
			switch s.tag {
			default:
				s = s.decimate()
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
	for i := range wo.pool {
		wo.pool[i] = Sorm{}
	}
	wo.pool = wo.pool[:0]
	for i := range wo.tmp {
		wo.tmp[i] = Sorm{}
	}
	wo.tmp = wo.tmp[:0]
	wo.nextn = 0

	wo.sinks = wo.sinks[:1]

	// z := wo.Trace[0]
	// if wo.MatchInNochoke(`Click(1):in`, geom.Rect(-11111, -11111, 11111, 11111)) {
	// 	wo.Trace[0] = z
	// 	println(collect(wo.Events.Trace, func(p EventPoint) int { return p.z }))
	// }

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
	for _, s := range s.kids(wo) {
		sormp(wo, s, tab+1)
	}
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
	window
}

func (wi *Window) Rect() image.Rectangle {
	return wi.rect()
}

func New(config Config) (wo World) {
	if config.WindowRect.Dx() == 0 {
		config.WindowRect.Max.X = config.WindowRect.Min.X + 1024
	}
	if config.WindowRect.Dy() == 0 {
		config.WindowRect.Max.Y = config.WindowRect.Min.Y + 768
	}

	runtime.LockOSThread()
	concretenew(config, &wo)

	wo.Events = NewEventTracer(wo.Window.window, config.ReplayReader)
	wo.sinks = make([]func(any), 1)
	wo.keys = map[any]*labelt{}
	return wo
}
