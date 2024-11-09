// Contraption: a simple framework for user interfaces.
//
// A good user interface framework must be an engine for a word processing game.
//
// TODO:
//	- Investigate better ways to specify size values
//		- Current: -1+20i — complex128, + concrete, - stretch, imaginary adds concrete change to stretchy
//			+ Almost transparent
//			+ Looks very good for static values
//			~ Simple
//			- Variable stretch is done through complex() constructor
//			- Arithmetic is static only
//		- Strings: "200", "1s+20px"
//			+ Most powerful
//			+ Can implement deferred arithmetics
//				+ 1dp+8cm-20%/1s
//			+ Can be extendable by user
//			+ Parsing can be quite fast if done right
//			- Construction from variables requires strconv
//		- [N]int
//			+ Easily extendable
//			- Looks gross and is not ergonomical
//		- Builders: St(1).Add(Px(20))
//			+ Powerful
//			+ Can implement deferred arithmetics
//				+ Dp(1).Add(Cm(8)).Sub(Percent(20)).Div(St(1))
//			- Less readable
//		- Add another two values to every Sorm constructor
//			- Requires major refactoring
//			- Looks overcomplicated
//	- Sequence buffer size hints
//		- 16 elements on complex and large items
//		- 128 on moderately complex and medium-sized items (messages in chat)
//		- 1024 on small and simple elements (rows of data, histograms)
//	- Round shape sizes inside aligners.
//		- Noround modifier
//		? Smooth rounding hint — don't round if in motion
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
//		+ wo.Sequence(seq Sequence) — a window to infinity!
//			+ Every returned Sorm is included to the parent
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
//	? Modifier to shape position independence
//	? Fix paint interface
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
	"reflect"
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
	flagIteratedNotScissor
	flagIteratedScissor
	flagDontDecimate
	flagDecimate
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
	tagIllustration
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
	tagPosttransform
	tagTransform
	tagScissor
	tagHshrink
	tagVshrink
	tagLimit // Limit must be executed after update of a matrix, but before aligners, because it influences size.
	tagVfollow
	tagHfollow
	tagDontDecimate
	tagDecimate
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
	shapeActions[tagIllustration] = illustrationrun

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

	preActions[-100-tagPosttransform] = posttransformrun
	preActions[-100-tagTransform] = transformrun
	preActions[-100-tagScissor] = scissorrun
	preActions[-100-tagVfollow] = vfollowrun
	preActions[-100-tagHfollow] = hfollowrun
	preActions[-100-tagHshrink] = hshrinkrun
	preActions[-100-tagVshrink] = vshrinkrun
	preActions[-100-tagLimit] = limitrun
	preActions[-100-tagDontDecimate] = dontdecimaterun
	preActions[-100-tagDecimate] = decimaterun

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
	gen int

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

	drag        any
	dragstart   geom.Point
	sinks       []func(any)
	dragEffects map[reflect.Type]func(interval [2]geom.Point, drag any) Sorm

	showOutlines bool
	f1           bool

	alloc func(n int) (left, right int)

	scissored  []*Sorm
	scissoring bool

	images map[io.Reader]imagestruct
}

type imagestruct struct {
	gen     int
	texid   int
	origsiz geom.Point
}

func AddDragEffect[T any](wo *World, convert func(interval [2]geom.Point, drag T) Sorm) {
	var z T
	wo.dragEffects[typeof(z)] = func(interval [2]geom.Point, drag any) Sorm {
		return convert(interval, drag.(T))
	}
}

func (wo *World) ResetDragEffects() {
	for k := range wo.dragEffects {
		delete(wo.dragEffects, k)
	}
}

