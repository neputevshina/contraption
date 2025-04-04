package main

import (
	"log"

	"github.com/neputevshina/contraption/nanovgo"
	"github.com/neputevshina/geom"
	"golang.org/x/exp/constraints"
)

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

var println = log.Println

func lerp(a, b, c float64) float64 {
	return a*(1-c) + b*c
}

func geom2nanovgo(g geom.Geom) nanovgo.TransformMatrix {
	return nanovgo.TransformMatrix{
		float32(g[0][0]), float32(g[0][1]),
		float32(g[1][0]), float32(g[1][1]),
		float32(g[2][0]), float32(g[2][1]),
	}
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
