package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/neputevshina/contraption"
	"github.com/neputevshina/contraption/nanovgo"
	"github.com/neputevshina/geom"
)

type World struct {
	*contraption.World
}

var wo *World = &World{}

var patchpath = flag.String("p", "", "path to pure data patch")

func main() {
	flag.Parse()

	wo.World = contraption.New(contraption.Config{})

	pdcontext := wo.Vgo.Sub()
	f, err := os.Open(*patchpath)
	if err != nil {
		panic(err)
	}
	pdfile, err := io.ReadAll(f)
	generateVectors(pdfile, pdcontext)

	for wo.Next() {
		wo.Root(
			wo.Compound(wo.Canvas(800, 600, func(vgo *contraption.Context, wt geom.Geom, rect geom.Rectangle) {
				vgo.Replay(pdcontext)
			})))
		wo.Develop()
	}
}

func generateVectors(pdfile []byte, vgo *contraption.Context) {
	const cap = 10
	const height = cap + 4*2

	buf := bytes.NewBuffer(pdfile)
	i := -1
	closes := ""
	type ioconf struct {
		x, y, w, maxin, maxout int
	}
	conns := [][4]int{}
	iocs := []ioconf{}

	vgo.SetStrokeWidth(1)
	vgo.SetStrokeColor(hex(`#000000`))

	for {
		i++
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		bs := bytes.Split(line, []byte{' '})
		if string(bs[0]) == `#N` {
			if i == 0 {
				continue
			} else {
				closes = "restore"
			}
		}
		if closes != "" {
			if string(bs[1]) == closes {
				closes = ""
			} else {
				continue
			}
		}
		op := string(bs[1])
		if op == `obj` || op == `restore` || op == `msg` {
			x, _ := strconv.Atoi(string(bs[2]))
			y, _ := strconv.Atoi(string(bs[3]))

			te := bs[3:]
			l := len(te) - 1
			w := len(bytes.Join(te, []byte{' '}))
			if l > 2 && string(te[l-1]) == "f" &&
				bytes.HasSuffix(te[l-2], []byte(",")) {
				b, _ := bytes.CutSuffix(te[l], []byte(";\n"))
				w, err = strconv.Atoi(string(b))
			}
			w *= 5

			iocs = append(iocs, ioconf{x: x, y: y, w: w})

			vgo.BeginPath()
			vgo.Rect(float64(x)+.5, float64(y)+.5, float64(w), height)
			vgo.Stroke()
		}
		if string(bs[1]) == `connect` {
			// #X connect 2 1 17 0;
			// id 2 outlet 1 to id 17 inlet 0
			from, _ := strconv.Atoi(string(bs[2]))
			outlet, _ := strconv.Atoi(string(bs[3]))
			to, _ := strconv.Atoi(string(bs[4]))
			inlet, _ := strconv.Atoi(string(bs[5]))

			conns = append(conns, [4]int{from, outlet, to, inlet})

			iocs[from].maxout = max(iocs[from].maxout, outlet)
			iocs[to].maxin = max(iocs[from].maxin, inlet)
		}
	}

	for i, e := range iocs {
		println(i, e)
	}

	for _, c := range conns {
		vgo.BeginPath()
		// float64(x), float64(y), float64(w), height
		from := iocs[c[0]]
		to := iocs[c[2]]
		vgo.MoveTo(
			lerp(float64(from.x)+2+.5, float64(from.x+from.y)-2+.5,
				float64(from.maxout-c[1])/float64(from.maxout+1)),
			float64(from.y+height))
		vgo.LineTo(
			lerp(float64(to.x)+2, float64(to.x+to.y)-2,
				float64(to.maxin-c[3])/float64(to.maxin+1)),
			float64(to.y))
		println(float64(from.maxout-c[1])/float64(from.maxout+1), c[1], from.maxout)
		println(float64(to.maxin-c[3])/float64(to.maxin+1), c[3], from.maxin)
		println()
		vgo.Stroke()
	}

}

var println = fmt.Println

func lerp(a, b, x float64) float64 {
	return a*(1-x) + b*x
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
