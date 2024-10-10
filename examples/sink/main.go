package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/neputevshina/contraption"
	"github.com/neputevshina/geom"
	"golang.org/x/image/font/gofont/goregular"
)

type Sorm = contraption.Sorm
type World struct {
	contraption.World
	Text func(size float64, str []rune) Sorm
}

var (
	yellow = hexpaint(`#ffdb2ca0`)
	dark   = hexpaint(`#f66b0080`)
)

var (
	scale = 1.0

	slider0 = 0.0
	slider1 = 1.0

	numbox = 0.3

	filename0 = ""
	filename1 = ""
)

func main() {
	wo := World{
		World: contraption.New(contraption.Config{}),
	}
	wo.Text = wo.NewVectorText(goregular.TTF)
	f, _ := contraption.NewFont(nil, goregular.TTF, "")
	println(f.Captoem(53))

	for wo.Next() {
		if wo.Match(`!Release(Ctrl)* Press(Ctrl)`) {
			if wo.Match(`Scroll(-1)`) {
				scale += 0.1
			}
			if wo.Match(`Scroll(+1)`) {
				scale -= 0.1
			}
			if wo.Match(`Press(E)`) {
				scale = 1
			}
		}

		wo.Root(
			wo.Pretransform(geom.Scale2d(scale, scale)),
			wo.Compound(
				wo.Sequence2(contraption.SliceSequence2(wo.Trace, func(i int) contraption.Sorm {
					p := `#00000000`
					h := complex128(1.0)
					switch wo.Trace[i].E.(type) {
					case contraption.Click:
						p = `#00ff00`
					case contraption.Unclick:
						p = `#00ff00`
						h = 0.5
					case contraption.Press:
						p = `#ff0000`
					case contraption.Release:
						p = `#ff0000`
						h = 0.5
					case contraption.Hover:
						p = `#00000000`
					case contraption.Drop:
						p = `#ff00ff`
					}
					return wo.Rectangle(-1, -h).Fill(hexpaint(p))
				})),
				wo.Limit(0, 10),
				wo.Hfollow(),
				wo.Valign(0)),
			wo.Compound(
				wo.Halign(0.5),
				wo.Valign(0.5),
				wo.Void(-1, -1),
				wo.Examples(),
			),
			wo.Compound(
				wo.Halign(0),
				wo.Valign(1),
				wo.Void(-1, -1),
				wo.Compound(
					wo.Vfollow(),
					wo.Compound(
						wo.Hfollow(),
						wo.Void(5, 0),
						wo.Label(`Ctrl+Scroll — scale interface, Ctrl+E — reset scale`),
						wo.Void(-1, 0),
						wo.Label(strconv.FormatFloat(scale*100, 'f', 0, 64), "%").
							Cond(func(m contraption.Matcher) {
								if m.Match(`Click:in`) {
									scale = 1.0
								}
							}),
						wo.Void(5, 0),
					),
					wo.Void(0, 5)),
			),
		)

		wo.Develop()
	}
}

