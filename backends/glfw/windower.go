package glfw

import (
	"image"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/neputevshina/contraption"
	"github.com/neputevshina/geom"
)

// type Windower interface {
// 	SetupInputCallbacks(emit func(ev any, pt geom.Point, t time.Time), u *contraption.Events)
// 	PollEvents(u *contraption.Events)
// 	WaitEvents(u *contraption.Events)
// }

type Windower struct {
	*glfw.Window
}

func New(windowRect image.Rectangle) *Windower {
	wer := &Windower{}

	_ = glfw.Init()
	glfw.WindowHint(glfw.Samples, 4)
	wer.Window, _ = glfw.CreateWindow(windowRect.Dx(), windowRect.Dy(), "", nil, nil)
	if windowRect.Min.X != 0 && windowRect.Min.Y != 0 {
		wer.Window.SetPos(windowRect.Min.X, windowRect.Min.Y)
	}
	wer.Window.MakeContextCurrent()
	gl.Init()

	if glfw.ExtensionSupported("GLX_EXT_swap_control_tear") || glfw.ExtensionSupported("WGL_EXT_swap_control_tear") {
		glfw.SwapInterval(-1)
	} else {
		glfw.SwapInterval(1)
	}

	return wer
}

func (wer *Windower) SetupInputCallbacks(emit func(ev any, pt geom.Point, t time.Time), u *contraption.Events) {
	w := wer.Window

	emit2 := func(ev interface{}) {
		// Take the coordinate from previous event.
		// Actually, Hover is the only event type that can update the cursor position.
		n := time.Now()
		emit(ev, u.Trace[0].Pt, n)
	}
	w.SetCursorPosCallback(func(_ *glfw.Window, xpos, ypos float64) {
		emit(contraption.Hover{}, geom.Pt(xpos, ypos), time.Now())
	})
	w.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
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
			emit2(contraption.Click(v))
		case glfw.Release:
			emit2(contraption.Unclick(v))
		case glfw.Repeat:
			emit2(contraption.Click(v))
			emit2(contraption.Unclick(v))
		}
	})
	w.SetScrollCallback(func(_ *glfw.Window, xoff, yoff float64) {
		if yoff != 0 {
			emit2(contraption.Scroll(-yoff))
		}
		if xoff != 0 {
			emit2(contraption.Sweep(xoff))
		}
	})
	w.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			emit2(contraption.Press{Key: key})
		case glfw.Release:
			emit2(contraption.Release{Key: key})
		case glfw.Repeat:
			emit2(contraption.Release{Key: key})
			emit2(contraption.Press{Key: key})
		}
	})
	w.SetCharCallback(func(w *glfw.Window, r rune) {
		tr := &u.Trace[0]
		switch e := tr.E.(type) {
		case contraption.Press:
			e.Rune = r
			tr.E = e
		default:
			// panic(`unreachable`)
		}
	})
	w.SetDropCallback(func(w *glfw.Window, names []string) {
		emit2(contraption.Drop{Paths: names})
	})
}

func (wer *Windower) WaitEvents(_ *contraption.Events) {
	glfw.WaitEvents()
}

func (wer *Windower) PollEvents(_ *contraption.Events) {
	glfw.PollEvents()
}

func (wer *Windower) Develop(_ *contraption.Events) {
	wer.Window.SwapBuffers()
}

func (wer *Windower) Next(_ *contraption.Events) (ok bool, w, h int, scale float64) {
	window := wer.Window
	if window.ShouldClose() {
		return false, 0, 0, 0
	}
	ok = true

	// TODO Decouple backend and put more stuff in that .next().

	w, h = window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(w), int32(h))

	w, h = window.GetSize()

	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)

	sc, _ := window.GetContentScale()
	scale = float64(sc)
	return
}
