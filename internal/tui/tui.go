package tui

import (
	"io"
	"log"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/seandheath/gomc/pkg/plugin"
)

const SLEEP_INTERVAL = time.Millisecond * 10

type window struct {
	*cview.TextView
	writer    io.Writer
	content   string
	scrolling bool
}

type TUI struct {
	app              *cview.Application
	grid             *cview.Grid
	windows          map[string]*window
	input            *cview.InputField
	inputHistory     []string
	inputHighlighted bool
	historyIndex     int
	dataReady        bool
	parse            func(string)
}

var mainWindow = plugin.Window{
	Row:           1,
	Col:           0,
	RowSpan:       1,
	ColSpan:       1,
	MinGridHeight: 0,
	MinGridWidth:  0,
	Border:        false,
	Scrollable:    true,
	MaxLines:      200,
}

func NewTUI(parse func(string)) *TUI {
	tui := &TUI{}
	tui.windows = make(map[string]*window)
	tui.grid = cview.NewGrid()
	tui.input = cview.NewInputField()
	tui.input.SetDoneFunc(tui.handleInput)
	tui.input.SetFieldBackgroundColor(tcell.ColorBlack)
	tui.inputHistory = []string{""}

	tui.parse = parse

	tui.app = cview.NewApplication()
	tui.app.EnableMouse(true)
	tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			tui.windows["main"].ScrollToEnd()
			//tui.scrollToEnd(tui.windows["main"])
		case tcell.KeyPgUp:
			tui.app.SetFocus(tui.windows["main"])
			//tui.scrollUp(tui.windows["main"])
		case tcell.KeyPgDn:
			tui.app.SetFocus(tui.windows["main"])
			//tui.scrollDown(tui.windows["main"])
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
			}
			tui.app.SetFocus(tui.input)
		}
		return event
	})
	// Default view just main window and input bar
	tui.AddWindow("main", mainWindow)
	tui.grid.AddItem(tui.input, 2, 0, 1, 1, 0, 0, true)
	tui.grid.SetRows(0, 1)

	return tui
}

func (t *TUI) AddWindow(name string, win plugin.Window) {
	var w *window
	if cw, ok := t.windows[name]; ok {
		w = cw
		t.grid.RemoveItem(w)
	} else {
		nw := cview.NewTextView()
		nw.SetBorder(win.Border)
		nw.SetScrollable(win.Scrollable)
		nw.SetDynamicColors(true)
		nw.SetMaxLines(win.MaxLines)
		nw.SetWordWrap(false)
		nw.SetRegions(false)
		nw.SetChangedFunc(func() {
			t.app.Draw()
		})
		wr := cview.ANSIWriter(nw)
		w = &window{
			TextView: nw,
			writer:   wr,
			content:  "",
		}
		t.windows[name] = w
	}
	t.grid.AddItem(w,
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
		go t.parse(text)
		t.inputHighlighted = true
		if t.inputHistory[len(t.inputHistory)-1] != text {
			t.inputHistory = append(t.inputHistory, text)
		}
	}
}

func (t *TUI) Print(name string, text string) {
	t.dataReady = true
	if w, ok := t.windows[name]; ok {
		w.content += text
		t.windows[name].writer.Write([]byte(text))
	}
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
	t.app.SetRoot(t.grid, true)
	t.app.SetFocus(t.input)
	if err := t.app.Run(); err != nil {
		log.Fatal(err)
	}
}
