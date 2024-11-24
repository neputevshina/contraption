package main

import "github.com/neputevshina/contraption/nanovgo"

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
