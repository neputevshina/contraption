?? wo.Readlimit(n int) — Limits the number of components read from sequences in this compound. Negatve values remove the limit. -1 by default.

- Sequences are not buferized interally. contraption.BufferSequence.
- Need scissor.
- contraption.ReverseSequence
- Impl:
	- Apply it
	- While applying allocate to `auxpool`
	- When allocated rewrite indices in original Sequence Sorm
	- When drawing distract to `auxpool`

```
type Scroll struct {
	Offset geom.Point // May be negative
	I int
}

// 
	wo.Compound(
		wo.Vscroll(scroll),
		wo.Vfollow(),
		wo.Valign(1),
		wo.Halign(0.5),
		wo.BetweenVoid(0, 10),
		wo.Sequence(contraption.Slice(...)))
```

Special behaviors:
- wo.Vscroll + wo.Vfollow + wo.Valign
	- Chat: Valign(1)
	- Page: Valign(0) (Default)
	- Funky: Valign(anything else)
- 

|               | S No scrolling          | F No scrolling          | F In-axis scrolling          | F Tangential scrolling |
| ------------- | ----------------------- | ----------------------- | ---------------------------- | ---------------------- |
| No sequence   |                         |                         |                              |                        |
| With sequence | Read whole, paste whole | Read whole, paste whole | Start reading from scroll.I, |                        |
*S — Stack, F — Follow*