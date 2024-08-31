package contraption

import (
	"encoding/gob"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/h2non/filetype"
	"github.com/neputevshina/geom"
)

type EventPoint struct {
	E     interface{}
	Pt    geom.Point
	T     time.Time
	z, zc int
}

func (ev EventPoint) valuestring() string {
	t := strings.Replace(fmt.Sprintf("%T", ev.E), "contraption.", "", -1)
	return fmt.Sprintf("%s(%v)", t, ev.E)
}

func (ev EventPoint) String() string {
	if ev.E == nil {
		return "<empty>"
	}
	return fmt.Sprintf("%v: %s@%v", ev.T, ev.valuestring(), ev.Pt)
}

// Click is a mouse click event, representing a button number.
type Click int

type Hover struct{}

//type Message func() interface{}

// Unclick is a mouse button release event, representing a button number.
type Unclick int

// Scroll is a vertical mouse scroll event.
// Positive values represent scrolling up, negative — scrolling down.
type Scroll int

// Sweep is a horizontal mouse scroll event.
// Negative values represent scrolling left, positive — scrolling right.
type Sweep int

// Press is a keyboard key press event, representing the key code
// with Rune.
//
// Rune part is not matched in regular expressions.
type Press struct {
	Key
	Rune rune
}

func (p Press) key() Key { return p.Key }

// Release is a keyboard key release event, representing the key code with
// optional Rune if releasing emitted a text enter.
//
// Rune part is not matched in regular expressions.
type Release struct {
	Key
	Rune rune
}

func (p Release) key() Key { return p.Key }

// Drag is an event of dragged into application file.
// TODO: change events backend to one that supports Drag events — GLFW does not.
// type Drag struct {
// 	Paths []string
// 	mime  string
// }

// Drag is an event of dragged and dropped into application file.
type Drop struct {
	Paths []string
	mime  string
}

// Every event type above must be registered in Gob for working record-replay functionality.
func init() {
	gob.Register(Click(0))
	gob.Register(Hover{})
	gob.Register(Unclick(0))
	gob.Register(Scroll(0))
	gob.Register(Sweep(0))
	gob.Register(Press{})
	gob.Register(Release{})
	gob.Register(Drop{})
}

// Special values of rune.
//
// If one needs to enter those code points, they either writing
// a terminal emulator and should insert such special characters
// using Ctrl+Shift+<Key> instead, or they are plainly doing
// something wrong.
const (
	RuneDelete    = '\x7f' // ASCII Delete
	RuneBackspace = '\b'   // ASCII Backspace
	RuneLeft      = '\x11' // ASCII Device Control 1 — Left cursor key
	RuneDown      = '\x12' // ASCII Device Control 2 — Down ...
	RuneUp        = '\x13' // ASCII Device Control 3 — Up   ...
	RuneRight     = '\x14' // ASCII Device Control 4 — Right ...
)

type dummyerr struct{}

func (*dummyerr) Error() string {
	return ""
}

const (
	anyShift = 1000000
	anyCtrl  = 1000001
)

// nameevent converts type name and value stirng into event value.
// May be later replaced with reflection and type registry.
func nameevent(typ string, value string) any {
	rechar := regexp.MustCompile("'.'")
	renum := regexp.MustCompile("[-+]?(0|[1-9][0-9]*)")
	intv := 0
	var keyv glfw.Key
	var err error = &dummyerr{}

	if value == "" {
		// Skip this whole else chain.
		err = nil
	} else if typ == `Drag` || typ == `Drop` {
		if !filetype.IsMIMESupported(value) {
			panic(`unsupported MIME type: ` + value)
		}
		err = nil
	} else if rechar.FindString(value) == value {
		// runev = []rune(value)[0]
		err = nil
	} else if renum.FindString(value) == value {
		intv, err = strconv.Atoi(value)
	} else {
		err = nil
		ok := false
		keyv, ok = keynames[value]
		if !ok {
			panic("no such constant “" + value + "”")
		}
	}
	if err != nil {
		panic("check parser for errors: " + err.Error())
	}

	v := any(nil)
	switch typ {
	case "Click":
		v = Click(intv)
	case "Unclick":
		v = Unclick(intv)
	// case "Rune":
	// 	v = Rune(runev)
	case "Hover":
		v = Hover{}
	case "Press":
		v = Press{Key: Key(keyv)}
	case "Release":
		v = Release{Key: Key(keyv)}
	case "Scroll":
		v = Scroll(intv)
	case "Sweep":
		v = Sweep(intv)
	// case "Drag":
	// 	v = Drag{mime: value}
	case "Drop":
		v = Drop{mime: strings.TrimSpace(value)}
	default:
		panic("unknown type: " + typ)
	}
	return v
}

// holdable returns true on every type that is meant to be held until
// it's complement value goes to the bottom of trace
func holdable(v any) bool {
	switch v.(type) {
	// case Drag:
	// 	return true
	case Click:
		return true
	case Press:
		return true
	}
	return false
}

// complement returns for the value of holdable type a value of a type that is complement to it
func complement(v any) any {
	switch v := v.(type) {
	// case Drag:
	// 	return Drop(v)
	// case Drop:
	// 	return Drag(v)
	case Click:
		return Unclick(v)
	case Unclick:
		return Click(v)
	case Press:
		return Release(v)
	case Release:
		return Press(v)
	}
	return nil
}
