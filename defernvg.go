package contraption

import (
	"image"
	"reflect"
	"time"

	"github.com/neputevshina/contraption/nanovgo"
	"github.com/neputevshina/contraption/nanovgo/fontstashmini"
	"github.com/neputevshina/contraption/op"
	"github.com/neputevshina/geom"
)

type RenderOp struct {
	Tag   op.Op
	Args  [10]float64
	Iargs [10]int

	Fontsiz float64
	Hfont   int
	Himage  int

	Lsp     float64
	Fblur   float64
	Strokep nanovgo.Paint
	Strokec nanovgo.Color
	Fillp   nanovgo.Paint
	Fillc   nanovgo.Color
	Strokew float64

	Linecap nanovgo.LineCap
	nanovgo.Direction
	nanovgo.Winding
	nanovgo.Align
	nanovgo.TransformMatrix
	nanovgo.GlyphPosition

	Runes []rune

	Left, Right int
}

type RenderImage struct {
	Deleted bool
	H       int

	image.Image
	Data []byte
	nanovgo.ImageFlags
	Wh image.Point
}

type RenderFont struct {
	Deleted  bool
	Hbackend int

	Data     []byte
	FreeData byte
}

// Context is a draw list with the interface of Nanovgo.
//
// It is interpreted by the backend to draw on screen or to do anything else.
// TODO All transformations must be resolved by Context, not by the backend.
type Context struct {
	publicContext
	state op.Op

	parent *Context
}

// Sub returns a persistent context that could be replayed later.
func (c *Context) Sub() *Context {
	// NOTE Subcontexts are not hashed.
	return &Context{
		state:  c.state, // TODO Subcontexts are not required to Begin/EndFrame
		parent: c,
		publicContext: publicContext{
			Fs: c.Fs,
		}}
}

// Replay replays the subordinating context of the current.
func (c *Context) Replay(sub *Context) {
	c.Log = append(c.Log, sub.Log...)
	// TODO Use opReplay and replay at the interpretation time.
}

func newContext() *Context {
	return &Context{
		publicContext: publicContext{
			Fs:     fontstashmini.New(512, 512),
			Images: []RenderImage{{}},
		},
	}
}

type SpriteUnit struct {
	Hfont int
	Clip  geom.Rectangle
	Tc    geom.Rectangle
}

// PublicContext is a draw list, publicly exposed for any (possibly custom) backend.
//
// TODO The only way to get PublicContext is to use the Process method on Context.
// type PublicContext publicContext

type publicContext struct {
	RenderOp // Current state

	Fs *fontstashmini.FontStash

	Log           []RenderOp
	Images        []RenderImage
	Fonts         []RenderFont
	SpriteUnits   []SpriteUnit
	devicePxRatio float64
}

func (c *Context) assertPathStarted() {
	if c.state != op.BeginPath {
		panic(`contraption.Context: vector operation before BeginPath`)
	}
}

func (c *Context) assertFrameStarted() {
	if c.parent != nil {
		return
	}
	if c.state != op.BeginFrame && c.state != 1 {
		panic(`contraption.Context: caller was called before BeginFrame`)
	}
}

func (c *Context) Block(block func()) {
	c.Save()
	block()
	c.Restore()
}

func (c *Context) DebugDumpPathCache() {
	panic(`unimplemented`)
}

func (c *Context) IntersectScissor(x, y, w, h float64) {
	c.assertPathStarted()
	_ = c.add(op.IntersectScissor, RenderOp{
		Args: [10]float64{x, y, w, h},
	})
}

func functag(f any) uintptr {
	return uintptr(reflect.ValueOf(f).UnsafePointer())
}

func (c *Context) add(f op.Op, op RenderOp) op.Op {
	op.Tag = f
	c.Log = append(c.Log, op)

	p := c
	for ; p.parent != nil; p = p.parent {
	}
	return op.Tag
}

/* Frame and context state */

