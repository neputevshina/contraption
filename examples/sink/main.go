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
				wo.Void(wo.Wwin, wo.Hwin),
				wo.Compound(
					wo.Hfollow(),
					wo.Valign(0),
					wo.Compound(
						wo.Limit(200, 200),
						wo.Compound(
							wo.Rectangle(-1, -1).Fill(hexpaint(`#dddd0080`)),
							wo.Compound(
								wo.Halign(0.5),
								wo.Valign(0.5),
								wo.Vshrink(),
								wo.Void(-1, 0).Voverride(),
								wo.Rectangle(-1, -1).Fill(hexpaint(`#ffff0080`)),
								wo.Text(16, []rune(`abc`)).Fill(hexpaint(`#000000`))))))))

		wo.Develop()
	}
}
