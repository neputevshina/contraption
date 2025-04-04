package contraption

import (
	"encoding/gob"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/neputevshina/geom"
)

// Event trace buffer length
const tracelen = 128

var endOfTimes = time.Unix(1<<63-62135596801, 999999999)

type EventTraceLast struct {
	StartedAt  time.Time // Time of a first matched event.
	FirstTouch geom.Point
	Duration   time.Duration  // Duration of match.
	Freshness  time.Duration  // Difference between start of match and now.
	Box        geom.Rectangle // Total cursor position distribution of a match.
	Choked     bool
}

type Events struct {
	Trace    []EventPoint
	Last     EventTraceLast
	Now      time.Time
	Dt       time.Duration
	Viewport geom.Point

	deadline time.Time

	tr         [tracelen * 2]EventPoint
	trcur      int
	held       [tracelen / 2]EventPoint
	heldcur    int
	evuser     [tracelen]EventPoint
	temp       [tracelen / 4]EventPoint
	tempcur    int
	regexps    map[string][]rinst
	MatchCount int

	// 0 — normal operation
	// 1 — recording
	// 2 — replaying
	// 3 — end of replay
	rec        int
	rec0       time.Time
	playc      chan time.Time
	RecordPath string
	records    []EventPoint
	future     EventPoint
}

// scrollPush pushes a new event point to the circular trace buffer and
// returns a slice of it where the new value is the first element.
func scrollPush(buffer []EventPoint, cursor *int, v EventPoint) (slice []EventPoint) {
	gl := len(buffer)
	buffer[(*cursor+gl/2)%gl] = v
	buffer[*cursor] = v
	*cursor--
	if *cursor < 0 {
		*cursor = gl/2 - 1
	}

	// *cursor+1 because in other case it will repeat event triggers after period of a buffer.
	return buffer[*cursor+1 : *cursor+gl/2]
}

func (u *Events) SetDeadline(t time.Time) {
	u.deadline = t
}

func (wo *Events) next() bool {
	now := time.Now()
	if wo.tempcur == 0 {
		// Don't update time on catching up events or else they won't be matched as fresh.
		wo.Dt = wo.Now.Sub(now)
		wo.Now = now // Current frame events are “from future” so they are fresh if deadline is “now”.
	}
	if wo.rec >= 2 {
		now = <-wo.playc
		// We need to poll events to get the window close event.
		// If it is not done, window can't be closed in replay mode.
		concretepoll(wo)
	}

	if wo.rec < 2 {
		if wo.tempcur > 0 {
			m := wo.temp[0]
			copy(wo.temp[:], wo.temp[1:])
			wo.trueemit(m.E, m.Pt, m.T)
			wo.tempcur--
		} else {
			if wo.Now.Compare(wo.deadline) >= 0 {
				concretewait(wo)
			} else {
				concretepoll(wo)
			}
		}
	}

	return true
}

func (u *Events) emit(ev interface{}, pt geom.Point, t time.Time) {
	m := EventPoint{ev, pt, t, 0, 0}
	// Skip events if application is lagging.
	u.tempcur = min(len(u.temp)-1, u.tempcur)
	u.temp[u.tempcur] = m
	u.tempcur++
}

func (u *Events) develop() {
	for i := range u.tr {
		u.tr[i].z = 0
		u.tr[i].zc = 0
	}
}

// trueemit pushes the new event to the trace.
func (u *Events) trueemit(ev interface{}, pt geom.Point, t time.Time) {
	if _, yes := ev.(EventPoint); yes {
		panic("can't emit EventPoint")
	}
	m := EventPoint{ev, pt, t, 0, 0}

	if u.rec == 1 {
		u.records = append(u.records, m)
	}

	// // Push paint deadline further.
	// u.SetDeadline(t.Add(100 * time.Millisecond))

	// If the just happened event is the same type and value as the latest event in Trace,
	// the latter is pushed to Details and replaced with former.
	// This behavior may be not the same for types of events that might be added in future.
	if ev == u.Trace[0].E && pt == u.Trace[0].Pt {
		u.Trace[0] = m
		return
	}

	base := scrollPush(u.tr[:], &u.trcur, m)
	relative := len(base) - u.heldcur

	// If our click is unclicked we can obviously stop holding it.
	// No matter in which place it is held.
	for i := range u.held[:u.heldcur] {
		if u.held[i].E == complement(base[relative].E) {
			copy(u.held[i:], u.held[i+1:])
			u.heldcur--
			// We can't have two same held events: a button can't be held twice.
			break
		}
	}

	// If the holdable type is in the bottom of the trace, it is copied to the hold trace,
	// which is then copied over the end of Trace.
	// This way user can accurately match, for instance, any action that begins with a press or
	// click and ends several other presses/clicks later.
	if bottom := base[relative-1]; holdable(bottom.E) {
		if u.heldcur < len(u.held) {
			copy(u.held[1:], u.held[:])
			u.held[0] = bottom
			u.heldcur++
		}
	}

	// FIXME -1 because rinterp needs a tail to get rmatch
	relative = len(u.Trace) - u.heldcur - 1
	u.Trace = u.evuser[:]
	copy(u.Trace, base)
	copy(u.Trace[relative:], u.held[:u.heldcur])
	u.Trace[len(u.Trace)-1] = EventPoint{}
	// println(u.Trace)
}

func (u *Events) In(r geom.Rectangle) Matcher {
	m := newMatcher(u)
	m.rect = r
	m.alwaysin = true
	return m
}

func (u *Events) Anywhere() Matcher {
	m := newMatcher(u)
	m.rect = geom.Rect(-123456, -132412, 142323, 123111)
	m.alwaysin = false
	return m
}

//syntax:type contraption-regexp
type Regexp string