func (c *Context) BeginFrame(windowWidth, windowHeight int, devicePixelRatio float64) {
	if c.state != 0 {
		panic(`contraption.Context: BeginFrame can only be called at the start of a frame`)
	}
	c.SpriteUnits = c.SpriteUnits[:0]

	st := c.add(op.BeginFrame, RenderOp{
		Iargs: [10]int{windowWidth, windowHeight},
		Args:  [10]float64{devicePixelRatio},
	})
	c.devicePxRatio = devicePixelRatio
	c.state = st
}
func (c *Context) EndFrame() (oldhash [512]byte) {
	_ = c.add(op.EndFrame, RenderOp{})
	c.state = 0
	return
}
func (c *Context) CancelFrame() {
	panic(`unimplemented`)
}
func (c *Context) Delete() {
	panic(`unimplemented`)
}

/* Images */

func (c *Context) CreateImageRGBA(w, h int, imageFlags nanovgo.ImageFlags, data []byte) int {
	if c.parent != nil {
		panic(`contraption.Context: can't create image from subcontext`)
	}
	c.Images = append(c.Images, RenderImage{
		Image:      nil,
		Data:       data,
		ImageFlags: imageFlags,
		Wh:         image.Pt(w, h),
	})
	_ = c.add(op.CreateImageRGBA, RenderOp{
		Himage: len(c.Images),
	})
	return len(c.Images)
}
func (c *Context) CreateImageFromGoImage(imageFlag nanovgo.ImageFlags, img image.Image) int {
	if c.parent != nil {
		panic(`contraption.Context: can't create image from subcontext`)
	}
	c.Images = append(c.Images, RenderImage{
		Image:      img,
		ImageFlags: imageFlag,
	})
	_ = c.add(op.CreateImageFromGoImage, RenderOp{
		Himage: len(c.Images) - 1,
	})
	return len(c.Images) - 1
}
func (c *Context) UpdateImage(img int, data []byte) error {
	panic(`unimplemented`)
}
func (c *Context) DeleteImage(img int) {
	if img < 0 {
		panic(`contraption.Context: incorrect image handle`)
	}
	if c.Images[img].Deleted {
		panic(`contraption.Context: double image delete`)
	}
	c.Images[img] = RenderImage{Deleted: true}
}
func (c *Context) ImageSize(img int) (int, int, error) {
	panic(`unimplemented`)
}

/* Fonts */

func (c *Context) CreateFontFromMemory(name string, data []byte, freeData uint8) int {
	if c.parent != nil {
		panic(`contraption.Context: can't create font from subcontext`)
	}
	c.Fonts = append(c.Fonts, RenderFont{
		Data:     data,
		FreeData: freeData,
	})
	c.Fs.AddFontFromMemory(name, data, freeData)
	return len(c.Fonts) - 1
}

/* Shapes */

func (c *Context) Circle(cx, cy, r float64) {
	c.assertPathStarted()
	_ = c.add(op.Circle, RenderOp{
		Args: [10]float64{cx, cy, r},
	})
}
func (c *Context) Rect(x, y, w, h float64) {
	c.assertPathStarted()
	_ = c.add(op.Rect, RenderOp{
		Args: [10]float64{x, y, w, h},
	})
}
func (c *Context) Ellipse(cx, cy, rx, ry float64) {
	c.assertPathStarted()
	_ = c.add(op.Ellipse, RenderOp{
		Args: [10]float64{cx, cy, rx, ry},
	})
}
func (c *Context) RoundedRect(x, y, w, h, r float64) {
	c.assertPathStarted()
	_ = c.add(op.RoundedRect, RenderOp{
		Args: [10]float64{x, y, w, h, r},
	})
}

/* Paths */

