package nanovgo

import (
	"github.com/neputevshina/contraption"
	"github.com/neputevshina/contraption/nanovgo"
	"github.com/neputevshina/contraption/op"
)

func New(f nanovgo.CreateFlags) *Renderer {
	rer := &Renderer{}
	var err error
	rer.Context, err = nanovgo.NewContext(f)
	if err != nil {
		panic(err)
	}
	return rer
}

type Renderer struct {
	*nanovgo.Context
}

func (rer *Renderer) Run(c *contraption.Context) {
	vgo := rer.Context
	vtxb := []nanovgo.Vertex{}
	for i := range c.Log {
		l := &c.Log[i]
		switch l.Tag {
		case op.Also:
			vgo.Also()
		case op.Arc:
			vgo.Arc(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]), float32(l.Args[4]), l.Direction)
		case op.ArcTo:
			vgo.ArcTo(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]), float32(l.Args[4]))
		case op.BeginFrame:
			vgo.BeginFrame(l.Iargs[0], l.Iargs[1], float32(l.Args[0]))
		case op.BeginPath:
			vgo.BeginPath()
		case op.BezierTo:
			vgo.BezierTo(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]), float32(l.Args[4]), float32(l.Args[5]))
		case op.Block:
			panic(`unimplemented`)
		case op.CancelFrame:
			vgo.CancelFrame()
		case op.Circle:
			vgo.Circle(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]))
		case op.ClosePath:
			vgo.ClosePath()
		case op.CreateFontFromMemory:
			f := &c.Fonts[l.Hfont]
			f.Hbackend = vgo.CreateFontFromMemory("", f.Data, f.FreeData)
		case op.CreateImageFromGoImage:
			m := &c.Images[l.Himage]
			m.H = vgo.CreateImageFromGoImage(m.ImageFlags, m.Image)
		case op.CreateImageRGBA:
			m := &c.Images[l.Himage]
			m.H = vgo.CreateImageRGBA(m.Wh.X, m.Wh.Y, m.ImageFlags, m.Data)
		case op.CurrentTransform:
			panic(`getter, unreachable`)
		case op.DebugDumpPathCache:
			vgo.DebugDumpPathCache()
		case op.Delete:
			vgo.Delete()
		case op.DeleteImage:
			vgo.DeleteImage(c.Images[l.Himage].H)
		case op.Ellipse:
			vgo.Ellipse(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]))
		case op.EndFrame:
			vgo.EndFrame()
		case op.Fill:
			vgo.Fill()
		case op.FindFont:
			panic(`unimplemented`)
		case op.FontBlur:
			panic(`getter, unreachable`)
		case op.FontFace:
			panic(`getter, unreachable`)
		case op.FontFaceID:
			panic(`getter, unreachable`)
		case op.FontSize:
			panic(`getter, unreachable`)
		case op.GlobalAlpha:
			panic(`getter, unreachable`)
		case op.ImageSize:
			panic(`getter, unreachable`)
		case op.IntersectScissor:
			vgo.IntersectScissor(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]))
		case op.LineCap:
			panic(`getter, unreachable`)
		case op.LineJoin:
			panic(`getter, unreachable`)
		case op.LineTo:
			vgo.LineTo(float32(l.Args[0]), float32(l.Args[1]))
		case op.MiterLimit:
			panic(`getter, unreachable`)
		case op.MoveTo:
			vgo.MoveTo(float32(l.Args[0]), float32(l.Args[1]))
		case op.PathWinding:
			vgo.PathWinding(l.Winding)
		case op.QuadTo:
			vgo.QuadTo(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]))
		case op.Rect:
			vgo.Rect(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]))
		case op.Reset:
			vgo.Reset()
		case op.ResetScissor:
			vgo.ResetScissor()
		case op.ResetTransform:
			vgo.ResetTransform()
		case op.Restore:
			vgo.Restore()
		case op.Rotate:
			vgo.Rotate(float32(l.Args[0]))
		case op.RoundedRect:
			vgo.RoundedRect(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]), float32(l.Args[4]))
		case op.Save:
			vgo.Save()
		case op.Scale:
			vgo.Scale(float32(l.Args[0]), float32(l.Args[1]))
		case op.Scissor:
			vgo.Scissor(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]))
		case op.SetFillColor:
			vgo.SetFillColor(l.Fillc)
		case op.SetFillPaint:
			p := l.Fillp
			p.Image = c.Images[p.Image].H
			vgo.SetFillPaint(p)
		case op.SetFontBlur:
			vgo.SetFontBlur(float32(l.Args[0]))
		case op.SetFontFace:
			panic(`unimplemented`)
		case op.SetFontFaceID:
			vgo.SetFontFaceID(l.Hfont)
		case op.SetFontSize:
			vgo.SetFontSize(float32(l.Args[0]))
		case op.SetGlobalAlpha:
			vgo.SetGlobalAlpha(float32(l.Args[0]))
		case op.SetLineCap:
			vgo.SetLineCap(l.Linecap)
		case op.SetLineJoin:
			vgo.SetLineJoin(l.Linecap)
		case op.SetMiterLimit:
			vgo.SetMiterLimit(float32(l.Args[0]))
		case op.SetStrokeColor:
			vgo.SetStrokeColor(l.Strokec)
		case op.SetStrokePaint:
			p := l.Strokep
			p.Image = c.Images[p.Image].H
			vgo.SetStrokePaint(p)
		case op.SetStrokeWidth:
			vgo.SetStrokeWidth(float32(l.Strokew))
		case op.SetTextAlign:
			vgo.SetTextAlign(l.Align)
		case op.SetTextLetterSpacing:
			vgo.SetTextLetterSpacing(float32(l.Args[0]))
		case op.SetTextLineHeight:
			vgo.SetTextLineHeight(float32(l.Args[0]))
		case op.SetTransform:
			vgo.SetTransform(l.TransformMatrix)
		case op.SetTransformByValue:
			vgo.SetTransformByValue(float32(l.Args[0]), float32(l.Args[1]), float32(l.Args[2]), float32(l.Args[3]), float32(l.Args[4]), float32(l.Args[5]))
		case op.SkewX:
			vgo.SkewX(float32(l.Args[0]))
		case op.SkewY:
			vgo.SkewY(float32(l.Args[0]))
		case op.Stroke:
			vgo.Stroke()
		case op.StrokeWidth:
			vgo.StrokeWidth()
		case op.TextAlign:
			panic(`getter, unreachable`)
		case op.TextBounds:
			panic(`getter, unreachable`)
		case op.TextLetterSpacing:
			panic(`getter, unreachable`)
		case op.TextLineHeight:
			panic(`getter, unreachable`)
		case op.TextMetrics:
			panic(`getter, unreachable`)
		case op.TextRune:
			// vgo.TextRune(c.fs, float32(l.Args[1]), float32(l.Args[2]), l.Runes)

			sus := c.SpriteUnits[l.Left:l.Right]
			bf := c.Fonts[l.Hfont]

			vtxb = vtxb[:0]
			vtxb = append(vtxb, make([]nanovgo.Vertex, 4*len(sus))...)
			vidx := 0
			invScale := l.Args[0]

			if l.Args[3] == 1 {
				// Reallocate atlas since we have new glyphs.
				_, ok := vgo.AllocTextAtlas(c.Fs, bf.Hbackend)
				if !ok {
					panic(``)
				}
			}
			for i := range sus {
				quad := sus[i]
				// Transform corners.
				t := vgo.CurrentTransform()
				c0, c1 := t.TransformPoint(float32(quad.Clip.Min.X*invScale), float32(quad.Clip.Min.Y*invScale))
				c2, c3 := t.TransformPoint(float32(quad.Clip.Max.X*invScale), float32(quad.Clip.Min.Y*invScale))
				c4, c5 := t.TransformPoint(float32(quad.Clip.Max.X*invScale), float32(quad.Clip.Max.Y*invScale))
				c6, c7 := t.TransformPoint(float32(quad.Clip.Min.X*invScale), float32(quad.Clip.Max.Y*invScale))
				//log.Printf("quad(%c) x0=%d, x1=%d, y0=%d, y1=%d, s0=%d, s1=%d, t0=%d, t1=%d\n", iter.CodePoint, int(quad.Clip.Min.X), int(quad.Clip.Max.X), int(quad.Clip.Min.Y), int(quad.Clip.Max.Y), int(1024*quad.Tc.Min.X), int(quad.Tc.Max.X*1024), int(quad.Tc.Min.Y*1024), int(quad.Tc.Max.Y*1024))
				// Create triangles
				vtx := func(x, y float32, u, v float64) nanovgo.Vertex {
					return nanovgo.Vertex{X: float32(x), Y: float32(y), U: float32(u), V: float32(v)}
				}
				vtxb[vidx] = vtx(c2, c3, quad.Tc.Max.X, quad.Tc.Min.Y)
				vtxb[vidx+1] = vtx(c0, c1, quad.Tc.Min.X, quad.Tc.Min.Y)
				vtxb[vidx+2] = vtx(c4, c5, quad.Tc.Max.X, quad.Tc.Max.Y)
				vtxb[vidx+3] = vtx(c6, c7, quad.Tc.Min.X, quad.Tc.Max.Y)
				vidx += 4
			}
			vgo.FlushTextTexture(c.Fs, bf.Hbackend)
			vgo.RenderText(vtxb)

		case op.Translate:
			vgo.Translate(float32(l.Args[0]), float32(l.Args[1]))
		case op.UpdateImage:
			f := c.Images[l.Himage]
			vgo.UpdateImage(l.Himage, f.Data)
		default:
			panic(`unreachable`)
		}
	}
}
