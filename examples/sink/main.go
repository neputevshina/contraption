package main

import (
	"fmt"

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
	slider0 = 0.0
	slider1 = 1.0
)

func main() {
	wo := World{
		World: contraption.New(contraption.Config{}),
	}
	wo.Text = wo.NewVectorText(goregular.TTF)

	for wo.Next() {
		wo.Root(
			wo.Compound(
				wo.Halign(0.5),
				wo.Valign(0.5),
				wo.Void(-1, -1),
				wo.Examples(),
			))

		wo.Develop()
	}
}

func (wo *World) Examples() Sorm {

	wo.DragEffect = func(interval [2]geom.Point, drag any) Sorm {
		if f, ok := drag.(*float64); ok {
			r := wo.Prevkey(Tag(f, 1)).Rectangle()
			*f = (interval[1].X - r.Min.X) / r.Dx()
			*f = max(0, min(*f, 1))
		}
		return Sorm{}
	}

	return wo.Compound(
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
	)
}

func (wo *World) Slider(v *float64) Sorm {
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
