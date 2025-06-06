package contraption

import (
	"crypto/sha1"
	"hash"
	"image"
	"reflect"
	"unicode/utf8"

	"github.com/neputevshina/contraption/nanovgo"
	"github.com/neputevshina/contraption/nanovgo/fontstashmini"
	"github.com/neputevshina/geom"
)

type nvguop struct {
	tag   uintptr
	args  [10]float64
	iargs [10]int

	fontsiz float64
	hfont   int
	himage  int

	lsp     float64
	fblur   float64
	strokep nanovgo.Paint
	strokec nanovgo.Color
	fillp   nanovgo.Paint
	fillc   nanovgo.Color
	strokew float64

	linecap nanovgo.LineCap
	nanovgo.Direction
	nanovgo.Winding
	nanovgo.Align
	nanovgo.TransformMatrix
	nanovgo.GlyphPosition

	runes []rune

	left, right int
}

type contextImage struct {
	deleted bool
	h       int

	image.Image
	data []byte
	nanovgo.ImageFlags
	wh image.Point
}

type contextFont struct {
	deleted  bool
	hbackend int

	data     []byte
	freeData byte
}

// Context is a draw list with the interface of Nanovgo.
//
// It is interpreted by the backend to draw on screen or to do anything else.
// TODO All transformations must be resolved by Context, not by the backend.
type Context struct {
	publicContext
	state  uintptr
	hash   [512]byte
	hasher hash.Hash

	parent *Context
}

