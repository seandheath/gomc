package tui

import (
	"io"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type window struct {
	*tview.TextView
	writer io.Writer
}

type TUI struct {
	app              *tview.Application
	grid             *tview.Grid
	windows          map[string]*window
	input            *tview.InputField
	inputHistory     []string
	inputHighlighted bool
	historyIndex     int
	Parse            func(string)
}

type Window struct {
	Row           int  `yaml:"row"`
	Col           int  `yaml:"col"`
	RowSpan       int  `yaml:"rowspan"`
	ColSpan       int  `yaml:"colspan"`
	MinGridHeight int  `yaml:"mingridheight"`
	MinGridWidth  int  `yaml:"mingridwidth"`
	Border        bool `yaml:"border"`
	Scrollable    bool `yaml:"scrollable"`
	MaxLines      int  `yaml:"maxlines"`
}

var defaultMainWindow = Window{
	Row:           0,
	Col:           0,
	RowSpan:       1,
	ColSpan:       1,
	MinGridHeight: 0,
	MinGridWidth:  0,
	Border:        false,
	Scrollable:    true,
	MaxLines:      100000,
}

func NewTUI() *TUI {
	tui := &TUI{}
	tui.windows = make(map[string]*window)
	tui.grid = tview.NewGrid()
	tui.input = tview.NewInputField().
		SetDoneFunc(tui.handleInput).
		SetFieldBackgroundColor(tcell.ColorBlack)

	tui.app = tview.NewApplication().
		EnableMouse(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				tui.windows["main"].ScrollToEnd()
			case tcell.KeyPgUp:
				tui.app.SetFocus(tui.windows["main"])
			case tcell.KeyPgDn:
				tui.app.SetFocus(tui.windows["main"])
			case tcell.KeyUp:
				if len(tui.inputHistory) > 0 {
					tui.historyIndex += 1
					if tui.historyIndex > len(tui.inputHistory) {
						tui.historyIndex = len(tui.inputHistory)
					}
					tui.input.SetText(tui.inputHistory[len(tui.inputHistory)-tui.historyIndex])
				}
			case tcell.KeyDown:
				tui.historyIndex -= 1
				if tui.historyIndex <= 0 {
					tui.historyIndex = 0
					tui.input.SetText("")
				} else {
					tui.input.SetText(tui.inputHistory[len(tui.inputHistory)-tui.historyIndex])
				}
			case tcell.KeyEnter:
				tui.app.SetFocus(tui.input)
			default:
				if tui.inputHighlighted {
					tui.inputHighlighted = false
					tui.input.SetText("")
					tui.input.SetBackgroundColor(tcell.ColorBlack)
				}
				tui.app.SetFocus(tui.input)
			}
			return event
		})
	// Default view just main window and input bar
	tui.AddWindow("main", defaultMainWindow)
	tui.grid.AddItem(tui.input, 1, 0, 1, 1, 0, 0, true).
		SetRows(0, 1)

	return tui
}

func (tui *TUI) AddWindow(name string, win Window) {
	var w *window
	if cw, ok := tui.windows[name]; ok {
		w = cw
		tui.grid.RemoveItem(cw)
	} else {
		nw := tview.NewTextView()
		nw.SetBorder(win.Border)
		nw.SetScrollable(win.Scrollable)
		nw.SetMaxLines(win.MaxLines)
		nw.SetDynamicColors(true)
		nw.SetChangedFunc(func() {
			tui.app.Draw()
		})
		wr := tview.ANSIWriter(nw)
		w = &window{nw, wr}
		tui.windows[name] = w
	}
	tui.grid.AddItem(w,
		win.Row,
		win.Col,
		win.RowSpan,
		win.ColSpan,
		win.MinGridHeight,
		win.MinGridWidth,
		false,
	)
}

func (t *TUI) handleInput(key tcell.Key) {
	text := t.input.GetText()
	switch key {
	case tcell.KeyEnter:
		t.historyIndex = 0
		t.Parse(text)
		t.inputHighlighted = true
		t.input.SetBackgroundColor(tcell.ColorLightBlue)
		t.inputHistory = append(t.inputHistory, text)
	}
}

func (t *TUI) Show(name string, text string) {
	t.windows[name].writer.Write([]byte(text))
}

func (t *TUI) FixInputLine(rows []int, cols []int) {
	// Reset input line
	t.grid.RemoveItem(t.input)
	t.grid.AddItem(t.input,
		len(rows), // row
		0,         // col
		1,         // rowSpan
		len(cols), // colSpan
		0,         // minGridHeight
		0,         // minGridWidth
		true,      // focus
	)
	if len(cols) > 0 {
		t.grid.SetColumns(cols...)
	}
	if len(rows) > 0 {
		// Add a row for the input line
		t.grid.SetRows(append(rows, 1)...)
	}

}

func (t *TUI) Run() {
	if err := t.app.SetRoot(t.grid, true).SetFocus(t.input).Run(); err != nil {
		log.Fatal(err)
	}
}
