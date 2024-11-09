package contraption

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/neputevshina/geom"
	"github.com/neputevshina/nanovgo"
	"golang.org/x/exp/constraints"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

const maxint = int(^uint(0) >> 1)

var (
	capmap        = map[int]float64{}
	testingBuffer sfnt.Buffer
)

func createFontFromMemory(vg *nanovgo.Context, name string, font []byte, capk float64) {
	id := vg.CreateFontFromMemory(name, font, 0)
	capmap[id] = capk
}

func setFontSize(vg *nanovgo.Context, size float64) {
	id := vg.FontFaceID()
	vg.SetFontSize(float32(size * capmap[id]))
}

func gray(v byte) (c nanovgo.Color) {
	c.R = float32(v) / 255.0
	c.G = float32(v) / 255.0
	c.B = float32(v) / 255.0
	c.A = 1
	return c
}

func hex(s string) (c nanovgo.Color) {
	if !(len(s) == 9 || len(s) == 7) || s[0] != '#' {
		panic("incorrect color")
	}
	nib := func(i int) uint8 {
		b := s[i]
		if b >= 'a' && b <= 'f' {
			return b - 'a' + 10
		}
		return b - '0'
	}
	c.R = float32(nib(1)<<4+nib(2)) / 255.0
	c.G = float32(nib(3)<<4+nib(4)) / 255.0
	c.B = float32(nib(5)<<4+nib(6)) / 255.0
	c.A = 1
	if len(s) == 9 {
		c.A = float32(nib(7)<<4+nib(8)) / 255.0
	}
	return c
}

func paint(c nanovgo.Color) nanovgo.Paint { return nanovgo.LinearGradient(0, 0, 1, 1, c, c) }

func hexpaint(s string) nanovgo.Paint { return paint(hex(s)) }

// func println(args ...any) {
// 	_ = log.Println
// 	buf := [2 << 10]byte{}
// 	n := runtime.Stack(buf[:], false)
// 	log.Println(string(buf[:n]))
// }

var println = log.Println
var print = fmt.Print
var sprint = fmt.Sprint

var typeof = reflect.TypeOf

func lighter(col nanovgo.Color) nanovgo.Color {
	col.R += 0.1
	col.G += 0.1
	col.B += 0.1
	return col
}

func umatchclick(u *Events, efrect geom.Rectangle) (bool, geom.Point) {
	_ = `!Unclick(1)* Click(1):in`
	for _, p := range u.Trace {
		switch v := p.E.(type) {
		case Click:
			return v == 1 && p.Pt.In(efrect), p.Pt
		case Unclick:
			return v != 1, geom.Point{}
		}
	}
	return false, geom.Point{}
}

func umatchclickend(u *Events) bool {
	for _, p := range u.Trace {
		switch v := p.E.(type) {
		case Unclick:
			return v == 1
		case Click:
			return v != 1
		}
	}
	return false
}

func GetSystemFont(err error) (string, float64, error) {
	if err != nil {
		return "", 0, err
	}
	desktop := os.Getenv(`XDG_CURRENT_DESKTOP`)
	home, err := os.UserHomeDir()
	if err != nil {
		return "", 0, err
	}

	fgtk2rc, gtk2err := os.Open(path.Join(home, `.gtkrc2.0`))
	gtk2rc, gtk2err := io.ReadAll(fgtk2rc)
	switch desktop {
	case `LXDE`:
		if gtk2err != nil {
			// TODO Try another way, if you know it
			return "", 0, err
		}
		m := regexp.MustCompile(`gtk-font-name = "(.*?)(\d+)"\n`).FindSubmatch(gtk2rc)
		f, err := strconv.ParseFloat(string(m[2]), 64)
		return string(m[1]), f, err
	}
	return "", 0, err
}

func elvis[T any](p *T) *T {
	if p == nil {
		var nu T
		return &nu
	}
	return p
}

func max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func as[T any](a any) T {
	v, _ := a.(T)
	return v
}

func zero[T any](v T) bool {
	var z T
	sz := unsafe.Sizeof(v)
	mem := func(p *T) []byte {
		return unsafe.Slice((*byte)(unsafe.Pointer(p)), sz)
	}
	return bytes.Equal(mem(&v), mem(&z))
}

// arif returns a sum of numbers starting from from and ending with to with increment inc.
func arif[T number](from, to, inc T) T {
	return (from + to) * ((to-from)/inc + 1) / 2
}

func bigintmin(a, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return a
	}
	return b
}

func bigintmax(a, b *big.Int) *big.Int {
	if a.Cmp(b) > 0 {
		return a
	}
	return b
}

func clamp[T constraints.Ordered](low, a, high T) T {
	return max(low, min(a, high))
}

type float interface {
	~float64 | ~float32
}

type number interface {
	~int | ~uint | float
}

func lerp[F float](a, b, c F) F {
	return a*(1-c) + b*c
}

func abs[N number](n N) (z N) {
	return max(n, -n)
}

func sign[N number](n N) (z N) {
	return n / abs(n)
}

func pcos(x float64) float64 {
	return (math.Cos(x) + 1) / 2
}

