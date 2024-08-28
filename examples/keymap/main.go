package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"

	"github.com/neputevshina/contraption"

	"github.com/neputevshina/geom"
	"github.com/neputevshina/nanovgo"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	AhkScriptPath    = `c:\Users\User\Desktop\celestial.ahk`
	AhkPath          = `c:\Program Files\AutoHotkey\v2\AutoHotkey64.exe`
	SoundProgramPath = `C:\Users\User\Documents\asdf\adf.exe`
)

var (
	buttons = map[rune]string{}
	ahkcmd  *exec.Cmd
)

type Sorm = contraption.Sorm
type World struct {
	contraption.World
	Text func(size float64, str []rune) Sorm
}

func (wo *World) Button(b rune) Sorm {
	black := hexpaint(`#000000`)
	return wo.Compound(
		wo.Halign(0.5),
		wo.Valign(0.5),
		wo.Roundrect(50, 50, 6).Stroke(black).Strokewidth(3).CondFill(func(rect geom.Rectangle) nanovgo.Paint {
			if wo.MatchIn(`Drop(audio/mpeg):in`, rect) {
				buttons[b] = wo.Events.Trace[0].E.(contraption.Drop).Paths[0]
				// saveahk(buttons, AhkScriptPath)
				// reloadahk()
			}
			if wo.MatchIn(`Click(3):in`, rect) {
				delete(buttons, b)
				// saveahk(buttons, AhkScriptPath)
				// reloadahk()
			}

			v, ok := buttons[b]
			if ok {
				return hexpaint(`#00ff00`)
			}
			_ = v
			return hexpaint(`#00000000`)
		}),
		wo.Compound(
			wo.Text(20, []rune{b}).Fill(black)))
}

func (wo *World) BetweenVoid(w, h float64) Sorm {
	return wo.Between(func() contraption.Sorm { return wo.Void(w, h) })
}

func main() {
	// buttons = loadahk(AhkScriptPath)
	buttons = map[rune]string{}
	reloadahk()

	wo := World{
		World: contraption.New(contraption.Config{}),
	}
	wo.Text = wo.NewVectorText(goregular.TTF)

	for wo.Next() {
		wo.Root(
			wo.Compound(
				wo.Halign(0.5),
				wo.Valign(0.5),
				wo.Void(wo.Wwin, wo.Hwin),
				wo.Compound(
					wo.Vfollow(),
					wo.BetweenVoid(0, 10),
					wo.Compound(
						wo.Hfollow(),
						wo.BetweenVoid(10, 0),
						wo.Button('1'),
						wo.Button('2'),
						wo.Button('3'),
						wo.Button('4'),
						wo.Button('5'),
						wo.Button('6'),
						wo.Button('7'),
						wo.Button('8'),
						wo.Button('9'),
						wo.Button('0'),
						wo.Button('-'),
						wo.Button('='),
					),
					wo.Compound(
						wo.Hfollow(),
						wo.BetweenVoid(10, 0),
						wo.Void(25, 0).Betweener(),
						wo.Button('Q'),
						wo.Button('W'),
						wo.Button('E'),
						wo.Button('R'),
						wo.Button('T'),
						wo.Button('Y'),
						wo.Button('U'),
						wo.Button('I'),
						wo.Button('O'),
						wo.Button('P'),
						wo.Button('['),
						wo.Button(']'),
					),
					wo.Compound(
						wo.Hfollow(),
						wo.BetweenVoid(10, 0),
						wo.Void(25+15, 0).Betweener(),
						wo.Button('A'),
						wo.Button('S'),
						wo.Button('D'),
						wo.Button('F'),
						wo.Button('G'),
						wo.Button('H'),
						wo.Button('J'),
						wo.Button('K'),
						wo.Button('L'),
						wo.Button(';'),
						wo.Button('\''),
					),
					wo.Compound(
						wo.Hfollow(),
						wo.BetweenVoid(10, 0),
						wo.Void(25+15+15, 0).Betweener(),
						wo.Button('Z'),
						wo.Button('X'),
						wo.Button('C'),
						wo.Button('V'),
						wo.Button('B'),
						wo.Button('N'),
						wo.Button('M'),
						wo.Button(','),
						wo.Button('.'),
						wo.Button('/')))))

		wo.Develop()
	}

	ahkcmd.Process.Kill()
}

func saveahk(buttons map[rune]string, savepath string) {
	file, err := os.Create(savepath)
	if err != nil {
		panic(err)
	}
	for k, v := range buttons {
		fmt.Fprintf(file, `~%s::Run "%s -f %s"%s`, string(k), SoundProgramPath, v, "\n")
	}
	file.Close()
}

var lineregexp = regexp.MustCompile(`~(.)::Run ".* -f (.*)"`)

func loadahk(savepath string) (buttons map[rune]string) {
	buttons = map[rune]string{}
	file, err := os.Open(savepath)
	if err != nil {
		panic(err)
	}
	bfr := bufio.NewReader(file)
	for {
		bs, err := bfr.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		matches := lineregexp.FindSubmatch(bs)
		buttons[rune(matches[1][0])] = string(matches[2])
	}
}

func reloadahk() {
	if ahkcmd != nil {
		err := ahkcmd.Process.Kill()
		if err != nil {
			println(err)
		}
	}
	ahkcmd = exec.Command(AhkPath, "/force", AhkScriptPath)
	ahkcmd.Stdout = os.Stdout
	err := ahkcmd.Start()
	if err != nil {
		println(err)
	}
}

var println = fmt.Println