type Sorm struct {
	z, z2, i, pi int
	tag          tagkind
	flags        flagval

	Size, p, l            point
	knowns, props, eprops point
	add, ialign           point

	r       float64
	m, prem geom.Geom
	aligner alignerkind
	kidsl, kidsr,
	modsl, modsr,
	presl, presr int

	index *Index

	scissori int
	scissor  geom.Rectangle

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
		vals = fmt.Sprint(f(s.Size.X), ", ", f(s.Size.Y), ", ", f(s.p.X), ", ", f(s.p.Y), ", _", f(s.l.X), ", _", f(s.l.Y))
	} else {
		vals = fmt.Sprint(f(s.Size.X), "×", f(s.Size.Y), ", ", f(s.p.X), "y", f(s.p.Y), ", ↓", f(s.l.X), "×", f(s.l.Y), " ", color(p.innerColor))
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

	return fmt.Sprint(z, strings.Repeat(" ", max(0, 10-digits)), seq, ovrx, ovry, btw, " ", s.tag.String(), " ", vals, ` `, s.props, ` `, scissor, key) //, " ", s.callerfile, ":", s.callerline)
}

func (s Sorm) decimate() Sorm {
	// FIXME Decimation is now done inside layout.
	if s.flags&flagDontDecimate > 0 {
		return s
	}
	s.p.X = math.Floor(s.p.X)
	s.p.Y = math.Floor(s.p.Y)
	s.Size.X = math.Ceil(s.Size.X)
	s.Size.Y = math.Ceil(s.Size.Y)
	return s
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

func (s Sorm) Resize(x, y float64) Sorm {
	s.Size.X = x
	s.Size.Y = y
	return s
}

func (s Sorm) Rectangle() geom.Rectangle {
	return geom.Rect(s.p.X, s.p.Y, s.p.X+s.Size.X, s.p.Y+s.Size.Y)
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
		s.l.X = wo.Wwin
		s.l.Y = wo.Hwin
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
			if wo.dragEffects[typeof(wo.drag)] != nil && wo.drag != nil {
				ps := [2]geom.Point{wo.dragstart, wo.Trace[0].Pt}
				return wo.dragEffects[typeof(wo.drag)](ps, wo.drag)
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
	wo.Vgo.SetFontSize(11)

	return true
}

func (wo *World) beginvirtual() (pool []Sorm) {
	if sameslice(wo.pool, wo.auxpool) {
		panic(`contraption: nested Sequences are not allowed`)
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

type kiargs struct {
	i              Index
	a              alignerkind
	notscissormark bool
}

func (s *Sorm) kidsiter(wo *World, a kiargs, f func(k *Sorm)) {
	// TODO Idea: take a limit in kidsiter, and if it is scissored, stop iteration when over it.
	// It should work since scissored compounds can't stretch kids.
	if s.tag == tagSequence {
		return
	}

	// mode := func(k *Sorm) bool {
	// 	if !wo.scissoring && k.flags&flagIteratedScissor == 0 {
	// 		k.flags |= flagIteratedNotScissor
	// 		return true
	// 	}
	// 	if wo.scissoring && k.flags&flagIteratedNotScissor == 0 {
	// 		k.flags |= flagIteratedScissor
	// 		return true
	// 	}
	// 	return false
	// }

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
				q.Get(wo, i, args[:])
				// TODO Call (*World).layout on every scissored compound right here.
				// TODO
				// var buf [1024]Sorm
				// for i := 0; i < q.Length(wo); i++ {
				//
				// }
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
				if a.notscissormark { // (3)
					k.flags |= flagIteratedNotScissor
				}
				if k.flags&flagBreakIteration > 0 { // (2)
					break out
				}
			}
		} else {
			f(k)                  // (1)
			if a.notscissormark { // (3)
				k.flags |= flagIteratedNotScissor
			}
			if k.flags&flagBreakIteration > 0 { // (2)
				break out
			}
		}
	}
}

func (wo *World) topbreadthiter(pool []Sorm, f func(s, _ *Sorm)) {
	for i := len(pool) - 1; i >= 0; i-- {
		s := &pool[i]
		if s.tag == tagSequence {
			p := wo.beginvirtual()
			pool := wo.auxpool[s.kidsl:s.kidsr]
			for i := len(pool) - 1; i >= 0; i-- {
				s := &pool[i]
				f(s, s)
			}
			wo.endvirtual(p)
		} else {
			f(s, s)
		}
	}
}

func (wo *World) bottombreadthiter(pool []Sorm, f func(s, _ *Sorm)) {
	for i := range pool {
		s := &pool[i]
		if s.tag == tagSequence {
			p := wo.beginvirtual()
			pool := wo.auxpool[s.kidsl:s.kidsr]
			for i := range pool {
				s := &pool[i]
				f(s, s)
			}
			wo.endvirtual(p)
		} else {
			f(s, s)
		}
	}
}

// If a Compound has no aligner set (“stack” layout)
func noaligner(wo *World, c *Sorm) {
	minp := point{}
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		c.Size.X = max(c.Size.X, k.Size.X)
		c.Size.Y = max(c.Size.Y, k.Size.Y)
		if k.Size.X < 0 {
			minp.X = min(minp.X, k.Size.X)
		}
		if k.Size.Y < 0 {
			minp.Y = min(minp.Y, k.Size.Y)
		}
	})
	if c.flags&flagHshrink > 0 {
		c.l.X = c.Size.X
	}
	if c.flags&flagVshrink > 0 {
		c.l.Y = c.Size.Y
	}
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		stretch := false
		if k.Size.X < 0 {
			k.Size.X = c.l.X * k.Size.X / minp.X
			k.Size.X += k.add.X
			stretch = true
		}
		if k.Size.Y < 0 {
			k.Size.Y = c.l.Y * k.Size.Y / minp.Y
			k.Size.Y += k.add.Y
			stretch = true
		}
		if stretch {
			k.l.Y = k.Size.Y
			k.l.X = k.Size.X
			wo.apply(c, k)
		}
	})
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		c.Size.X = max(c.Size.X, k.Size.X)
		c.Size.Y = max(c.Size.Y, k.Size.Y)
	})
}

