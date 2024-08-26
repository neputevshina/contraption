package contraption

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/neputevshina/geom"
)

//go:generate go run github.com/pointlander/peg ./regexp.peg

const (
	rchar = iota
	rmatch
	rjmp
	rsplit
	rsave
	rjoker
	rnotchar
)

type rinst struct {
	opcode int

	e          any
	typeonly   bool
	in, out    bool
	begin, end bool

	x, y int
}

func (i *rinst) String() string {
	var in, begin, end string
	if i.in {
		in = " in"
	}
	if i.out {
		in = " out"
	}
	if i.begin {
		begin = " begin"
	}
	if i.end {
		end = " end"
	}
	return fmt.Sprint("{", map[int]string{
		rchar:    "rchar",
		rjmp:     "rjmp",
		rjoker:   "rjoker",
		rmatch:   "rmatch",
		rsplit:   "rsplit",
		rnotchar: "rnotchar",
	}[i.opcode], " ", i.x, " ", i.y, in, begin, end, " ", reflect.TypeOf(i.e).Name(), " (", i.e, ")}")
}

func rcompile(pattern string) ([]rinst, error) {
	var err error
	peg := &rpeg{}
	err = peg.Init(func(r *rpeg) error { r.Buffer = pattern; return nil })
	if err != nil {
		log.Panicln(err)
	}
	err = peg.Parse()
	if err != nil {
		return nil, err
	}
	//rprint(0, peg.AST(), pattern)
	//peg.AST().PrettyPrint(os.Stdout, pattern)
	v := rvisit(peg.AST().up, nil, peg.buffer)
	v = append(v, rinst{opcode: rmatch})
	return v, nil
}

type rlist struct {
	rinst
	next *rinst
}

func rprint(tab int, n *node32, pattern string) {
	if n == nil {
		return
	}

	fmt.Print(strings.Repeat("\t", tab))
	fmt.Println(pattern[n.begin:n.end], "("+n.token32.String()+")")
	rprint(tab+1, n.up, pattern)
	rprint(tab, n.next, pattern)
}

func rvisit(n, p *node32, pattern []rune) []rinst {
	if n == nil {
		return nil
	}
	adjust := func(r []rinst, n int) {
		for i := range r {
			r[i].x += n
			r[i].y += n
		}
	}

	switch n.token32.pegRule {
	case ruleAlter:
		left := rvisit(n.up, n, pattern)
		right := rvisit(n.up.next, n, pattern)
		v := []rinst{
			{opcode: rsplit,
				x: 1,
				y: 2 + len(left), // one for this, and one for next rjmp
			},
		}
		adjust(left, 1) // for previous rsplit
		v = append(v, left...)
		v = append(v, rinst{
			opcode: rjmp,
			x:      1 + len(left) + 1 + len(right),
		})
		adjust(right, 2+len(left)) // rsplit+left+rjmp
		v = append(v, right...)
		return v

	case ruleCat:
		left := rvisit(n.up, n, pattern)
		right := rvisit(n.up.next, n, pattern)
		adjust(right, len(left))
		return append(left, right...)

	case rulePoint:
		pt := rinst{}
		left := n.up.up // Concrete -> {} here
		typ := ""
		val := ""
		for left != nil {
			switch left.pegRule {
			case ruleRect:
				switch left.up.pegRule {
				case ruleOut:
					pt.out = true
				case ruleIn:
					pt.in = true
				}
			case ruleTime:
				switch left.up.pegRule {
				case ruleBegin:
					pt.begin = true
				case ruleEnd:
					pt.end = true
				}
			case ruleType:
				typ = string(pattern[left.begin:left.end])
			case ruleValue:
				val = string(pattern[left.begin:left.end])
			}
			left = left.next
		}
		if typ == "." {
			pt.opcode = rjoker
		} else {
			pt.opcode = rchar
			if typ[0] == '!' { // TODO make an operator
				typ = typ[1:]
				pt.opcode = rnotchar
			}
			pt.e = nameevent(typ, val)
			if val == "" {
				pt.typeonly = true
			}
		}
		left = n.up // Concrete
		if left.next != nil {
			switch left.next.pegRule {
			case ruleMaybe:
				v := []rinst{
					{opcode: rsplit,
						x: 1,
						y: 2,
					},
					pt,
				}
				return v

			case ruleAny:
				// TODO groups
				v := []rinst{
					{opcode: rsplit,
						x: 1,
						y: 3},
					pt,
					{opcode: rjmp,
						x: 0},
				}
				return v

			case ruleLazyAny:
				v := []rinst{
					{opcode: rsplit,
						x: 3,
						y: 1},
					pt,
					{opcode: rjmp,
						x: 0},
				}
				return v

			case ruleSeveral:
				panic("unimplemented, currently not needed")
			default:
				panic("check parser for errors")

			}
		}
		return []rinst{pt}
	}
	return []rinst{}
}

