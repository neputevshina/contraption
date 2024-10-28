package contraption

import "unicode/utf8"

// Sequence is the thing that can generate elements for a scroll-enabled compound.
//
// All the logic to differentiate scroll's elements must be in returned Sorms.
type Sequence interface {
	Get(wo *World, j int, buf []Sorm) (n int)
	Length(wo *World) int
}

type adhocSequence struct {
	get    func(j int) Sorm
	length func() int
}

func (s *adhocSequence) Get(wo *World, j int, buf []Sorm) (n int) {
	sl := s.length()
	for i := 0; i < min(len(buf), sl); i++ {
		buf[i] = s.get(j + i)
	}
	return sl - j
}

func (s *adhocSequence) Length(wo *World) int {
	return s.length()
}

func AdhocSequence(get func(i int) Sorm, length func() int) Sequence {
	return &adhocSequence{get: get, length: length}
}

func SliceSequence[T any](sl []T, produce func(T) Sorm) Sequence {
	return AdhocSequence(func(i int) Sorm { return produce(sl[i]) }, func() int { return len(sl) })
}

func SliceSequence2[T any](sl []T, produce func(int) Sorm) Sequence {
	return AdhocSequence(func(i int) Sorm { return produce(i) }, func() int { return len(sl) })
}

type Scrollptr struct {
	Index  int
	Offset float64
	y      float64

	Dirty bool
}

// type appendSequence struct {
// 	pre   bool
// 	affix []Sorm
// 	orig  Sequence
// }

// func (s *appendSequence) Get(wo *World, i int) Sorm {
// 	if s.pre {
// 		al := len(s.affix)
// 		if i < al {
// 			return s.affix[i]
// 		} else {
// 			return s.orig.Get(wo, al-i)
// 		}
// 	} else {
// 		al := s.orig.Length(wo)
// 		if i < al {
// 			return s.orig.Get(wo, i)
// 		} else {
// 			return s.affix[al-i]
// 		}
// 	}
// }

// func (s *appendSequence) Length(wo *World) int {
// 	return s.orig.Length(wo) + len(s.affix)
// }

// func PrependSeq(seq Sequence, wo *World, s ...Sorm) Sequence {
// 	return &appendSequence{
// 		pre:   true,
// 		affix: wo.stash(s),
// 		orig:  seq,
// 	}
// }

// func AppendSeq(seq Sequence, wo *World, s ...Sorm) Sequence {
// 	return &appendSequence{
// 		pre:   false,
// 		affix: wo.stash(s),
// 		orig:  seq,
// 	}
// }

type stringSeq struct {
	string
	len, consumed int
	produce       func(rune) Sorm
}

func (s *stringSeq) Get(wo *World, j int, buf []Sorm) (n int) {
	for i := 0; i < min(len(buf), s.len); i++ {
		r, sz := utf8.DecodeRuneInString(s.string[s.consumed:])
		buf[i] = s.produce(r)
		s.consumed += sz
	}
	return s.len - j
}

func (s *stringSeq) Length(wo *World) int {
	return s.len
}

func StringSeq(s string, produce func(rune) Sorm) Sequence {
	return &stringSeq{
		string:  s,
		produce: produce,
		len:     utf8.RuneCountInString(s),
	}
}
