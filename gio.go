package contraption

import (
	"reflect"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/h2non/filetype"
)

// func concretenew(config Config, wo *World) {
// 	var err error

// 	_ = glfw.Init()
// 	glfw.WindowHint(glfw.Samples, 4)
// 	wo.Window.window, _ = glfw.CreateWindow(config.WindowRect.Dx(), config.WindowRect.Dy(), "", nil, nil)
// 	if config.WindowRect.Min.X != 0 && config.WindowRect.Min.Y != 0 {
// 		wo.Window.SetPos(config.WindowRect.Min.X, config.WindowRect.Min.Y)
// 	}
// 	wo.Window.MakeContextCurrent()
// 	gl.Init()

// 	// FIXME

// 	if glfw.ExtensionSupported("GLX_EXT_swap_control_tear") || glfw.ExtensionSupported("WGL_EXT_swap_control_tear") {
// 		println("tear control is supported")
// 		glfw.SwapInterval(-1)
// 	} else {
// 		glfw.SwapInterval(1)
// 	}

// 	wo.cctx, err = nanovgo.NewContext(0)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func concretewait(_ *Events) {
// 	glfw.WaitEvents()
// }

// func concretepoll(_ *Events) {
// 	glfw.PollEvents()
// }

// func (wo *World) windowDevelop() {
// 	if wo.BeforeVgo != nil {
// 		wo.BeforeVgo()
// 	}
// 	wo.BeforeVgo = nil
// 	_ = wo.Vgo.EndFrame()
// 	if wo.Events.tempcur == 0 {
// 		// Retain if was not changed
// 		wo.rer.Run(wo.Vgo)
// 		wo.wer.Develop(wo.Events)
// 	}
// }

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
	// No Super/Win keys because they must be reserved for user's desktop environment.

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

	// TODO Adapt parser to support these.
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
			return pk.key() == keynames[`LShift`] ||
				pk.key() == keynames[`RShift`]
		case anyCtrl:
			return pk.key() == keynames[`LCtrl`] ||
				pk.key() == keynames[`RCtrl`]
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