func vfollowaligner(wo *World, c *Sorm) {
	sequencealigner(wo, c, false)
}

func hfollowaligner(wo *World, c *Sorm) {
	sequencealigner(wo, c, true)
}

func swapxy(p *point) {
	p.X, p.Y = p.Y, p.X
}

func axis(h bool) (begin, end func(k *Sorm)) {
	swap := func(k *Sorm) {
		if h {
			swapxy(&k.Size)
			swapxy(&k.p)
			swapxy(&k.l)
			swapxy(&k.props)
			swapxy(&k.eprops)
			swapxy(&k.knowns)

			swapxy(&k.scissor.Min)
			swapxy(&k.scissor.Max)
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
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		beginaxis(k)
		if k.tag == 0 {
			c.props.X = max(c.props.X, k.eprops.X)
			c.props.Y += k.eprops.Y
		} else {
			k.eprops = geom.Pt(-min(0, k.Size.X), -min(0, k.Size.Y))
			if k.Size.X >= 0 {
				c.knowns.X = max(c.knowns.X, k.Size.X)
			} else {
				c.props.X = max(c.props.X, -k.Size.X)
			}
			if k.Size.Y >= 0 {
				c.knowns.Y += k.Size.Y
			} else {
				c.props.Y += -k.Size.Y
			}
		}
		endaxis(k)
	})
	// Don't set if overriden by Limit with negative size.
	c.eprops = trypos(c.eprops, c.props)
	endaxis(c)
}

