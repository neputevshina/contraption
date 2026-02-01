package contraption

import (
	"image"
	"reflect"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/h2non/filetype"
	"github.com/neputevshina/contraption/nanovgo"
	"github.com/neputevshina/contraption/op"
	"github.com/neputevshina/geom"
)

func concretenew(config Config, wo *World) {
	var err error

	_ = glfw.Init()
	glfw.WindowHint(glfw.Samples, 4)
	wo.Window.window, _ = glfw.CreateWindow(config.WindowRect.Dx(), config.WindowRect.Dy(), "", nil, nil)
	if config.WindowRect.Min.X != 0 && config.WindowRect.Min.Y != 0 {
		wo.Window.SetPos(config.WindowRect.Min.X, config.WindowRect.Min.Y)
	}
	wo.Window.MakeContextCurrent()
	gl.Init()

	// FIXME

	if glfw.ExtensionSupported("GLX_EXT_swap_control_tear") || glfw.ExtensionSupported("WGL_EXT_swap_control_tear") {
		println("tear control is supported")
		glfw.SwapInterval(-1)
	} else {
		glfw.SwapInterval(1)
	}

	wo.cctx, err = nanovgo.NewContext(0)
	if err != nil {
		panic(err)
	}
}

func concretewait(_ *Events) {
	glfw.WaitEvents()
}

func concretepoll(_ *Events) {
	glfw.PollEvents()
}

func (wi *Window) rect() (r image.Rectangle) {
	r.Min = image.Pt(wi.window.GetPos())
	r.Max = r.Min.Add(image.Pt(wi.window.GetSize()))
	return
}

func (wo *World) windowDevelop() {
	if wo.BeforeVgo != nil {
		wo.BeforeVgo()
	}
	wo.BeforeVgo = nil
	_ = wo.Vgo.EndFrame()
	if wo.Events.tempcur == 0 {

		// Retain if was not changed
		runcontext(wo.cctx, wo.Vgo)
		wo.Window.SwapBuffers()

	}
	wo.Vgo.Log = wo.Vgo.Log[:0]
	// wo.Vgo.EndFrame()
}

type window = *glfw.Window

type Key = glfw.Key

func keyname(name string) (Key, bool) {
	k, ok := keynames[name]
	return Key(k), ok
}

// keynames maps key names to keycodes.
var keynames = map[string]glfw.Key{
	"LCtrl":  glfw.KeyLeftControl,
	"LShift": glfw.KeyLeftShift,
	"LAlt":   glfw.KeyLeftAlt,
	"RCtrl":  glfw.KeyRightControl,
	"RShift": glfw.KeyRightShift,
	"RAlt":   glfw.KeyRightAlt,
	"Ctrl":   anyCtrl,
	"Shift":  anyShift,
	// No Super/Win keys because they must be reserved for user's desktop.

	"Escape":      glfw.KeyEscape,
	"PrintScreen": glfw.KeyPrintScreen,
	"ScrollLock":  glfw.KeyScrollLock,
	"NumLock":     glfw.KeyNumLock,
	"CapsLock":    glfw.KeyCapsLock,
	"Pause":       glfw.KeyPause,
	"Insert":      glfw.KeyInsert,
	"Delete":      glfw.KeyDelete,
	"Home":        glfw.KeyHome,
	"End":         glfw.KeyEnd,
	"PageUp":      glfw.KeyPageUp,
	"PageDown":    glfw.KeyPageUp,
	"Backspace":   glfw.KeyBackspace,
	"Return":      glfw.KeyEnter, "Enter": glfw.KeyEnter,
	"Tab":  glfw.KeyTab,
	"Menu": glfw.KeyMenu, "Context": glfw.KeyMenu,

	"Q": glfw.KeyQ, "W": glfw.KeyW, "E": glfw.KeyE, "R": glfw.KeyR, "T": glfw.KeyT, "Y": glfw.KeyY, "U": glfw.KeyU, "I": glfw.KeyI, "O": glfw.KeyO, "P": glfw.KeyP,
	"A": glfw.KeyA, "S": glfw.KeyS, "D": glfw.KeyD, "F": glfw.KeyF, "G": glfw.KeyG, "H": glfw.KeyH, "J": glfw.KeyJ, "K": glfw.KeyK, "L": glfw.KeyL,
	"Z": glfw.KeyZ, "X": glfw.KeyX, "C": glfw.KeyC, "V": glfw.KeyV, "B": glfw.KeyB, "N": glfw.KeyN, "M": glfw.KeyM,

	"0": glfw.Key0, "1": glfw.Key1, "2": glfw.Key2, "3": glfw.Key3, "4": glfw.Key4,
	"5": glfw.Key5, "6": glfw.Key6, "7": glfw.Key7, "8": glfw.Key8, "9": glfw.Key9,

	// TODO Adapt parser to support this.
	// Hint: Token rule in regexp.peg
	`Comma`: glfw.KeyComma, `,`: glfw.KeyComma,
	`Period`: glfw.KeyPeriod, `.`: glfw.KeyPeriod,
	`Slash`: glfw.KeySlash, `/`: glfw.KeySlash,
	`Backslash`: glfw.KeyBackslash, `\`: glfw.KeyBackslash,
	`Semicolon`: glfw.KeySemicolon, `;`: glfw.KeySemicolon,
	`Apostrophe`: glfw.KeyApostrophe, `'`: glfw.KeyApostrophe,
	`Grave`: glfw.KeyGraveAccent, `Backtick`: glfw.KeyGraveAccent, "`": glfw.KeyGraveAccent,
	"-": glfw.KeyMinus,
	"=": glfw.KeyEqual,

	"←": glfw.KeyLeft, "Left": glfw.KeyLeft,
	"→": glfw.KeyRight, "Right": glfw.KeyRight,
	"↑": glfw.KeyUp, "Up": glfw.KeyUp,
	"↓": glfw.KeyDown, "Down": glfw.KeyDown,

	"F1": glfw.KeyF1, "F2": glfw.KeyF2, "F3": glfw.KeyF3, "F4": glfw.KeyF4,
	"F5": glfw.KeyF5, "F6": glfw.KeyF6, "F7": glfw.KeyF7, "F8": glfw.KeyF8,
	"F9": glfw.KeyF9, "F10": glfw.KeyF10, "F11": glfw.KeyF11, "F12": glfw.KeyF12,
}