// rinterp is the threaded regular expression bytecode vm taken from https://swtch.com/~rsc/regexp/regexp2.html#thompsonvm.
func rinterp(program []rinst, trace []EventPoint, rect geom.Rectangle, dur time.Duration, deadline time.Time, z int) (bool, EventTraceLast) {
	type rthread = int
	cs := make([]rthread, 0, len(program))
	ns := make([]rthread, 0, len(program))
	var left, right time.Time
	box := geom.Rectangle{}

	cs = append(cs, 0)
	var prev EventPoint
	choked := false
	for j, sv := range trace {
		sv := sv
		for i := 0; i < len(cs); i++ {
			// j := j
			pc := cs[i]
			v := &program[pc]

			joker := func() {
				defer func() {
					if z <= 0 {
						return
					}
					trace[j].zold = trace[j].z
					trace[j].z = z
				}()
				// // This comparison is dependent on the Cond evaluation order in (*World).Develop.
				if v.in && !sv.Pt.In(rect) {
					return
				}
				if v.out && sv.Pt.In(rect) {
					return
				}
				// If left (ending) time is not specified first match will be it.
				if left == (time.Time{}) {
					left = sv.T
				}
				if z > 0 && trace[j].z > z {
					choked = true
					return
				}

				if v.begin {
					left = sv.T
				} else if v.end {
					right = sv.T
				}
				box.Min.X = min(box.Min.X, sv.Pt.X)
				box.Min.Y = min(box.Min.Y, sv.Pt.Y)
				box.Max.X = min(box.Max.X, sv.Pt.X)
				box.Max.Y = min(box.Max.Y, sv.Pt.Y)
				ns = append(ns, cs[i]+1)
				prev = sv
			}

			switch v.opcode {
			case rnotchar:
				if requals(sv, v) {
					break
				}
				joker()

			case rchar:
				if requals(sv, v) {
					joker()
				}

			case rjoker:
				joker()

			case rmatch:
				// If right (starting time) is not specified end of the match will be on its end, duh.
				if right == (time.Time{}) {
					right = prev.T
				}
				e := EventTraceLast{
					StartedAt:  right,
					Duration:   left.Sub(right),
					Box:        box,
					FirstTouch: sv.Pt,
				}
				// Make zold not reset if match was successful.
				defer func() {
					for j := range trace {
						if trace[j].z == z {
							trace[j].zold = z
						}
					}
				}()
				// Deadline is one who is changing there.
				return left.Sub(right) <= dur && deadline.Before(left), e

			case rjmp:
				cs = append(cs, v.x)

			case rsplit:
				cs = append(cs, v.x)
				cs = append(cs, v.y)
			}
		}
		cs, ns = ns, cs
		ns = ns[:0]
	}
	// Match was unsuccessful, reset orders.
	for j := range trace {
		trace[j].z = trace[j].zold
	}
	return false, EventTraceLast{Choked: choked}
}

type keyer interface {
	key() glfw.Key
}

// Moved to nanovgo+glfw.go
// func requals(p EventPoint, inst *rinst) bool {
// 	teq := reflect.TypeOf(p.E) == reflect.TypeOf(inst.e)
// 	if inst.typeonly {
// 		return teq
// 	}
// 	pk, ok1 := p.E.(keyer)    // Key in event
// 	ik, ok2 := inst.e.(keyer) // Key in regexp rule
// 	// If p.E or inst.e is not keyer, they both will fail on default condition.
// 	if ok1 && ok2 && teq {
// 		switch ik.key() {
// 		case anyShift:
// 			return pk.key() == glfw.KeyLeftShift ||
// 				pk.key() == glfw.KeyRightShift
// 		case anyCtrl:
// 			return pk.key() == glfw.KeyLeftControl ||
// 				pk.key() == glfw.KeyRightControl
// 		}
// 	}
// 	return p.E == inst.e
// }

/*
// pike's version of above, supports submatches https://swtch.com/~rsc/regexp/regexp2.html#pike
func rinterp(program []rinst, trace []EventPoint, rect image.Rectangle, dur time.Duration) (match bool, submatches [10][]EventPoint) {
	type rthread struct {
		pc    *rinst
		saved [10][]EventPoint
	}

	clist := make([]rthread, 0, len(program))
	nlist := make([]rthread, 0, len(program))
	var begin, end time.Time
	in := true

	clist = append(clist, rthread{&program[0], submatches})
	for si, sv := range trace {
		for i, t := range clist {
			switch t.pc.opcode {
			case rchar:
				if t.pc.begin {
					begin = sv.T
				} else if t.pc.end {
					end = sv.T
				}
				if t.pc.in {
					in = in && sv.Pt.In(rect)
				}
				if !requals(sv, t.pc) {
					break
				}
				clist = append(clist, rthread{&program[i+1], t.saved})

			case rmatch:
				submatches = t.saved
				return in && end.Sub(begin) <= dur, submatches

			case rjmp:
				clist = append(clist, rthread{t.pc.x, t.saved})

			case rsplit:
				clist = append(clist, rthread{t.pc.x, t.saved})
				clist = append(clist, rthread{t.pc.y, t.saved})

			case rsave:
				t.saved[t.pc.i] = trace[si:]
				clist = append(clist, rthread{t.pc.x, t.saved})
			}
		}
		clist, nlist = nlist, clist
		nlist = nlist[:0]
	}
	return false, submatches
}
*/