func (c *Context) BeginPath() {
	c.assertFrameStarted()
	// if c.state != functag(opClosePath) {
	// 	panic(`contraption.Context: can't BeginPath before ClosePath`)
	// }
	st := c.add(op.BeginPath, RenderOp{})
	c.state = st
}
func (c *Context) ClosePath() {
	c.assertPathStarted()
	_ = c.add(op.ClosePath, RenderOp{})
	c.state = op.BeginFrame
}
func (c *Context) Fill() {
	if c.state == 1 {
		panic(`contraption.Context: another Fill can be called only after Also`)
	}
	_ = c.add(op.Fill, RenderOp{})
	c.state = 1
}
func (c *Context) Stroke() {
	if c.state == 1 {
		panic(`contraption.Context: another Stroke can be called only after Also`)
	}
	_ = c.add(op.Stroke, RenderOp{})
	c.state = 1
}
func (c *Context) Also() {
	if c.state != 1 {
		panic(`contraption.Context: can't use Also before Stroke or Fill`)
	}
	_ = c.add(op.Also, RenderOp{})
	c.state = 2
}
func (c *Context) Arc(cx, cy, r, a0, a1 float64, dir nanovgo.Direction) {
	c.assertPathStarted()
	_ = c.add(op.Arc, RenderOp{
		Args:      [10]float64{cx, cy, r, a0, a1},
		Direction: dir,
	})
}
func (c *Context) ArcTo(x1, y1, x2, y2, radius float64) {
	c.assertPathStarted()
	_ = c.add(op.ArcTo, RenderOp{
		Args: [10]float64{x1, y1, x2, y2, radius},
	})
}
func (c *Context) BezierTo(c1x, c1y, c2x, c2y, x, y float64) {
	c.assertPathStarted()
	_ = c.add(op.BezierTo, RenderOp{
		Args: [10]float64{c1x, c1y, c2x, c2y, x, y},
	})
}
func (c *Context) LineTo(x, y float64) {
	c.assertPathStarted()
	_ = c.add(op.LineTo, RenderOp{
		Args: [10]float64{x, y},
	})
}
func (c *Context) MoveTo(x, y float64) {
	c.assertPathStarted()
	_ = c.add(op.MoveTo, RenderOp{
		Args: [10]float64{x, y},
	})
}
func (c *Context) QuadTo(cx, cy, x, y float64) {
	c.assertPathStarted()
	_ = c.add(op.QuadTo, RenderOp{
		Args: [10]float64{cx, cy, x, y},
	})
}
func (c *Context) PathWinding(winding nanovgo.Winding) {
	c.assertPathStarted()
	_ = c.add(op.PathWinding, RenderOp{
		Winding: winding,
	})
}

/* State management */

func (c *Context) Reset() {
	_ = c.add(op.Reset, RenderOp{})
}
func (c *Context) ResetScissor() {
	_ = c.add(op.ResetScissor, RenderOp{})
}
func (c *Context) ResetTransform() {
	_ = c.add(op.ResetTransform, RenderOp{})
}
func (c *Context) Restore() {
	_ = c.add(op.Restore, RenderOp{})
}
func (c *Context) Save() {
	_ = c.add(op.Save, RenderOp{})
}

// TODO All transformations must be resolved by Context, not by the backend.
/* Transformation mutators */

func (c *Context) Rotate(angle float64) {
	_ = c.add(op.Rotate, RenderOp{})
}
func (c *Context) Scale(x, y float64) {
	panic(`unimplemented`)
}
func (c *Context) Scissor(x, y, w, h float64) {
	c.assertFrameStarted()
	_ = c.add(op.Scissor, RenderOp{
		Args: [10]float64{x, y, w, h},
	})
}
func (c *Context) SkewX(angle float64) {
	panic(`unimplemented`)
}
func (c *Context) SkewY(angle float64) {
	panic(`unimplemented`)
}
func (c *Context) SetTransform(t nanovgo.TransformMatrix) {
	c.TransformMatrix = t
	_ = c.add(op.SetTransform, RenderOp{
		TransformMatrix: t,
	})
}
func (cx *Context) SetTransformByValue(a, b, c, d, e, f float64) {
	panic(`unimplemented`)
}
func (c *Context) Translate(x, y float64) {
	panic(`unimplemented`)
}

