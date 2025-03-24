Крутилка
```
wo.Solve(
	func(mark int) contraption.SormId {
		if mark == 1 {
			return wo.Compound(
				wo.Rotate(PI/4),
				wo.Rectangle(10,10),
				wo.Void(0,0),
				wo.Halign(0.5),
				wo.Valign(0.5))
		}
	},
	wo.Compound(
		wo.Rotate(mix(-3*math.PI/2, 0, pos)),
		wo.Rotate(-PI/4),
		wo.Compound(
			wo.Translate(0, knobRadius),
			wo.Mark(1)),
	))
```

Параллельная диагонали палка
```
wo.Solve(func(so *contraption.Solver) {
	return so.Solve(
		so.Square(1001),
		so.Expand(1001, 1, 2, 3, 4),
		so.Line(1002, 2, 4),
		so.Line(1003, 0, 0),
		so.Parallel(1003, 1002),
		so.Point(1004),
		so.Eq(1004, x, y),
		so.Through(1004, 1003),
		so.Clip(1003, 1001))
})
```
```
wo.Compound(
	wo.Solve(func(so *contraption.Solver) {
		sq := so.Square()
		diag := so.Line(sq.TR(), sq.BL())
		t := so.Empty.Line().
			Parallel(diag).
			Through(geom.Pt(x, y)).
			And(sq)

		so.Setmark(1, t.A)
		so.Setmark(2, t.B)
	}),
	
)
```
Кнопки тонского
```
t1 := so.Text("Cancel")
t2 := so.Text("Apply")
t3 := so.Text("OK")

m := func(Sorm c) { 
	so.Component(
		func(SormFunc in) Sorm { return wo.Button(-1, -1, in) }, 
		c)
}

```