func sequencealigner(wo *World, c *Sorm, h bool) {
	c.Size.X, c.Size.Y = 0, 0
	beginaxis, endaxis := axis(h)

	// Calculate unknowns and apply kids which sizes were unknown,
	// lay out the sequence then.
	y := 0.0
	beginaxis(c)
	e := .0
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		stretch := false

		beginaxis(k)
		// Y is the main axis, X is perpendicular.
		if k.eprops.Y > 0 {
			// max() means we don't stretch if we're out of limit.
			k.Size.Y = max(0, (c.l.Y-c.knowns.Y)/c.props.Y*k.eprops.Y)
			// Round lengths and sizes to the nearest integer.
			// TODO Accumulate error in kids and try to make it strobe less.
			if c.flags&flagDecimate > 0 || !(c.flags&flagDontDecimate > 0) {
				e += k.Size.Y - math.Round(k.Size.Y)
				k.Size.Y = math.Round(k.Size.Y)
			}
			k.l.Y = k.Size.Y
			stretch = true
		}
		if k.eprops.X > 0 {
			k.Size.X = c.l.X / c.props.X * k.eprops.X
			k.l.X = k.Size.X
			stretch = true
		}
		endaxis(k)

		if stretch {
			wo.apply(c, k)
		}

		beginaxis(k)
		k.p.Y = y
		y += k.Size.Y
		c.Size.X = max(c.Size.X, k.Size.X)
		endaxis(k)
	})
	c.Size.Y = y + e
	if c.flags&flagDecimate > 0 || !(c.flags&flagDontDecimate > 0) {
		c.Size.Y = math.Round(c.Size.Y)
	}
	endaxis(c)
}

func (wo *World) apply(p *Sorm, c *Sorm) {
	// Apply is called on every shape, so ignore anything that is not Compound.
	if c.tag != 0 {
		return
	}

	// apply presumes that premodifiers are already sorted.
	for _, m := range c.pres(wo) {
		if m.tag == tagLimit && (m.Size.X >= 0 || m.Size.Y >= 0) {
			preActions[-100-m.tag](wo, c, &m)
		}
	}

	// Set scissor to limit if needed.
	if c.flags&flagScissor > 0 {
		c.scissor = geom.Rect(0, 0, c.l.X, c.l.Y)
	}
	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		// NOTE Aligner is called after these assignments.
		// 	So this can't influence limits at later stages.
		k.l.X = c.l.X
		k.l.Y = c.l.Y
		// Inherit scissor.
		if c.flags&flagScissor > 0 {
			k.scissori = c.i
		} else {
			k.scissori = c.scissori
		}

		// Apply scale.
		ns := k.m.ApplyPt(geom.Pt(k.Size.X, k.Size.Y))
		ims := k.m.ApplyPt(geom.Pt(k.add.X, k.add.Y))
		// Don't scale stretch coefficients.
		k.Size.X = cond(k.Size.X >= 0, ns.X, k.Size.X)
		k.Size.Y = cond(k.Size.Y >= 0, ns.Y, k.Size.Y)
		// But scale imaginaries.
		k.add.X = ims.X
		k.add.Y = ims.Y

		// Process only kids which sizes are known first.
		if k.Size.X >= 0 && k.Size.Y >= 0 {
			wo.apply(c, k)
		}
	})

	alignerActions[c.aligner](wo, c)

	// Apply postorder/anyorder modifiers.
	for _, m := range c.mods(wo) {
		modActions[-m.tag](wo, c, &m)
	}

	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		// Apply size override to c if it has one.
		if k.flags&flagOvrx > 0 {
			c.Size.X = k.Size.X
			c.p.X = min(c.p.X, -k.p.X)
		}
		if k.flags&flagOvry > 0 {
			c.Size.Y = k.Size.Y
			c.p.Y = min(c.p.Y, -k.p.Y)
		}
	})

	if c.flags&flagScissor > 0 {
		c.Size.X = min(c.Size.X, c.l.X)
		c.Size.Y = min(c.Size.Y, c.l.Y)
	}
}