/* Miscellaneous mutators */

func (c *Context) SetFillColor(color nanovgo.Color) {
	c.Fillc = color
	_ = c.add(op.SetFillColor, RenderOp{
		Fillc: color,
	})
}
func (c *Context) SetFillPaint(paint nanovgo.Paint) {
	c.Fillp = paint
	_ = c.add(op.SetFillPaint, RenderOp{
		Fillp: paint,
	})
}
func (c *Context) SetFontBlur(blur float64) {
	panic(`unimplemented`)
}
func (c *Context) SetFontFace(font string) {
	panic(`unimplemented`)
}
func (c *Context) SetFontFaceID(font int) {
	c.Hfont = font
	_ = c.add(op.SetFontFaceID, RenderOp{
		Hfont: font,
	})
}
func (c *Context) SetFontSize(size float64) {
	// TODO Convert from cap to em in other backends
	c.Fontsiz = size
	_ = c.add(op.SetFontSize, RenderOp{
		Fontsiz: size,
	})
}
func (c *Context) SetGlobalAlpha(alpha float64) {
	panic(`unimplemented`)
}
func (c *Context) SetLineCap(cap nanovgo.LineCap) {
	panic(`unimplemented`)
}
func (c *Context) SetLineJoin(joint nanovgo.LineCap) {
	panic(`unimplemented`)
}
func (c *Context) SetMiterLimit(limit float64) {
	panic(`unimplemented`)
}
func (c *Context) SetStrokeColor(color nanovgo.Color) {
	c.Strokec = color
	_ = c.add(op.SetStrokeColor, RenderOp{
		Strokec: color,
	})
}
func (c *Context) SetStrokePaint(paint nanovgo.Paint) {
	c.Strokep = paint
	_ = c.add(op.SetStrokePaint, RenderOp{
		Strokep: paint,
	})
}
func (c *Context) SetStrokeWidth(width float64) {
	c.Strokew = width
	_ = c.add(op.SetStrokeWidth, RenderOp{
		Strokew: width,
	})
}
func (c *Context) SetTextAlign(align nanovgo.Align) {
	c.Align = align
	_ = c.add(op.SetTextAlign, RenderOp{
		Align: align,
	})
}
func (c *Context) SetTextLetterSpacing(spacing float64) {
	panic(`unimplemented`)
}
func (c *Context) SetTextLineHeight(lineHeight float64) {
	panic(`unimplemented`)
}

/* Online getters */

func (c *Context) FindFont(name string) int {
	panic(`unimplemented`)
}
func (c *Context) FontBlur() float64 {
	panic(`unimplemented`)
}
func (c *Context) FontFace() string {
	panic(`unimplemented`)
}
func (c *Context) FontFaceID() int {
	panic(`unimplemented`)
}
func (c *Context) FontSize() float64 {
	panic(`unimplemented`)
}
func (c *Context) GlobalAlpha() float64 {
	panic(`unimplemented`)
}
func (c *Context) LineCap() nanovgo.LineCap {
	panic(`unimplemented`)
}
func (c *Context) LineJoin() nanovgo.LineCap {
	panic(`unimplemented`)
}
func (c *Context) MiterLimit() float64 {
	panic(`unimplemented`)
}
func (c *Context) StrokeWidth() float64 {
	panic(`unimplemented`)
}
func (c *Context) TextLetterSpacing() float64 {
	panic(`unimplemented`)
}
func (c *Context) TextLineHeight() float64 {
	panic(`unimplemented`)
}
func (c *Context) TextMetrics() (float64, float64, float64) {
	panic(`unimplemented`)
}
func (c *Context) TextAlign() nanovgo.Align {
	panic(`unimplemented`)
}
func (c *Context) CurrentTransform() nanovgo.TransformMatrix {
	return c.TransformMatrix
}

