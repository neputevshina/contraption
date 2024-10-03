package contraption

import "unicode/utf8"

// Sequence is the thing that can generate elements for a scroll-enabled compound.
//
// All the logic to differentiate scroll's elements must be in returned Sorms.
type Sequence interface {
	Get(i int) Sorm
	Length() int
}

type adhocSequence struct {
	get    func(i int) Sorm
	length func() int
}

func (s *adhocSequence) Get(i int) Sorm {
	return s.get(i)
}

func (s *adhocSequence) Length() int {
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

type appendSequence struct {
	pre   bool
	affix []Sorm
	orig  Sequence
}

func (s *appendSequence) Get(i int) Sorm {
	if s.pre {
		al := len(s.affix)
		if i < al {
			return s.affix[i]
		} else {
			return s.orig.Get(al - i)
		}
	} else {
		al := s.orig.Length()
		if i < al {
			return s.orig.Get(i)
		} else {
			return s.affix[al-i]
		}
	}
}

func (s *appendSequence) Length() int {
	return s.orig.Length() + len(s.affix)
}

func PrependSeq(seq Sequence, wo *World, s ...Sorm) Sequence {
	return &appendSequence{
		pre:   true,
		affix: wo.stash(s),
		orig:  seq,
	}
}

func AppendSeq(seq Sequence, wo *World, s ...Sorm) Sequence {
	return &appendSequence{
		pre:   false,
		affix: wo.stash(s),
		orig:  seq,
	}
}

type stringSeq struct {
	string
	len, consumed int
	produce       func(rune) Sorm
}

func (s *stringSeq) Get(i int) Sorm {
	r, sz := utf8.DecodeRuneInString(s.string[s.consumed:])
	s.consumed += sz
	return s.produce(r)
}

func (s *stringSeq) Length() int {
	return s.len
}

func StringSeq(s string, produce func(rune) Sorm) Sequence {
	return &stringSeq{
		string:  s,
		produce: produce,
		len:     utf8.RuneCountInString(s),
	}
}
