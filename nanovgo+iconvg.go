package contraption

import (
	"github.com/neputevshina/contraption/nanovgo"
	"golang.org/x/exp/shiny/iconvg"
)

type NanovgoDestination struct {
	vg   *nanovgo.Context
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
	n.vg.MoveTo(x, y)
}

func (n *NanovgoDestination) ClosePathEndPath() {
	n.vg.ClosePath()
}

func (n *NanovgoDestination) ClosePathAbsMoveTo(x, y float32) {
	n.x, n.y = x, y
	n.vg.ClosePath()
	n.vg.MoveTo(n.x, n.y)
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
	n.vg.LineTo(n.x, n.y)
}

func (n *NanovgoDestination) RelLineTo(x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsSmoothQuadTo(x, y float32) {
	n.vg.QuadTo(n.x, n.y, x, y)
	n.x, n.y = x, y
}

func (n *NanovgoDestination) RelSmoothQuadTo(x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsQuadTo(x1, y1, x, y float32) {
	n.vg.QuadTo(n.x, n.y, x, y)
	n.x, n.y = x, y
}

func (n *NanovgoDestination) RelQuadTo(x1, y1, x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsSmoothCubeTo(x2, y2, x, y float32) {
	n.vg.BezierTo(n.x, n.y, x2, y2, x, y)
	n.x, n.y = x, y
}

func (n *NanovgoDestination) RelSmoothCubeTo(x2, y2, x, y float32) {
	panic(`unimplemented`)
}

func (n *NanovgoDestination) AbsCubeTo(x1, y1, x2, y2, x, y float32) {
	n.vg.BezierTo(n.x, n.y, x2, y2, x, y)
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
