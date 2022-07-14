package tui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
)

type TUI struct {
	views.BoxLayout

	app     *views.Application
	windows map[string]*Window
	parse   func(string)
	input   *views.Text
}

type Window struct {
	*views.SimpleStyledText

	x, y, width, height int
	content             string
}

func NewTUI(parse func(string)) *TUI {
	t := &TUI{}
	t.SetOrientation(views.Vertical)

	t.app = &views.Application{}
	t.app.SetStyle(tcell.StyleDefault)

	s, err := tcell.NewTerminfoScreen()
	if err != nil {
		log.Fatal(err)
	}
	t.app.SetScreen(s)

	t.windows = map[string]*Window{}
	t.AddWindow("main")
	t.Print("main", "main window")

	t.input = views.NewText()
	t.input.SetText("> ")
	t.input.SetAlignment(views.VAlignBottom | views.HAlignLeft)

	t.AddWidget(t.windows["main"], 1)
	t.AddWidget(t.input, 0)

	t.app.SetRootWidget(t)

	t.parse = parse

	return t
}
func (t *TUI) Print(name string, text string) {
	if w, ok := t.windows[name]; ok {
		//t.app.PostFunc(func() {
		w.content += text
		w.SetText(w.content)
		//t.app.Refresh()
		//})
	}
}
func (t *TUI) AddWindow(name string) *Window {

	w := &Window{
		content:          "",
		SimpleStyledText: views.NewSimpleStyledText(),
	}
	t.windows[name] = w
	return w
}
func (t *TUI) RemoveWindow(name string) {
	delete(t.windows, name)
}
func (t *TUI) Run() {
	if err := t.app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

var highlighted = false

func (t *TUI) HandleEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyCtrlL:
			t.app.Refresh()
			return true
		case tcell.KeyCtrlC:
			t.app.Quit()
			return true
		case tcell.KeyDEL, tcell.KeyBackspace:
			if highlighted {
				highlighted = false
				t.input.SetText("> ")
				t.input.SetStyle(tcell.StyleDefault)
			} else {
				if t.input.Text() != "> " {
					t.input.SetText(t.input.Text()[:len(t.input.Text())-1])
				}
			}
		case tcell.KeyEnter:
			t.parse(t.input.Text()[2:])
			highlighted = true
			//t.Print("main", t.input.Text()[2:])
			t.input.SetStyle(tcell.StyleDefault.Background(tcell.ColorDimGray))
		default:
			if ev.Key() == tcell.KeyRune {
				if highlighted {
					highlighted = false
					t.input.SetText("> ")
					//t.input.SetStyle(tcell.StyleDefault)
				}
				t.input.SetText(t.input.Text() + string(ev.Rune()))
			} else {
				highlighted = false
			}
		}
	}
	return t.BoxLayout.HandleEvent(ev)
}