func (wo *World) Examples() Sorm {
	// TODO Pluggable DragEffects
	wo.DragEffect = func(interval [2]geom.Point, drag any) Sorm {
		if f, ok := drag.(*float64); ok {
			r := wo.Prevkey(Tag(f, 1)).Rectangle()
			*f = (interval[1].X - r.Min.X) / r.Dx()
			*f = max(0, min(*f, 1))
		}
		return Sorm{}
	}

	return wo.Compound(
		wo.Vfollow(),
		wo.BetweenVoid(0, 64),
		wo.Compound(
			wo.Hfollow(),
			wo.BetweenVoid(32, 0),
			wo.Example(`Align child's center to container's top`, func() Sorm {
				return wo.Compound(
					wo.Limit(100, 100),
					wo.Compound(
						wo.Rectangle(-1, -1).Stroke(dark),
						wo.Compound(
							wo.Halign(0.5),
							wo.Valign(0.5),
							wo.Vshrink(),
							wo.Void(-1, 0).Voverride(),
							wo.Rectangle(-1, -1).Fill(yellow),
							wo.Label(`abc`))))
			}),
			wo.Example(`Align child's center to container's top`, func() Sorm {
				return wo.Compound(
					wo.Limit(100, 100),
					wo.Compound(
						wo.Valign(0.5),
						wo.Rectangle(-1, -1).Stroke(dark),
						wo.Compound(
							wo.Halign(0.5),
							wo.Valign(0),
							wo.Vshrink(),
							wo.Void(-1, 0).Voverride(),
							wo.Rectangle(-1, -1).Fill(yellow),
							wo.Label(`abc`))))
			}),
			wo.Example(`Horizontal and vertical together`, func() Sorm {
				return wo.Compound(
					wo.Limit(100, 100),
					wo.Compound(
						wo.Halign(0.5),
						wo.Valign(0.5),
						wo.Rectangle(-1, -1).Stroke(dark),
						wo.Compound(
							wo.Halign(0.5),
							wo.Valign(0),
							wo.Vshrink(),
							wo.Hshrink(),
							wo.Rectangle(-1, -1).Fill(yellow),
							wo.Label(`abc`))))
			}),
		),
		wo.Compound(
			wo.Hfollow(),
			wo.BetweenVoid(32, 0),
			wo.Example(`Stretch`, func() Sorm {
				return wo.Compound(
					wo.Limit(100, 100),
					wo.Compound(
						wo.Hfollow(),
						wo.Rectangle(-1, -1).Stroke(dark),
						wo.Rectangle(-2, -1).Stroke(dark),
						wo.Compound(
							// TODO Ways to determine size of a compound if all of
							// its members are stretchy.
							wo.Hfollow(),
							wo.Limit(20, 0),
							wo.Rectangle(-1, -1).Stroke(dark),
							wo.Rectangle(-1, -1).Stroke(dark),
						),
					))
			}),
			wo.Example(`Sliders`, func() Sorm {
				return wo.Compound(
					wo.Vfollow(),
					wo.BetweenVoid(0, 8),
					wo.Compound(
						wo.Limit(200, 20),
						wo.Slider(&slider0),
					),
					wo.Compound(
						wo.Limit(100, 20),
						wo.Slider(&slider1)))
			}),
			wo.Example(`Numbox`, func() Sorm {
				return wo.Compound(
					wo.Limit(100, 20),
					wo.Numbox(&numbox))
			}),
			wo.Compound(
				wo.Vfollow(),
				wo.BetweenVoid(0, 8),
				wo.Example(`File drop`, func() Sorm {
					return wo.Drop(&filename0, ``)
				}),
				wo.Example(`File drop, but only for MP3s`, func() Sorm {
					return wo.Drop(&filename1, `audio/mpeg`)
				}),
			),
			wo.Example(`Labels`, func() Sorm {
				return wo.Compound(
					wo.Vfollow(),
					wo.BetweenVoid(0, 15),
					wo.Compound(
						wo.Hfollow(),
						wo.BetweenVoid(10, 0),
						wo.Label(`Q`),
						wo.Label(`W`),
						wo.Label(`E`),
						wo.Label(`R`),
						wo.Label(`T`)),
					wo.Compound(
						wo.Hfollow(),
						wo.BetweenVoid(10, 0),
						wo.Label(`Y`),
						wo.Label(`U`),
						wo.Label(`I`),
						wo.Label(`O`),
						wo.Label(`P`),
						wo.Label(`x`),
					),
					wo.Compound(
						wo.Hfollow(),
						wo.BetweenVoid(10, 0),
						wo.Label(`A`),
						wo.Label(`H`),
						wo.Label(`AH`),
						wo.Label(`Ah`),
						wo.Label(`Ap`),
						wo.Label(`p`),
						wo.Label(`pp`),
					))
			}),
		),
		wo.Compound(
			wo.Vfollow(),
			wo.BetweenVoid(0, 64),
			wo.Example(`Scissor`, func() Sorm {
				return wo.Compound(
					wo.Limit(100, 100),
					wo.Scissor(),
					wo.Compound(
						wo.Vfollow(),
						wo.Rectangle(100, 100).Fill(dark),
						wo.Rectangle(100, 100).Fill(yellow)))
			}),
		))
}

