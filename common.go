package contraption

import (
	"regexp"

	"github.com/neputevshina/geom"
)

type bang = struct{}

// type PixelFace struct {
// 	font.Face
// 	// Cache values are preserved only between two consecutive frames.
// 	Rasterized map[rune]*struct {
// 		*gel.Texture
// 		Seen bool
// 	}
// }

// // Next is the new frame handler for PixelFace.
// func (pf PixelFace) Next() {
// 	for r, t := range pf.Rasterized {
// 		if t.Seen == false {
// 			delete(pf.Rasterized, r)
// 		}
// 	}
// }

// func (pf PixelFace) Rune(r rune) *gel.Texture {
// 	t, ok := pf.Rasterized[r]
// 	if !ok {
// 		// Almost copy-paste from unto/text.go
// 		dr, mask, mp, _, ok := pf.Face.Glyph(fixed.P(0, 0), r)
// 		if !ok {
// 			println("MISSING GLYPH", r)
// 		}
// 		w := int32(dr.Dx())
// 		h := int32(dr.Dy())
// 		// spaces and non-printable characters
// 		if w == 0 || h == 0 {
// 			return nil
// 		}
// 		su := image.NewRGBA(image.Rectangle{Max: dr.Size()})

// 		co := premultiply(hex(`#000000`))
// 		draw.DrawMask(su, su.Bounds(), &image.Uniform{co}, image.Point{}, mask, mp, draw.Over)
// 		te := gel.UploadUnfilteredTexture(su)
// 		t = &struct {
// 			*gel.Texture
// 			Seen bool
// 		}{Texture: te, Seen: true}
// 		pf.Rasterized[r] = t
// 	}
// 	t.Seen = true
// 	return t.Texture
// }

// func premultiply(straight unto.Color) (premultiplied color.RGBA) {
// 	premultiplied.R = (byte)(int(straight.R) * int(straight.A) / 255)
// 	premultiplied.G = (byte)(int(straight.G) * int(straight.A) / 255)
// 	premultiplied.B = (byte)(int(straight.B) * int(straight.A) / 255)
// 	premultiplied.A = straight.A
// 	return
// }

// // image/color.Color takes premultiplied color, so the image/draw.
// // formula from https://microsoft.github.io/Win2D/WinUI3/html/PremultipliedAlpha.htm
// var premulBlendMode = sdl.ComposeCustomBlendMode(
// 	sdl.BLENDFACTOR_ONE,
// 	sdl.BLENDFACTOR_ONE_MINUS_SRC_ALPHA,
// 	sdl.BLENDOPERATION_ADD,
// 	sdl.BLENDFACTOR_ZERO,
// 	sdl.BLENDFACTOR_ONE,
// 	sdl.BLENDOPERATION_ADD,
// )

var _ = func() {
}

var sizeRegexp = regexp.MustCompile(`[\pZ]*(\d+\.?\d*|\.\d+)(mm|pt|cm)[\pZ]*`)

func rectBox(r geom.Rectangle, thickness float64) (rs [4]geom.Rectangle) {
	for i := range rs {
		rs[i] = r
	}

	rs[0].Max.Y = rs[0].Min.Y
	rs[0].Min.Y -= thickness
	rs[0].Max.X += thickness

	rs[1].Min.X = rs[1].Max.X
	rs[1].Max.X += thickness
	rs[1].Max.Y += thickness

	rs[2].Min.Y = rs[2].Max.Y
	rs[2].Max.Y += thickness
	rs[2].Max.X += thickness

	rs[3].Max.X = rs[3].Min.X
	rs[3].Min.Y -= thickness
	rs[3].Min.X -= thickness

	return
}

func repeat[T any](count int, t T) (sl []T) {
	for i := 0; i < count; i++ {
		sl = append(sl, t)
	}
	return
}

type Sormer[T interface{ BaseWorld() *World }] interface {
	Sorm(T) Sorm
}
