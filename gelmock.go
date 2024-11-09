package contraption

import (
	"fmt"
	"io"

	"github.com/neputevshina/geom"
)

// Segment is a segment of a vector path.
// It is based on golang.org/x/image/font/sfnt.Segment.
type Segment struct {
	// Op is the operator.
	// 	'M' — move, one coordinate
	// 	'L' — line, one coordinate
	//	'Q' — quadratic bezier, two coordinates
	// 	'C' — cubic bezier, three coordinates
	Op byte
	// Args is up to three (x, y) coordinates. The Y axis increases down.
	Args [3]geom.Point
}

func (s Segment) LastComponent() geom.Point {
	switch s.Op {
	case 'M':
		return s.Args[0]
	case 'L':
		return s.Args[0]
	case 'Q':
		return s.Args[1]
	case 'C':
		return s.Args[2]
	default:
		panic(`malformed Segment`)
	}
}

func (s Segment) String() string {
	switch s.Op {
	case 'M':
		return fmt.Sprint("{M ", s.Args[0], "}")
	case 'L':
		return fmt.Sprint("{L ", s.Args[0], "}")
	case 'Q':
		return fmt.Sprint("{Q ", s.Args[0], " ", s.Args[1], "}")
	case 'C':
		return fmt.Sprint("{Q ", s.Args[0], " ", s.Args[1], " ", s.Args[2], "}")
	default:
		return "{invalid op “" + string(s.Op) + "”}"
	}
}

type (
	point     = geom.Point
	rectangle = geom.Rectangle
)

var pt = geom.Pt

func validSize(r []rune) bool {
	return sizeRegexp.MatchReader(untoRuneScanner(r))
}

type runeBuf struct {
	buf []rune
	off int
}

func untoRuneScanner(b []rune) io.RuneScanner {
	return &runeBuf{buf: b}
}

func (s *runeBuf) ReadRune() (r rune, size int, err error) {
	if s.off >= len(s.buf) {
		err = io.EOF
		return
	}
	r = s.buf[s.off]
	s.off++
	size = 1
	err = nil
	return
}

func (s *runeBuf) UnreadRune() error {
	if s.off > 0 {
		s.off--
		return nil
	}
	panic("unread at 0")
}

func ApplySegment(g geom.Geom, s Segment) Segment {
	for i := range s.Args {
		s.Args[i] = g.ApplyPt(s.Args[i])
	}
	return s
}