func (wo *World) Drop(filename *string, mime string) Sorm {
	return wo.Compound(
		wo.Hshrink(),
		wo.Halign(0.5),
		wo.Rectangle(-1, 20).Fill(yellow),
		wo.Compound(
			wo.Hfollow(),
			wo.Valign(0.5),
			wo.Void(10, 20),
			wo.Label(*filename),
			wo.Void(10, 20),
		)).
		Cond(func(m contraption.Matcher) {
			pat := `Drop:in`
			if mime != `` {
				pat = `Drop(` + mime + `):in`
			}
			if m.Match(pat) {
				*filename = wo.Trace[0].E.(contraption.Drop).Paths[0]
			}
		})
}

func (wo *World) Numbox(v *float64) Sorm {
	return wo.Compound(
		wo.Halign(0.5),
		wo.Valign(0.5),
		wo.Rectangle(-1, 20).
			Fill(dark).
			Override(),
		wo.Compound(
			wo.Void(-1, -1),
			wo.Rectangle(complex(-*v, 0), -1).Fill(yellow)),
		wo.Void(-1+8i, -1+8i).Cond(func(m contraption.Matcher) {
			if m.Duration(300 * time.Millisecond).Match(`Click(1):in .* Unclick(1):in Click(1):in`) {
				*v = 0
			}
			if m.Match(`Click(1):in`) {
				wo.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
				*wo.Key(v) = wo.Trace[0].Pt.Y
			}
			if m.Match(`!Unclick(1)* Click(1):in`) {
				c := 0.001
				if m.Nochoke().Match(`!Release(Shift)* Press(Shift)`) {
					c = 0.0001
				}
				if *wo.Key(v) != nil {
					d := (*wo.Key(v)).(float64) - wo.Trace[0].Pt.Y
					*v += d * c
					*v = max(0, min(*v, 1))
				}
				*wo.Key(v) = wo.Trace[0].Pt.Y
			}
			if m.Nochoke().Match(`Unclick(1)`) {
				wo.Window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
			}
		}),
		wo.Label(strconv.FormatFloat(*v, 'f', 2, 64)))
}

func (wo *World) Slider(v *float64) Sorm {
	// Here I use drag-and-drop mechanics to implement Slider, but it can be implemented with Cond,
	// not using DragEffect.
	// Feedback is used in both cases.
	return wo.Compound(
		wo.Halign(*v),
		wo.Valign(0.5),
		wo.Void(-1, 20),
		wo.Rectangle(-1, 1).Fill(dark),
		wo.Identity(Tag(v, 1)),
		wo.Compound(
			// Center the knob relative to the rail, so it won't slide under the cursor.
			wo.Void(0, 0).Override(),
			wo.Halign(0.5),
			wo.Valign(0.5),
			wo.Compound(
				wo.Identity(v),
				wo.Source(),
				wo.Circle(20).Fill(yellow))))
}

func (wo *World) Label(v ...any) Sorm {
	// TODO wo.Text takes io.RuneReader, define and use SprintRuneReader
	return wo.Text(10, []rune(fmt.Sprint(v...))).Fill(hexpaint(`#000000`))
}

func (wo *World) Example(label string, ex func() Sorm) Sorm {
	return wo.Compound(
		wo.Vfollow(),
		ex(),
		wo.Void(0, 8),
		wo.Label(label),
	)
}

func Tag[T comparable](value T, i int) any {
	type tag struct {
		t T
		int
	}
	return tag{value, i}
}