func (wo *World) resolvepremods(_ *Sorm, c *Sorm) {
	if c.tag != 0 {
		return
	}
	// Premodifiers have order of execution.
	slices.SortFunc(c.pres(wo), func(a, b Sorm) int {
		return int(b.tag - a.tag)
	})

	for _, m := range c.pres(wo) {
		if m.tag == tagLimit && (m.Size.X >= 0 || m.Size.Y >= 0) {
			continue
		}
		preActions[-100-m.tag](wo, c, &m)
	}

	// Defer scissored compounds for another pass.
	// if c.flags&flagScissor > 0 && !wo.scissoring {
	// 	wo.scissored = append(wo.scissored, c)
	// 	return
	// }

	c.kidsiter(wo, kiargs{}, func(k *Sorm) {
		// Cascade matrices and some flags
		k.m = c.m
		k.flags |= c.flags & flagDontDecimate
		k.flags |= c.flags & flagDecimate
		wo.resolvepremods(c, k)
		// Resolve sprite text widths based off a real font size
		// TODO Broken, probably because of incorrect scaling with matrices
		// TODO Vertical
		// TODO flagNonlinear — set the size of an element only after setting up the matrix.
		//	Another element with this property is Equation.
		if k.tag == tagText {
			s := k
			wo.Vgo.SetFontFaceID(s.fontid)
			wo.Vgo.SetFontSize(float32(k.Size.Y))
			_, abcd := wo.Vgo.TextBounds(0, 0, s.key.(string))
			_, space := wo.Vgo.TextBounds(0, 0, " ")
			s.Size.X = float64(abcd[2]-abcd[0]) - float64(space[2]-space[0])
			if s.Size.X < 0 {
				s.Size.X = 0
			}
		}
	})
}

