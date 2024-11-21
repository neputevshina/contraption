package contraption

import (
	"github.com/neputevshina/geom"
	"github.com/neputevshina/nanovgo"
)

type Oscilloscope struct {
	on   bool
	text func(size float64, text string) Sorm
}

func (wo *World) displayOscilloscope() Sorm {
	if wo.Oscilloscope.on {
		return wo.Compound(
			// wo.Limit(wo.Wwin+1, wo.Hwin+1),
			wo.Compound(
				wo.Rectangle(complex(wo.Wwin+1, 0), complex(wo.Hwin+1, 0)).Fill(hexpaint(`#00000060`)), // FIXME Must be -1, -1, is not working even with limit.
				wo.Canvas(-1, -1, func(vgo *Context, _ geom.Geom, rect geom.Rectangle) {
					vgo.SetFillColor(hex(`#ffffff`))
					vgo.SetFontSize(12)
					vgo.SetTextAlign(nanovgo.AlignLeft)
					y := 0.0
					for _, s := range collect(wo.Events.Trace, func(e EventPoint) string { return sprint(e.String()) }) {
						vgo.Text(12, 14+y, s)
						y += 14
					}
				})).Cond(func(m Matcher) {
				println(m.Rect())
				// FIXME {0y0, 0y0}
				m.Match(`.*`)
			}))
	}
	return wo.Void(0, 0)
}