// Matcher is a builder interface for regular expressions.
type Matcher struct {
	u        *Events
	pattern  string
	rect     geom.Rectangle
	dur      time.Duration
	deadline time.Time
	z        int
	alwaysin bool
}

func newMatcher(e *Events) Matcher {
	return Matcher{
		u:        e,
		rect:     geom.Rect(-99999, -99999, 99999, 99999),
		dur:      time.Duration(^uint64(0) >> 1),
		deadline: e.Now,
		z:        0,
	}
}

func (m Matcher) Rect() geom.Rectangle {
	return m.rect
}

func (m Matcher) Z() int {
	return m.z
}

func (m Matcher) WithZ(z int) Matcher {
	m.z = z
	return m
}

func (m Matcher) Duration(d time.Duration) Matcher {
	m.dur = d
	return m
}

func (m Matcher) Deadline(t time.Time) Matcher {
	m.deadline = t
	return m
}

func (m Matcher) Nochoke() Matcher {
	m.z = 0
	return m
}

func (m Matcher) Anywhere() Matcher {
	m.alwaysin = false
	return m
}

func (m Matcher) Indef() Matcher {
	m.deadline = time.Time{}
	return m
}

func (m Matcher) Match(pattern Regexp) bool {
	return m.u.match(pattern, m.rect, m.dur, m.deadline, m.z, m.alwaysin)
}

func (u *Events) Match(pattern Regexp) bool {
	return u.match(pattern, geom.Rect(-99999, -99999, 99999, 99999), time.Duration(^uint64(0)>>1), u.Now, 0, false)
}

// Hint: :in. And add MatchAllIn later, for fuck's sake.
func (u *Events) MatchIn(pattern Regexp, r geom.Rectangle) bool {
	return u.match(pattern, r, time.Duration(^uint64(0)>>1), u.Now, 0, false)
}

func (u *Events) MatchInNochoke(pattern Regexp, r geom.Rectangle) bool {
	// TODO
	return u.match(pattern, r, time.Duration(^uint64(0)>>1), u.Now, maxint, false)
}

func (u *Events) MatchIndef(pattern Regexp) bool {
	return u.match(pattern, geom.Rect(-99999, -99999, 99999, 99999), time.Duration(^uint64(0)>>1), time.Time{}, 0, false)
}

func (u *Events) MatchInIndef(pattern Regexp, rect geom.Rectangle) bool {
	return u.match(pattern, rect, time.Duration(^uint64(0)>>1), time.Time{}, 0, false)
}

func (u *Events) MatchInFreshness(pattern Regexp, rect geom.Rectangle, freshness time.Duration) bool {
	return u.match(pattern, rect, time.Duration(^uint64(0)>>1), u.Now.Add(-freshness), 0, false)
}

func (u *Events) MatchInDuration(pattern Regexp, rect geom.Rectangle, duration time.Duration) bool {
	return u.match(pattern, rect, duration, time.Time{}, 0, false)
}

func (u *Events) MatchFreshness(pattern Regexp, freshness time.Duration) bool {
	return u.match(pattern, geom.Rect(-99999, -99999, 99999, 99999), time.Duration(^uint64(0)>>1), u.Now.Add(-freshness), 0, false)
}

func (u *Events) MatchDeadline(pattern Regexp, deadline time.Time) bool {
	return u.match(pattern, geom.Rect(-99999, -99999, 99999, 99999), time.Duration(^uint64(0)>>1), deadline, 0, false)
}

func (u *Events) MatchInDeadline(pattern Regexp, rect geom.Rectangle, deadline time.Time) bool {
	return u.match(pattern, rect, time.Duration(^uint64(0)>>1), deadline, 0, false)
}

func (u *Events) match(p Regexp, rect geom.Rectangle, dur time.Duration, deadline time.Time, z int, alwaysin bool) bool {
	pattern := string(p)
	u.MatchCount++
	r, ok := u.regexps[pattern]
	if !ok {
		var err error
		r, err = rcompile(pattern)
		if err != nil {
			log.Println("error parsing pattern “" + pattern + "”")
			panic(err)
		}
		u.regexps[pattern] = r
	}

	ok, last := rinterp(r, u.Trace, rect, dur, deadline, z, alwaysin)

	if ok {
		u.Last = last
		u.Last.Freshness = u.Now.Sub(u.Last.StartedAt)
	}
	return ok
}

func NewEventTracer(w *glfw.Window, replay io.Reader) *Events {
	var u Events
	u.heldcur = 0
	u.regexps = map[string][]rinst{}
	if replay == nil {
		setupcallbacks(&u, w)
	} else {
		// Replay mode disables vsync.
		// This is the simplest way to synchronize recorded events and state.
		glfw.SwapInterval(0)
		u.rec = 2
		u.rec0 = time.Now()
		u.playc = make(chan time.Time)
		go u.replay(replay)
	}

	u.Trace = u.tr[tracelen/2 : tracelen-1]

	return &u
}

func (u *Events) replay(r io.Reader) {
	err := gob.NewDecoder(r).Decode(&(u.records))
	if err != nil && err != io.EOF {
		panic(err)
	}

	re := u.records
	// Last event is F5 press, skip it.
	d := 0 * time.Second
	for i, r := range re[:len(re)-1] {
		d = re[i+1].T.Sub(r.T)
		u.future = re[i+1]
		u.emit(r.E, r.Pt, r.T.Add(1*time.Millisecond))
		u.playc <- r.T
		time.Sleep(d)
	}
	u.rec = 3
	os.Exit(0)
	// t := last(re).T
	// for {
	// 	m := 100 * time.Millisecond
	// 	t = t.Add(m)
	// 	u.playc <- t
	// 	time.Sleep(m)
	// }
}
