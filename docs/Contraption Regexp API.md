```
wo.Match -> wo.Trigger
wo.MatchIndef -> wo.Match

wo.MatchInIndef -> wo.MatchIn
wo.MatchIn -> wo.TriggerIn

.:out !Unclick(1)* Click(1):in -> .:out !Unclick(1)*:anywhere Click(1)

wo.MatchIn(`Click(1):in`, rect) -> wo.TriggerIn(`Click(1)`, rect)

Не, херня.
Внутри *In работают по умолчанию как :in, внутри обычных работают как :anywhere. При помощи соответствующих тегов можно менять.

Press(Z) .* Click(1):begin !Unclick(1)* Click(1):end 
-> 
Press(Z) .* [ Click(1) !Unclick(1)* Click(1) ]

wo.TriggerDeadlineIn(pat, )

wo.Matcher
.WithZ(int) ...
.In(rect) ...
.Anywhere() ...
.Match(pat) bool
.Trigger(pat) bool

```