/* Online text operations */

const maxFontTextures = 4

func (c *Context) TextRune(x, y float64, runes []rune) float64 {
	c.assertFrameStarted()

	p := c
	for ; p.parent != nil; p = p.parent {
	}

	scale := float64(min(c.CurrentTransform().GetAverageScale(), 4)) * c.devicePxRatio // TODO Extract the diagonal from current transform.
	invScale := 1.0 / scale
	if c.Hfont < 0 {
		return 0
	}

	c.Fs.SetSize(float32(c.Fontsiz * scale))
	c.Fs.SetSpacing(float32(c.Lsp * scale))
	c.Fs.SetBlur(float32(c.Fblur * scale))
	c.Fs.SetAlign(fontstashmini.ALIGN_LEFT)
	c.Fs.SetFont(c.Hfont)

	left := len(p.SpriteUnits)
	right := left + max(2, len(runes)) // Not less than two quads.
	p.SpriteUnits = append(p.SpriteUnits, make([]SpriteUnit, right-left)...)

	iter := c.Fs.TextIterForRunes(float32(x*scale), float32(y*scale), runes)
	prevIter := iter

	reallocateImage := false
	i := 0
	for {
		quad, ok := iter.Next()
		if !ok {
			break
		}
		if iter.PrevGlyph == nil || iter.PrevGlyph.Index == -1 {
			reallocateImage = true
		}
		// TODO -1 means 'do kerning'
		if iter.PrevGlyph.Index == -1 && c.Hfont < maxFontTextures-1 {
			iter = prevIter
			quad, _ = iter.Next() // try again
		}
		prevIter = iter
		p.SpriteUnits[left:right][i] = SpriteUnit{
			Hfont: c.Hfont,
			Clip:  geom.Rect(float64(quad.X0), float64(quad.Y0), float64(quad.X1), float64(quad.Y1)),
			Tc:    geom.Rect(float64(quad.S0), float64(quad.T0), float64(quad.S1), float64(quad.T1)),
		}
		i++
	}

	_ = c.add(op.TextRune, RenderOp{
		Left:  left,
		Right: right,
		Runes: runes,
		Args:  [10]float64{invScale, x, y, cond(reallocateImage, 1.0, 0)},
	})

	return float64(iter.X)
}

func (c *Context) TextBounds(x, y float64, runes []rune) (float64, geom.Rectangle) {
	scale := 1.0 // * c.devicePxRatio
	invScale := 1.0 / scale
	if c.Hfont < 0 {
		return 0, geom.Rectangle{}
	}

	c.Fs.SetSize(float32(c.Fontsiz * scale))
	c.Fs.SetSpacing(float32(c.Lsp * scale))
	c.Fs.SetBlur(float32(c.Fblur * scale))
	c.Fs.SetFont(c.Hfont)

	width, bounds, ok := c.Fs.TextBoundsOfRunes(float32(x*scale), float32(y*scale), runes)
	if !ok {
		bounds.Min.Y, bounds.Max.Y = c.Fs.LineBounds(float32(y * scale))
		bounds.Max = bounds.Max.Mul(invScale)
		bounds.Min = bounds.Min.Mul(invScale)
	}
	return float64(width) * invScale, bounds
}

// Renderer is an object that takes a context and interprets every graphics operation
// in it to obtain some graphical output.
// Typically, Renderer is called with the periodicity of the screen refresh.
type Renderer interface {
	Run(c *Context)
}

// Windower is an object that enables Contraption to communicate with the operating system.
// Windower receives input events and manages the window state.
//
// Windower does not handle screen refresh blocking, it must be handled by the Renderer.
type Windower interface {
	SetupInputCallbacks(emit func(ev any, pt geom.Point, t time.Time), u *Events)
	PollEvents(u *Events)
	WaitEvents(u *Events)
	Next(u *Events) (ok bool, w, h int, scale float64)
	Develop(u *Events)
}
