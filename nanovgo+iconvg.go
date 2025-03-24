package contraption

import (
	"golang.org/x/exp/shiny/iconvg"
)

type NanovgoDestination struct {
	vg   *Context
	x, y float32
	iconvg.Palette
}

func (n *NanovgoDestination) Reset(m iconvg.Metadata) {
	n.vg.Reset()
	n.Palette = m.Palette
	n.vg.SetFillColor(hex(`#000000`))
}

func (n *NanovgoDestination) SetCSel(cSel uint8) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) SetNSel(nSel uint8) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) SetCReg(adj uint8, incr bool, c iconvg.Color) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) SetNReg(adj uint8, incr bool, f float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) SetLOD(lod0, lod1 float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) StartPath(adj uint8, x, y float32) {
	n.vg.BeginPath()
	n.vg.MoveTo(float64(x), float64(y))
}

func (n *NanovgoDestination) ClosePathEndPath() {
	n.vg.ClosePath()
}

func (n *NanovgoDestination) ClosePathAbsMoveTo(x, y float32) {
	n.x, n.y = x, y
	n.vg.ClosePath()
	n.vg.MoveTo(float64(n.x), float64(n.y))
}

func (n *NanovgoDestination) ClosePathRelMoveTo(x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsHLineTo(x float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) RelHLineTo(x float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsVLineTo(y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) RelVLineTo(y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsLineTo(x, y float32) {
	n.x, n.y = x, y
	n.vg.LineTo(float64(n.x), float64(n.y))
}

func (n *NanovgoDestination) RelLineTo(x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsSmoothQuadTo(x, y float32) {
	n.vg.QuadTo(float64(n.x), float64(n.y), float64(x), float64(y))
	n.x, n.y = x, y
}

func (n *NanovgoDestination) RelSmoothQuadTo(x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsQuadTo(x1, y1, x, y float32) {
	n.vg.QuadTo(float64(n.x), float64(n.y), float64(x), float64(y))
	n.x, n.y = x, y
}

func (n *NanovgoDestination) RelQuadTo(x1, y1, x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsSmoothCubeTo(x2, y2, x, y float32) {
	n.vg.BezierTo(float64(n.x), float64(n.y), float64(x2), float64(y2), float64(x), float64(y))
	n.x, n.y = x, y
}

func (n *NanovgoDestination) RelSmoothCubeTo(x2, y2, x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsCubeTo(x1, y1, x2, y2, x, y float32) {
	n.vg.BezierTo(float64(n.x), float64(n.y), float64(x2), float64(y2), float64(x), float64(y))
	n.x, n.y = x, y
}

func (n *NanovgoDestination) RelCubeTo(x1, y1, x2, y2, x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) RelArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32) {
	panic(`unimplemented`)
}

var _ iconvg.Destination = &NanovgoDestination{}