func (wo *World) layout(pool []Sorm, root *Sorm) {
	// Resolve premodifiers and stack negative sizes.
	wo.resolvepremods(nil, root)
	wo.bottombreadthiter(pool, func(s, _ *Sorm) {
		switch s.aligner {
		case alignerVfollow:
			sequencedivider(wo, s, false)
		case alignerHfollow:
			sequencedivider(wo, s, true)
		}
	})

	if wo.Events.Match(`Press(F4)`) {
		println("@ Tree before applying stretches")
		sormp(wo, *last(pool), 0)
		println()
	}

	// Do the layout.
	wo.apply(nil, root)

	// Inherit moves. Cascade scissors, fills and strokes.
	// TODO Maybe wo.apply should do it? Kind of makes more sense.
	// The reason it is separated is because we don't know how absolute component
	// sizes and coordinates till the very end.
	wo.topbreadthiter(pool, func(c, efc *Sorm) {
		c.kidsiter(wo, kiargs{}, func(k *Sorm) {
			k.p.X += efc.p.X
			k.p.Y += efc.p.Y
			k.scissor = k.scissor.Add(geom.Pt(efc.p.X, efc.p.Y))
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
	})

	// Apply scissors.
	wo.bottombreadthiter(pool, func(s, _ *Sorm) {
		if s.scissori > 0 {
			s.scissor = pool[s.scissori].Rectangle()
		}
	})

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
}

// Develop applies the layout and renders the next frame.
// See package description for preferred use of Contraption.
func (wo *World) Develop() {
	vgo := wo.Vgo

	vgo.ResetTransform()
	root := &wo.pool[len(wo.pool)-1]
	root.l.X = wo.Wwin
	root.l.Y = wo.Hwin

	pool := wo.pool[0:wo.rend]

	for i := len(pool) - 1; i >= 0; i-- {
		s := &pool[i]
		s.i = i
		s.kidsiter(wo, kiargs{}, func(k *Sorm) {
			k.pi = s.i
		})
	}

	wo.layout(pool, last(pool))
	// for i := range wo.scissored {
	// 	wo.layout(pool, wo.scissored[i])
	// }

	// Sort in draw order.
	slices.SortFunc(pool, func(a, b Sorm) int {
		return a.z - b.z
	})
	// After this point, (*Sorm).kidsiter won't work because indices are broken.

	// Apply conditional paints, match drag-and-drop events.
	for i := len(pool) - 1; i >= 0; i-- {
		s := &pool[i]
		r := geom.Rect(s.p.X, s.p.Y, s.p.X+s.Size.X, s.p.Y+s.Size.Y)
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
	// FIXME auxpool is not sorted.
	wo.bottombreadthiter(pool, func(s, _ *Sorm) {
		if s.tag <= 0 {
			return
		}

		if s.scissori > 0 && s.scissor.Intersect(s.Rectangle()) == (geom.Rectangle{}) {
			return
		}

		*s = s.decimate()
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
			s.p.X -= 0.5
			s.p.Y -= 0.5
		case tagTopDownText:
			fallthrough
		case tagBottomUpText:
			fallthrough
		case tagText:
			// s.pos.X -= 1
			// s.pos.Y -= 1

		case tagVectorText:
		case tagCanvas:
		}
		shapeActions[s.tag](wo, s)
		vgo.ResetTransform()
	})

	wo.Vgo.Reset()

	if wo.Match(`Press(F2)`) {
		wo.showOutlines = !wo.showOutlines
	}
	if wo.showOutlines {
		vgo.SetStrokeWidth(1)
		vgo.SetStrokePaint(hexpaint(`#00000020`))
		vgo.BeginPath()
		for _, s := range pool {
			if s.Size.Y < .5 && s.Size.X < .5 {
				continue
			}
			s := s.decimate()
			s.p.X -= 0.5
			s.p.Y -= 0.5
			if s.tag >= 0 {
				wo.Vgo.Rect(float32(s.p.X), float32(s.p.Y), float32(s.Size.X), float32(s.Size.Y))
			}
		}
		for _, s := range wo.auxpool {
			if s.Size.Y < .5 && s.Size.X < .5 {
				continue
			}
			s := s.decimate()
			s.p.X -= 0.5
			s.p.Y -= 0.5
			if s.tag >= 0 {
				wo.Vgo.Rect(float32(s.p.X), float32(s.p.Y), float32(s.Size.X), float32(s.Size.Y))
			}
		}
		vgo.ClosePath()
		vgo.Stroke()

		vgo.SetStrokeWidth(0)
		vgo.SetFillPaint(hexpaint(`#ff000020`))
		vgo.BeginPath()
		for _, s := range pool {
			if s.scissori == 0 {
				continue
			}
			if s.Size.Y < .5 && s.Size.X < .5 {
				continue
			}
			s := s.decimate()
			s.p.X -= 0.5
			s.p.Y -= 0.5
			if s.tag >= 0 {
				wo.Vgo.Rect(float32(s.scissor.Min.X), float32(s.scissor.Min.Y),
					float32(s.scissor.Dx()), float32(s.scissor.Dy()))
			}
		}
		for _, s := range wo.auxpool {
			if s.scissori == 0 {
				continue
			}
			if s.Size.Y < .5 && s.Size.X < .5 {
				continue
			}
			s := s.decimate()
			s.p.X -= 0.5
			s.p.Y -= 0.5
			if s.tag >= 0 {
				wo.Vgo.Rect(float32(s.scissor.Min.X), float32(s.scissor.Min.Y),
					float32(s.scissor.Dx()), float32(s.scissor.Dy()))
			}
		}
		vgo.ClosePath()
		vgo.Fill()
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

	zeroandclear(&wo.pool)
	wo.nextn = 0
	zeroandclear(&wo.auxpool)
	wo.auxn = 0
	zeroandclear(&wo.tmp)
	zeroandclear(&wo.scissored)

	wo.sinks = wo.sinks[:1]

	wo.Events.develop()
	wo.windowDevelop()

	for k, v := range wo.keys {
		v.counter--
		if v.counter <= 0 {
			delete(wo.keys, k)
		}
	}

	for k, v := range wo.images {
		if v.gen < wo.gen-2 {
			wo.Vgo.DeleteImage(v.texid)
			delete(wo.images, k)
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
	s.kidsiter(wo, kiargs{}, func(k *Sorm) {
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
	wo.images = map[io.Reader]imagestruct{}
	wo.dragEffects = map[reflect.Type]func(interval [2]geom.Point, drag any) Sorm{}
	wo.alloc = wo.allocmain
	return wo
}