// Sub returns a persistent context that could be replayed later.
func (c *Context) Sub() *Context {
	// NOTE Subcontexts are not hashed.
	return &Context{
		state:  c.state, // TODO Subcontexts are not required to Begin/EndFrame
		parent: c,
		publicContext: publicContext{
			fs: c.fs,
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
			fs:     fontstashmini.New(512, 512),
			Images: []contextImage{{}},
		},
		hasher: sha1.New(),
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
	nvguop // Current state

	fs *fontstashmini.FontStash

	Log           []nvguop
	Images        []contextImage
	Fonts         []contextFont
	SpriteUnits   []SpriteUnit
	devicePxRatio float64
}

func (c *Context) assertPathStarted() {
	if c.state != functag((*Context).BeginPath) {
		panic(`contraption.Context: vector operation before BeginPath`)
	}
}

func (c *Context) assertFrameStarted() {
	if c.parent != nil {
		return
	}
	if c.state != functag((*Context).BeginFrame) && c.state != 1 {
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
	_ = c.add((*Context).IntersectScissor, nvguop{
		args: [10]float64{x, y, w, h},
	})
}

func functag(f any) uintptr {
	return uintptr(reflect.ValueOf(f).UnsafePointer())
}

func (c *Context) add(f any, op nvguop) uintptr {
	op.tag = functag(f)
	c.Log = append(c.Log, op)
	// TODO FNV-1a adds >2% of cpu usage, without any benefits
	// Also previous testing shows that there are several hundreds of megabytes per second
	// passing through the hash when actively moving a mouse.
	p := c
	for ; p.parent != nil; p = p.parent {
	}
	ap := func(bs []byte) {
		p.hasher.Write(bs)
	}
	if op.runes != nil {
		for _, r := range op.runes {
			bb := [8]byte{}
			c.hasher.Write(bb[:utf8.EncodeRune(bb[:], r)])
		}
	}
	ap(asbs(op.tag))
	ap(asbs(op.args))
	ap(asbs(op.iargs))
	ap(asbs(op.fillc))
	ap(asbs(op.fillp))
	ap(asbs(op.strokec))
	ap(asbs(op.strokep))
	ap(asbs(op.fontsiz))
	ap(asbs(op.hfont))
	ap(asbs(op.himage))
	return op.tag
}

type deferOp int

const (
	opBeginFrame deferOp = iota + 1
	opEndFrame
	opCancelFrame
	opCreateImageRGBA
	opCreateImageFromGoImage
	opUpdateImage
	opDeleteImage
	opCreateFontFromMemory
	opCircle
	opRect
	opEllipse
	opRoundedRect
	opBeginPath
	opClosePath
	opFill
	opStroke
	opAlso
	opArc
	opArcTo
	opBezierTo
	opLineTo
	opMoveTo
	opQuadTo
	opPathWinding
	opReset
	opResetScissor
	opResetTransform
	opRestore
	opSave
	opRotate
	opScale
	opScissor
	opSkewX
	opSkewY
	opSetTransform
	opSetTransformByValue
	opTranslate
	opSetFillColor
	opSetFillPaint
	opSetFontBlur
	opSetFontFace
	opSetFontFaceID
	opSetFontSize
	opSetGlobalAlpha
	opSetLineCap
	opSetLineJoin
	opSetMiterLimit
	opSetStrokeColor
	opSetStrokePaint
	opSetStrokeWidth
	opSetTextAlign
	opSetTextLetterSpacing
	opSetTextLineHeight
	opTextRune
	opReplay
)

/* Frame and context state */

func (c *Context) BeginFrame(windowWidth, windowHeight int, devicePixelRatio float64) {
	if c.state != 0 {
		panic(`contraption.Context: BeginFrame can only be called at the start of a frame`)
	}
	c.SpriteUnits = c.SpriteUnits[:0]

	st := c.add((*Context).BeginFrame, nvguop{
		iargs: [10]int{windowWidth, windowHeight},
		args:  [10]float64{devicePixelRatio},
	})
	c.devicePxRatio = devicePixelRatio
	c.state = st
}
func (c *Context) EndFrame() (oldhash [512]byte) {
	_ = c.add((*Context).EndFrame, nvguop{})
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
	c.Images = append(c.Images, contextImage{
		Image:      nil,
		data:       data,
		ImageFlags: imageFlags,
		wh:         image.Pt(w, h),
	})
	_ = c.add((*Context).CreateImageRGBA, nvguop{
		himage: len(c.Images),
	})
	return len(c.Images)
}
func (c *Context) CreateImageFromGoImage(imageFlag nanovgo.ImageFlags, img image.Image) int {
	if c.parent != nil {
		panic(`contraption.Context: can't create image from subcontext`)
	}
	c.Images = append(c.Images, contextImage{
		Image:      img,
		ImageFlags: imageFlag,
	})
	_ = c.add((*Context).CreateImageFromGoImage, nvguop{
		himage: len(c.Images) - 1,
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
	if c.Images[img].deleted {
		panic(`contraption.Context: double image delete`)
	}
	c.Images[img] = contextImage{deleted: true}
}
func (c *Context) ImageSize(img int) (int, int, error) {
	panic(`unimplemented`)
}

/* Fonts */

func (c *Context) CreateFontFromMemory(name string, data []byte, freeData uint8) int {
	if c.parent != nil {
		panic(`contraption.Context: can't create font from subcontext`)
	}
	c.Fonts = append(c.Fonts, contextFont{
		data:     data,
		freeData: freeData,
	})
	c.fs.AddFontFromMemory(name, data, freeData)
	return len(c.Fonts) - 1
}

/* Shapes */

func (c *Context) Circle(cx, cy, r float64) {
	c.assertPathStarted()
	_ = c.add((*Context).Circle, nvguop{
		args: [10]float64{cx, cy, r},
	})
}
func (c *Context) Rect(x, y, w, h float64) {
	c.assertPathStarted()
	_ = c.add((*Context).Rect, nvguop{
		args: [10]float64{x, y, w, h},
	})
}
func (c *Context) Ellipse(cx, cy, rx, ry float64) {
	c.assertPathStarted()
	_ = c.add((*Context).Ellipse, nvguop{
		args: [10]float64{cx, cy, rx, ry},
	})
}
func (c *Context) RoundedRect(x, y, w, h, r float64) {
	c.assertPathStarted()
	_ = c.add((*Context).RoundedRect, nvguop{
		args: [10]float64{x, y, w, h, r},
	})
}

/* Paths */

func (c *Context) BeginPath() {
	c.assertFrameStarted()
	// if c.state != functag((*Context).ClosePath) {
	// 	panic(`contraption.Context: can't BeginPath before ClosePath`)
	// }
	st := c.add((*Context).BeginPath, nvguop{})
	c.state = st
}
func (c *Context) ClosePath() {
	c.assertPathStarted()
	_ = c.add((*Context).ClosePath, nvguop{})
	c.state = functag((*Context).BeginFrame)
}
func (c *Context) Fill() {
	if c.state == 1 {
		panic(`contraption.Context: another Fill can be called only after Also`)
	}
	_ = c.add((*Context).Fill, nvguop{})
	c.state = 1
}
func (c *Context) Stroke() {
	if c.state == 1 {
		panic(`contraption.Context: another Stroke can be called only after Also`)
	}
	_ = c.add((*Context).Stroke, nvguop{})
	c.state = 1
}
func (c *Context) Also() {
	if c.state != 1 {
		panic(`contraption.Context: can't use Also before Stroke or Fill`)
	}
	_ = c.add((*Context).Also, nvguop{})
	c.state = 2
}
func (c *Context) Arc(cx, cy, r, a0, a1 float64, dir nanovgo.Direction) {
	c.assertPathStarted()
	_ = c.add((*Context).Arc, nvguop{
		args:      [10]float64{cx, cy, r, a0, a1},
		Direction: dir,
	})
}
func (c *Context) ArcTo(x1, y1, x2, y2, radius float64) {
	c.assertPathStarted()
	_ = c.add((*Context).ArcTo, nvguop{
		args: [10]float64{x1, y1, x2, y2, radius},
	})
}
func (c *Context) BezierTo(c1x, c1y, c2x, c2y, x, y float64) {
	c.assertPathStarted()
	_ = c.add((*Context).BezierTo, nvguop{
		args: [10]float64{c1x, c1y, c2x, c2y, x, y},
	})
}
func (c *Context) LineTo(x, y float64) {
	c.assertPathStarted()
	_ = c.add((*Context).LineTo, nvguop{
		args: [10]float64{x, y},
	})
}
func (c *Context) MoveTo(x, y float64) {
	c.assertPathStarted()
	_ = c.add((*Context).MoveTo, nvguop{
		args: [10]float64{x, y},
	})
}
func (c *Context) QuadTo(cx, cy, x, y float64) {
	c.assertPathStarted()
	_ = c.add((*Context).QuadTo, nvguop{
		args: [10]float64{cx, cy, x, y},
	})
}
func (c *Context) PathWinding(winding nanovgo.Winding) {
	c.assertPathStarted()
	_ = c.add((*Context).PathWinding, nvguop{
		Winding: winding,
	})
}

/* State management */

func (c *Context) Reset() {
	_ = c.add((*Context).Reset, nvguop{})
}
func (c *Context) ResetScissor() {
	_ = c.add((*Context).ResetScissor, nvguop{})
}
func (c *Context) ResetTransform() {
	_ = c.add((*Context).ResetTransform, nvguop{})
}
func (c *Context) Restore() {
	_ = c.add((*Context).Restore, nvguop{})
}
func (c *Context) Save() {
	_ = c.add((*Context).Save, nvguop{})
}

// TODO All transformations must be resolved by Context, not by the backend.
/* Transformation mutators */

func (c *Context) Rotate(angle float64) {
	_ = c.add((*Context).Rotate, nvguop{})
}
func (c *Context) Scale(x, y float64) {
	panic(`unimplemented`)
}
func (c *Context) Scissor(x, y, w, h float64) {
	c.assertFrameStarted()
	_ = c.add((*Context).Scissor, nvguop{
		args: [10]float64{x, y, w, h},
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
	_ = c.add((*Context).SetTransform, nvguop{
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
	c.fillc = color
	_ = c.add((*Context).SetFillColor, nvguop{
		fillc: color,
	})
}
func (c *Context) SetFillPaint(paint nanovgo.Paint) {
	c.fillp = paint
	_ = c.add((*Context).SetFillPaint, nvguop{
		fillp: paint,
	})
}
func (c *Context) SetFontBlur(blur float64) {
	panic(`unimplemented`)
}
func (c *Context) SetFontFace(font string) {
	panic(`unimplemented`)
}
func (c *Context) SetFontFaceID(font int) {
	c.hfont = font
	_ = c.add((*Context).SetFontFaceID, nvguop{
		hfont: font,
	})
}
func (c *Context) SetFontSize(size float64) {
	// TODO Convert from cap to em in other backends
	c.fontsiz = size
	_ = c.add((*Context).SetFontSize, nvguop{
		fontsiz: size,
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
	c.strokec = color
	_ = c.add((*Context).SetStrokeColor, nvguop{
		strokec: color,
	})
}
func (c *Context) SetStrokePaint(paint nanovgo.Paint) {
	c.strokep = paint
	_ = c.add((*Context).SetStrokePaint, nvguop{
		strokep: paint,
	})
}
func (c *Context) SetStrokeWidth(width float64) {
	c.strokew = width
	_ = c.add((*Context).SetStrokeWidth, nvguop{
		strokew: width,
	})
}
func (c *Context) SetTextAlign(align nanovgo.Align) {
	c.Align = align
	_ = c.add((*Context).SetTextAlign, nvguop{
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
	if c.hfont < 0 {
		return 0
	}

	c.fs.SetSize(float32(c.fontsiz * scale))
	c.fs.SetSpacing(float32(c.lsp * scale))
	c.fs.SetBlur(float32(c.fblur * scale))
	c.fs.SetAlign(fontstashmini.ALIGN_LEFT)
	c.fs.SetFont(c.hfont)

	left := len(p.SpriteUnits)
	right := left + max(2, len(runes)) // Not less than two quads.
	p.SpriteUnits = append(p.SpriteUnits, make([]SpriteUnit, right-left)...)

	iter := c.fs.TextIterForRunes(float32(x*scale), float32(y*scale), runes)
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
		if iter.PrevGlyph.Index == -1 && c.hfont < maxFontTextures-1 {
			iter = prevIter
			quad, _ = iter.Next() // try again
		}
		prevIter = iter
		p.SpriteUnits[left:right][i] = SpriteUnit{
			Hfont: c.hfont,
			Clip:  geom.Rect(float64(quad.X0), float64(quad.Y0), float64(quad.X1), float64(quad.Y1)),
			Tc:    geom.Rect(float64(quad.S0), float64(quad.T0), float64(quad.S1), float64(quad.T1)),
		}
		i++
	}

	_ = c.add((*Context).TextRune, nvguop{
		left:  left,
		right: right,
		runes: runes,
		args:  [10]float64{invScale, x, y, cond(reallocateImage, 1.0, 0)},
	})

	return float64(iter.X)
}

func (c *Context) TextBounds(x, y float64, runes []rune) (float64, geom.Rectangle) {
	scale := 1.0 // * c.devicePxRatio
	invScale := 1.0 / scale
	if c.hfont < 0 {
		return 0, geom.Rectangle{}
	}

	c.fs.SetSize(float32(c.fontsiz * scale))
	c.fs.SetSpacing(float32(c.lsp * scale))
	c.fs.SetBlur(float32(c.fblur * scale))
	c.fs.SetFont(c.hfont)

	width, bounds, ok := c.fs.TextBoundsOfRunes(float32(x*scale), float32(y*scale), runes)
	if !ok {
		bounds.Min.Y, bounds.Max.Y = c.fs.LineBounds(float32(y * scale))
		bounds.Max = bounds.Max.Mul(invScale)
		bounds.Min = bounds.Min.Mul(invScale)
	}
	return float64(width) * invScale, bounds
}
