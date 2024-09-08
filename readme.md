# Contraption
A real-time programmable vector graphics editor and GUI framework.

## Intro
```
package main

import (
	"github.com/neputevshina/contraption"
)

// Create your world (a scene), where all GUI state is placed.
// It is not required, but will be really handy when you define your own components.
func World struct {
	*contraption.World
}

func main() {
	wo := World{World: contraption.New(contraption.Config{})}
	
	for wo.Next() {
		wo.Root(
			wo.Compound(
				wo.Halign(0.5),
				wo.Valign(0.5),
				// hex() is a one of recommended util functions. See below.
				wo.SampleText(`Sample text`).ColorFill(hex(`#ff00ff`))))
				
		wo.Develop()
	}
}
```
## Sorm (shape or modifier)
- Describe it
- **We use function call mechanics for topological sort of components to minimize copying. This means that sometimes you should use closures to place component correctly in a tree. If you not do this you will get an incorrect layout.**
- Actual list of shapes and modifiers is on Godoc.

## State
Contraption uses an unusual approach to state in components. Instead of storing an internal DOM and diffing it or updating it by signals, it stores only latest events history, activator (focus object) stack and a tree from previous update cycle. All other state is external and managed by end user. Components can use regular expressions on events trace. It is enough for 80% of internal UI state, such like double-clicks, hover effects of *static components* and text field state. 

However, it is a compromise. For example, you can't do button release effect on a scrollable pane (at least yet, there is no scrollable pane also), because object relies on it's previous position to do it. I see using this system instead of VDOM is like using screen-space techniques instead of path tracing.
## Events trace and regular expressions
- Regexps are matched left to right →
- **In a trace and in regexp syntax, last event is at left ←**. `(*World).Events.Trace[0]` is also the latest event.
- Modifiers: `:in :out :before :after`
- `!` negates a symbol — match anything except this.
- `*` and `?` work like intended, `+` is not there for reasons I don't remember.
# Recipes
### Hotkey
From Contraption itself:
```
if wo.Events.Match(`!Release(Ctrl)* Press(Ctrl)`) {
	if wo.Events.Match(`!Release(Shift)* Press(Shift)`) {
		if wo.Events.Match(`Press(I)`) {
			wo.Goggles.on = !wo.Goggles.on
		}
	}
}
```
This will work too, but is order-dependent:
```
if wo.Events.Match(`Press(I) !Release(Ctrl)* !Release(Shift)* Press(Shift) !Release(Ctrl)* Press(Ctrl)`) {
	wo.Goggles.on = !wo.Goggles.on
}
```
# TBD
- Documentation, examples
- Scrolling
- Rotations
- Multiline text, sensible text interface
- Refine layout model so Transform makes sense
- Refine regexp interface (key word is `:anywhere`)

See also top comment in `contraption.go` 

# Acknowledgements

Данный проект разрабатывался как часть реализации гранта УМНИК [Фонда Содействия Инновациям](https://fasie.ru/programs/programma-umnik/). Номер договора 18384ГУ/2023.