func mix(a, b nanovgo.Color, c float32) nanovgo.Color {
	m := lerp[float32]
	a.R = m(a.R, b.R, c)
	a.G = m(a.G, b.G, c)
	a.B = m(a.B, b.B, c)
	a.A = m(a.A, b.A, 0.5)
	return a
}

func ptlerpx(a, b geom.Point, c float64) geom.Point {
	x := a.X*(1-c) + b.X*c
	a.X = x
	return a
}

func ptlerpy(a, b geom.Point, c float64) geom.Point {
	y := a.Y*(1-c) + b.Y*c
	a.Y = y
	return a
}

func ptlerp(a, b geom.Point, c float64) geom.Point {
	return a.Mul(1 - c).Add(b.Mul(c))
}

func dup[T any](a T) *T {
	return &a
}

func geom2nanovgo(g geom.Geom) nanovgo.TransformMatrix {
	return nanovgo.TransformMatrix{
		float32(g[0][0]), float32(g[0][1]),
		float32(g[1][0]), float32(g[1][1]),
		float32(g[2][0]), float32(g[2][1]),
	}
}

func probe[T any](t T) T {
	_, f, l, _ := runtime.Caller(1)
	println(sprint(f, ":", l, " ", t))
	return t
}

func probes(t ...any) {
	_, f, l, _ := runtime.Caller(1)
	println(sprint(append([]any{f, ":", l, " "}, t...)...))
}

func cond[T any](p bool, a T, b T) T {
	if p {
		return a
	} else {
		return b
	}
}

func last[T any](s []T) *T {
	return &s[len(s)-1]
}

func fixeddiv(p, q fixed.Int26_6) fixed.Int26_6 {
	return fixed.Int26_6(int64(p) << 6 / int64(q))
}

func collect[A, B any](over []A, f func(a A) B) []B {
	new := make([]B, len(over))
	for i := range over {
		new[i] = f(over[i])
	}
	return new
}

func fold[T, J any](over []T, initial J, f func(a J, b T) (c J)) J {
	out := initial
	for i := range over {
		out = f(out, over[i])
	}
	return out
}

func fold0[T, J any](over []T, f func(a J, b T) (c J)) J {
	var out J
	for i := range over {
		out = f(out, over[i])
	}
	return out
}

func fold1[T any](over []T, f func(a T, b T) (c T)) T {
	out := over[0]
	for i := range over {
		out = f(out, over[i])
	}
	return out
}

func sum[T constraints.Float](s []T) T {
	return fold0(s, func(a, b T) T { return a + b })
}

func andfold[T any](s []T, f func(T) bool) bool {
	return fold0(s, func(b bool, a T) bool { return f(a) && b })
}

func orfold[T any](s []T, f func(T) bool) bool {
	return fold0(s, func(b bool, a T) bool { return f(a) || b })
}

func filter[A any](over []A, f func(a A) bool) []A {
	new := make([]A, 0, len(over))
	for _, a := range over {
		if f(a) {
			new = append(new, a)
		}
	}
	return new
}

type valuet[T comparable] struct {
	v T
}

func value[T comparable](v T) valuet[T] { return valuet[T]{v: v} }

func (v valuet[T]) oneof(alt ...T) bool {
	for _, n := range alt {
		if v.v == n {
			return true
		}
	}
	return false
}

func rectpt(pt geom.Point) geom.Rectangle {
	return geom.Rectangle{Min: pt, Max: pt}
}

func runelen(s string) int {
	return len([]rune(s))
}

// qtoc converts a quadratic Bezier segment to a cubic one.
// Start is the previous pen position.
// Needed for PDF export.
func qtoc(q Segment, start geom.Point) (c Segment) {
	cp1 := start.Add(q.Args[0].Sub(start).Mul(2.0 / 3))
	cp2 := q.Args[1].Add(q.Args[0].Sub(q.Args[1]).Mul(2.0 / 3))
	return Segment{
		Op:   'C',
		Args: [3]geom.Point{cp1, cp2, q.Args[1]},
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type Symbol int

const (
	Silence Symbol = iota
	Unfocus
	Ack
)

func mmtopt(mm float64) float64 {
	return mm * 254 / 720
}

func pttomm(pt float64) float64 {
	return pt * 720 / 254
}

func cceil(c complex128) complex128 {
	return complex(math.Ceil(real(c)), imag(c))
}

func rect2nvgxywh(r geom.Rectangle) (x, y, w, h float32) {
	x = float32(r.Min.X)
	y = float32(r.Min.Y)
	w = float32(r.Dx())
	h = float32(r.Dy())
	return
}

func sameslice[T any](a, b []T) bool {
	switch {
	case cap(a) == 0:
		return false
	case cap(b) == 0:
		return false
	case &a[:1][0] == &b[:1][0]:
		return true
	default:
		return false
	}
}

func oneof[T comparable](a T, as ...T) bool {
	for _, v := range as {
		if v == a {
			return true
		}
	}
	return false
}

func zeroandclear[T any](pool *[]T) {
	var z T
	for i := range *pool {
		(*pool)[i] = z
	}
	(*pool) = (*pool)[:0]
}

func roundmodf(f float64) (float64, float64) {
	r := math.Round(f)
	return f - r, r
}
