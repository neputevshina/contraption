package main

import (
	"github.com/neputevshina/contraption"
	"golang.org/x/image/font/gofont/goregular"
)

type Sorm = contraption.Sorm
type World struct {
	contraption.World
	Text func(size float64, str []rune) Sorm
}

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
	yellow := hexpaint(`#ffdb2ca0`)
	dark := hexpaint(`#f66b0080`)
	return wo.Compound(
		wo.Hfollow(),
		wo.BetweenVoid(32, 0),
		wo.Example(`Align child's center to container's top`, func() Sorm {
			return wo.Compound(
				wo.Limit(100, 100),
				wo.Compound(
					wo.Rectangle(-1, -1).Fill(dark),
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
					wo.Rectangle(-1, -1).Fill(dark),
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
					wo.Rectangle(-1, -1).Fill(dark),
					wo.Compound(
						wo.Halign(0.5),
						wo.Valign(0),
						wo.Vshrink(),
						wo.Hshrink(),
						wo.Rectangle(-1, -1).Fill(yellow),
						wo.Label(`abc`))))
		}),
	)
}

func (wo *World) Label(s string) Sorm {
	return wo.Text(10, []rune(s)).Fill(hexpaint(`#000000`))
}

func (wo *World) Example(label string, ex func() Sorm) Sorm {
	return wo.Compound(
		wo.Vfollow(),
		ex(),
		wo.Void(0, 8),
		wo.Label(label),
	)
}