func setupcallbacks(u *Events, window any) {
	w := window.(*glfw.Window)

	emit2 := func(ev interface{}) {
		// Take the coordinate from previous event.
		// Actually, Hover is the only event type that can update the cursor position.
		n := time.Now()
		u.emit(ev, u.Trace[0].Pt, n)
		// println(n, u.Now)
	}
	w.SetCursorPosCallback(func(_ *glfw.Window, xpos, ypos float64) {
		u.emit(Hover{}, geom.Pt(xpos, ypos), time.Now())
	})
	w.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
		// println(`CLICK`, time.Now())
		v := int(button)
		switch button {
		case glfw.MouseButtonLeft:
			v = 1
		case glfw.MouseButtonMiddle:
			v = 2
		case glfw.MouseButtonRight:
			v = 3
		}
		switch action {
		case glfw.Press:
			emit2(Click(v))
		case glfw.Release:
			emit2(Unclick(v))
		case glfw.Repeat:
			emit2(Click(v))
			emit2(Unclick(v))
		}
	})
	w.SetScrollCallback(func(_ *glfw.Window, xoff, yoff float64) {
		if yoff != 0 {
			emit2(Scroll(-yoff))
		}
		if xoff != 0 {
			emit2(Sweep(xoff))
		}
	})
	w.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			emit2(Press{Key: key})
		case glfw.Release:
			emit2(Release{Key: key})
		case glfw.Repeat:
			emit2(Release{Key: key})
			emit2(Press{Key: key})
		}
	})
	w.SetCharCallback(func(w *glfw.Window, r rune) {
		tr := &u.Trace[0]
		switch e := tr.E.(type) {
		case Press:
			e.Rune = r
			tr.E = e
		default:
			// panic(`unreachable`)
		}
	})
	w.SetDropCallback(func(w *glfw.Window, names []string) {
		emit2(Drop{Paths: names})
	})
}

func requals(p EventPoint, inst *rinst) bool {
	sametype := reflect.TypeOf(p.E) == reflect.TypeOf(inst.e)
	if inst.typeonly {
		return sametype
	}
	pk, ok1 := p.E.(keyer)    // Key in event
	ik, ok2 := inst.e.(keyer) // Key in regexp rule
	// If p.E or inst.e is not keyer, they both will fail on default condition.
	if ok1 && ok2 && sametype {
		switch ik.key() {
		case anyShift:
			return pk.key() == glfw.KeyLeftShift ||
				pk.key() == glfw.KeyRightShift
		case anyCtrl:
			return pk.key() == glfw.KeyLeftControl ||
				pk.key() == glfw.KeyRightControl
		}
	}
	pd, ok1 := p.E.(Drop)
	id, ok2 := inst.e.(Drop)
	if ok1 && ok2 && sametype {
		if pd.mime == `` {
			return true
		}
		for _, f := range pd.Paths {
			t, err := filetype.MatchFile(f)
			if err != nil {
				panic(err)
			}
			if id.mime != "" {
				return t.MIME.Value == id.mime
			}
		}
	}
	return p.E == inst.e
}

