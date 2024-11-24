package contraption

import (
	"math"

	"github.com/neputevshina/geom"
	"github.com/neputevshina/contraption/nanovgo"
	"golang.org/x/exp/slices"
)

const impEps = 0.1 // pixels

// impXpd returns partial derivative of an equation with respect to x.
func impXpd(eqn Eqn) (xpd Eqn) {
	return func(pt geom.Point) (dist float64) {
		a := eqn(geom.Pt(pt.X+impEps, pt.Y))
		b := eqn(geom.Pt(pt.X, pt.Y))
		return (a - b) / impEps
	}
}

// impXpd returns partial derivative of an equation with respect to y.
func impYpd(eqn Eqn) (xpd Eqn) {
	return func(pt geom.Point) (dist float64) {
		a := eqn(geom.Pt(pt.X, pt.Y+impEps))
		b := eqn(geom.Pt(pt.X, pt.Y))
		return (a - b) / impEps
	}
}

func impPoint(eqn, eqnx, eqny Eqn, near geom.Point) geom.Point {
	for range [10]struct{}{} { // Limit to 10 iterations.
		old := near
		// Saving this result is crucial here: compiler can't optimize this because of possible side-effects.
		sdfxnear := eqnx(near)
		sdfynear := eqny(near)
		dist := eqn(near) / (sdfxnear*sdfxnear + sdfynear*sdfynear)
		const impEpsCubed = impEps * impEps * impEps

		near = near.Sub(geom.Pt(sdfxnear, sdfynear).Mul(dist))
		if old.Sub(near).Length() <= impEps {
			break
		}
	}
	return near
}

type impLine struct {
	a, b geom.Point
}

// todo? https://github.com/prideout/par/blob/master/par_msquares.h
func impMarch(vg *nanovgo.Context, eqn Eqn, w, h float64) (points []geom.Point) {
	sign := math.Signbit
	const N = 22.0
	var _lines [N * N]impLine
	lines := _lines[:0:len(_lines)]

	// Step 1: sample the curve using marching squares algorithm.
	ptc := func(y, x int) geom.Point {
		fy, fx := float64(y), float64(x)
		return geom.Pt(w*fx/(N-1), h*fy/(N-1))
	}
	eqnx, eqny := impXpd(eqn), impYpd(eqn)
	var prevrow [N]float64
	for y := range [N]struct{}{} {
		var currrow [N]float64
		for x := range [N]struct{}{} {
			currrow[x] = eqn(ptc(y, x))
			if y > 0 && x > 0 {
				lt, ltd := func() geom.Point { return ptc(y-1, x-1) }, prevrow[x-1]
				rt, rtd := func() geom.Point { return ptc(y-1, x) }, prevrow[x]
				lb, lbd := func() geom.Point { return ptc(y, x-1) }, currrow[x-1]
				rb, rbd := func() geom.Point { return ptc(y, x) }, currrow[x]

				// TODO Fix interpolation
				// lm := func(p int) geom.Point { return ptlerp(lt(), lb(), clamp((1-lbd)/(ltd-lbd), 0, 1)) }
				// tc := func(p int) geom.Point { return ptlerp(lt(), rt(), clamp((1-ltd)/(rtd-ltd), 0, 1)) }
				// rm := func(p int) geom.Point { return ptlerp(rt(), rb(), clamp((1-rbd)/(rtd-rbd), 0, 1)) }
				// bc := func(p int) geom.Point { return ptlerp(lb(), rb(), clamp((1-lbd)/(rbd-lbd), 0, 1)) }

				// FIXME Temporary and costly solution: use Newton's step to get nearest value of a function.
				lm := func(p int) geom.Point { return impPoint(eqn, eqnx, eqny, ptlerp(lt(), lb(), 0.5)) }
				tc := func(p int) geom.Point { return impPoint(eqn, eqnx, eqny, ptlerp(lt(), rt(), 0.5)) }
				rm := func(p int) geom.Point { return impPoint(eqn, eqnx, eqny, ptlerp(rt(), rb(), 0.5)) }
				bc := func(p int) geom.Point { return impPoint(eqn, eqnx, eqny, ptlerp(lb(), rb(), 0.5)) }

				// apd := func(a, b geom.Point) { pairs = append(pairs, [2]geom.Point{a, b}) }
				apd := func(a, b geom.Point) {
					lines = append(lines, impLine{a: a, b: b})
					// vg.MoveTo(float32(a.X), float32(a.Y))
					// vg.LineTo(float32(b.X), float32(b.Y))
				}

				v := [4]bool{sign(ltd), sign(rtd), sign(lbd), sign(rbd)}
				switch v {
				case [4]bool{true, false, false, false}: // [' ]
					apd(tc(0), lm(0))
				case [4]bool{false, true, false, false}: // [ ']
					apd(tc(0), rm(0))
				case [4]bool{false, false, true, false}: // [. ]
					apd(lm(0), bc(0))
				case [4]bool{false, false, false, true}: // [ .]
					apd(bc(0), rm(0))

				case [4]bool{true, true, false, false}: // ['']
					apd(rm(0), lm(0))
				case [4]bool{false, false, true, true}: // [..]
					apd(rm(0), lm(0))
				case [4]bool{true, false, true, false}: // [: ]
					// println(ltd, rtd, lbd, rbd)
					// println(bc(0), tc(1))
					apd(bc(0), tc(0))
				case [4]bool{false, true, false, true}: // [ :]
					apd(tc(0), bc(0))
				case [4]bool{true, false, false, true}: // ['.]
					apd(lm(0), tc(0))
					apd(bc(0), rm(0))
				case [4]bool{false, true, true, false}: // [.']
					apd(lm(0), bc(0))
					apd(tc(0), rm(0))

				case [4]bool{false, true, true, true}: // [.:]
					apd(tc(0), lm(0))
				case [4]bool{true, false, true, true}: // [:.]
					apd(tc(0), rm(0))
				case [4]bool{true, true, false, true}: // [':]
					apd(lm(0), bc(0))
				case [4]bool{true, true, true, false}: // [:']
					apd(bc(0), rm(0))
				}
			}
		}
		prevrow = currrow
	}

	// Step 2: connect the lines.
	// TODO maybe it is better to open up shader of Nanovgo/use multiple
	// passes guided by z-buffer and draw meshes than connecting the lines
	// and making Nanovgo triangulate it again?
	var _nu [N * N]impLine
	nu := _nu[:0:len(_nu)]
	abs := math.Abs
	epseq := func(a, b geom.Point) bool { return abs(a.X-b.X) <= impEps && abs(a.Y-b.Y) <= impEps }

	for len(lines) > 0 {
		nu = append(nu, lines[len(lines)-1])
		lines = lines[:len(lines)-1]
		for !epseq(nu[len(nu)-1].b, nu[0].a) {
			end := nu[len(nu)-1].b
			for i, n := range lines {
				if epseq(end, n.b) {
					n.a, n.b = n.b, n.a
				}
				if epseq(end, n.a) {
					nu = append(nu, n)
					lines = slices.Delete(lines, i, i+1)
					break
				}
				if i == len(lines)-1 {
					panic(`unconnected lines`)
				}
			}
		}
	}

	points = append(points, nu[0].a)
	for _, p := range nu {
		points = append(points, p.b)
	}

	// vg.MoveTo(float32(nu[0].a.X), float32(nu[0].a.Y))
	// for _, n := range nu {
	// 	vg.LineTo(float32(n.b.X), float32(n.b.Y))
	// }
	return
}