func runcontext(concrete any, c *Context) {
	vgo := concrete.(*nanovgo.Context)
	vtxb := []nanovgo.Vertex{}
	for i := range c.Log {
		l := &c.Log[i]
		switch l.tag {
		case op.Also:
			vgo.Also()
		case op.Arc:
			vgo.Arc(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]), float32(l.args[4]), l.Direction)
		case op.ArcTo:
			vgo.ArcTo(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]), float32(l.args[4]))
		case op.BeginFrame:
			vgo.BeginFrame(l.iargs[0], l.iargs[1], float32(l.args[0]))
		case op.BeginPath:
			vgo.BeginPath()
		case op.BezierTo:
			vgo.BezierTo(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]), float32(l.args[4]), float32(l.args[5]))
		case op.Block:
			panic(`unimplemented`)
		case op.CancelFrame:
			vgo.CancelFrame()
		case op.Circle:
			vgo.Circle(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]))
		case op.ClosePath:
			vgo.ClosePath()
		case op.CreateFontFromMemory:
			f := &c.Fonts[l.hfont]
			f.hbackend = vgo.CreateFontFromMemory("", f.data, f.freeData)
		case op.CreateImageFromGoImage:
			m := &c.Images[l.himage]
			m.h = vgo.CreateImageFromGoImage(m.ImageFlags, m.Image)
		case op.CreateImageRGBA:
			m := &c.Images[l.himage]
			m.h = vgo.CreateImageRGBA(m.wh.X, m.wh.Y, m.ImageFlags, m.data)
		case op.CurrentTransform:
			panic(`getter, unreachable`)
		case op.DebugDumpPathCache:
			vgo.DebugDumpPathCache()
		case op.Delete:
			vgo.Delete()
		case op.DeleteImage:
			vgo.DeleteImage(c.Images[l.himage].h)
		case op.Ellipse:
			vgo.Ellipse(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]))
		case op.EndFrame:
			vgo.EndFrame()
		case op.Fill:
			vgo.Fill()
		case op.FindFont:
			panic(`unimplemented`)
		case op.FontBlur:
			panic(`getter, unreachable`)
		case op.FontFace:
			panic(`getter, unreachable`)
		case op.FontFaceID:
			panic(`getter, unreachable`)
		case op.FontSize:
			panic(`getter, unreachable`)
		case op.GlobalAlpha:
			panic(`getter, unreachable`)
		case op.ImageSize:
			panic(`getter, unreachable`)
		case op.IntersectScissor:
			vgo.IntersectScissor(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]))
		case op.LineCap:
			panic(`getter, unreachable`)
		case op.LineJoin:
			panic(`getter, unreachable`)
		case op.LineTo:
			vgo.LineTo(float32(l.args[0]), float32(l.args[1]))
		case op.MiterLimit:
			panic(`getter, unreachable`)
		case op.MoveTo:
			vgo.MoveTo(float32(l.args[0]), float32(l.args[1]))
		case op.PathWinding:
			vgo.PathWinding(l.Winding)
		case op.QuadTo:
			vgo.QuadTo(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]))
		case op.Rect:
			vgo.Rect(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]))
		case op.Reset:
			vgo.Reset()
		case op.ResetScissor:
			vgo.ResetScissor()
		case op.ResetTransform:
			vgo.ResetTransform()
		case op.Restore:
			vgo.Restore()
		case op.Rotate:
			vgo.Rotate(float32(l.args[0]))
		case op.RoundedRect:
			vgo.RoundedRect(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]), float32(l.args[4]))
		case op.Save:
			vgo.Save()
		case op.Scale:
			vgo.Scale(float32(l.args[0]), float32(l.args[1]))
		case op.Scissor:
			vgo.Scissor(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]))
		case op.SetFillColor:
			vgo.SetFillColor(l.fillc)
		case op.SetFillPaint:
			p := l.fillp
			p.Image = c.Images[p.Image].h
			vgo.SetFillPaint(p)
		case op.SetFontBlur:
			vgo.SetFontBlur(float32(l.args[0]))
		case op.SetFontFace:
			panic(`unimplemented`)
		case op.SetFontFaceID:
			vgo.SetFontFaceID(l.hfont)
		case op.SetFontSize:
			vgo.SetFontSize(float32(l.args[0]))
		case op.SetGlobalAlpha:
			vgo.SetGlobalAlpha(float32(l.args[0]))
		case op.SetLineCap:
			vgo.SetLineCap(l.linecap)
		case op.SetLineJoin:
			vgo.SetLineJoin(l.linecap)
		case op.SetMiterLimit:
			vgo.SetMiterLimit(float32(l.args[0]))
		case op.SetStrokeColor:
			vgo.SetStrokeColor(l.strokec)
		case op.SetStrokePaint:
			p := l.strokep
			p.Image = c.Images[p.Image].h
			vgo.SetStrokePaint(p)
		case op.SetStrokeWidth:
			vgo.SetStrokeWidth(float32(l.strokew))
		case op.SetTextAlign:
			vgo.SetTextAlign(l.Align)
		case op.SetTextLetterSpacing:
			vgo.SetTextLetterSpacing(float32(l.args[0]))
		case op.SetTextLineHeight:
			vgo.SetTextLineHeight(float32(l.args[0]))
		case op.SetTransform:
			vgo.SetTransform(l.TransformMatrix)
		case op.SetTransformByValue:
			vgo.SetTransformByValue(float32(l.args[0]), float32(l.args[1]), float32(l.args[2]), float32(l.args[3]), float32(l.args[4]), float32(l.args[5]))
		case op.SkewX:
			vgo.SkewX(float32(l.args[0]))
		case op.SkewY:
			vgo.SkewY(float32(l.args[0]))
		case op.Stroke:
			vgo.Stroke()
		case op.StrokeWidth:
			vgo.StrokeWidth()
		case op.TextAlign:
			panic(`getter, unreachable`)
		case op.TextBounds:
			panic(`getter, unreachable`)
		case op.TextLetterSpacing:
			panic(`getter, unreachable`)
		case op.TextLineHeight:
			panic(`getter, unreachable`)
		case op.TextMetrics:
			panic(`getter, unreachable`)
		case op.TextRune:
			// vgo.TextRune(c.fs, float32(l.args[1]), float32(l.args[2]), l.runes)

			sus := c.SpriteUnits[l.left:l.right]
			bf := c.Fonts[l.hfont]

			vtxb = vtxb[:0]
			vtxb = append(vtxb, make([]nanovgo.Vertex, 4*len(sus))...)
			vidx := 0
			invScale := l.args[0]

			if l.args[3] == 1 {
				// Reallocate atlas since we have new glyphs.
				_, ok := vgo.AllocTextAtlas(c.fs, bf.hbackend)
				if !ok {
					panic(``)
				}
			}
			for i := range sus {
				quad := sus[i]
				// Transform corners.
				t := vgo.CurrentTransform()
				c0, c1 := t.TransformPoint(float32(quad.Clip.Min.X*invScale), float32(quad.Clip.Min.Y*invScale))
				c2, c3 := t.TransformPoint(float32(quad.Clip.Max.X*invScale), float32(quad.Clip.Min.Y*invScale))
				c4, c5 := t.TransformPoint(float32(quad.Clip.Max.X*invScale), float32(quad.Clip.Max.Y*invScale))
				c6, c7 := t.TransformPoint(float32(quad.Clip.Min.X*invScale), float32(quad.Clip.Max.Y*invScale))
				//log.Printf("quad(%c) x0=%d, x1=%d, y0=%d, y1=%d, s0=%d, s1=%d, t0=%d, t1=%d\n", iter.CodePoint, int(quad.Clip.Min.X), int(quad.Clip.Max.X), int(quad.Clip.Min.Y), int(quad.Clip.Max.Y), int(1024*quad.Tc.Min.X), int(quad.Tc.Max.X*1024), int(quad.Tc.Min.Y*1024), int(quad.Tc.Max.Y*1024))
				// Create triangles
				vtx := func(x, y float32, u, v float64) nanovgo.Vertex {
					return nanovgo.Vertex{X: float32(x), Y: float32(y), U: float32(u), V: float32(v)}
				}
				vtxb[vidx] = vtx(c2, c3, quad.Tc.Max.X, quad.Tc.Min.Y)
				vtxb[vidx+1] = vtx(c0, c1, quad.Tc.Min.X, quad.Tc.Min.Y)
				vtxb[vidx+2] = vtx(c4, c5, quad.Tc.Max.X, quad.Tc.Max.Y)
				vtxb[vidx+3] = vtx(c6, c7, quad.Tc.Min.X, quad.Tc.Max.Y)
				vidx += 4
			}
			vgo.FlushTextTexture(c.fs, bf.hbackend)
			vgo.RenderText(vtxb)

		case op.Translate:
			vgo.Translate(float32(l.args[0]), float32(l.args[1]))
		case op.UpdateImage:
			f := c.Images[l.himage]
			vgo.UpdateImage(l.himage, f.data)
		default:
			panic(`unreachable`)
		}
	}
}
