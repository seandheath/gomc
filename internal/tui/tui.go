package tui

import (
	"log"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/gdamore/tcell/v2"
	"github.com/seandheath/gomc/pkg/plugin"
)

const SLEEP_INTERVAL = time.Millisecond * 10

type window struct {
	*cview.TextView
	content    []byte
	scrolling  bool
	scrollable bool
}

type TUI struct {
	app              *cview.Application
	grid             *cview.Grid
	windows          map[string]*window
	input            *cview.InputField
	inputHistory     []string
	inputHighlighted bool
	historyIndex     int
	parse            func(string)
}

var mainWindow = plugin.Window{
	Row:           0,
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
	//tui.app.SetAfterResizeFunc(func(w int, h int) {
	//for _, w := range tui.windows {
	//if w.scrollable {
	//resizeWindow(w)
	//}
	//}
	//resizeWindow(tui.windows["main"])
	//})
	tui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			tui.app.SetFocus(tui.windows["main"])
			tui.scroll("end", tui.windows["main"])
		case tcell.KeyPgUp:
			tui.app.SetFocus(tui.windows["main"])
			tui.scroll("up", tui.windows["main"])
		case tcell.KeyPgDn:
			tui.app.SetFocus(tui.windows["main"])
			tui.scroll("down", tui.windows["main"])
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
		case tcell.KeyDEL, tcell.KeyBackspace:
			if tui.inputHighlighted {
				tui.inputHighlighted = false
				tui.input.SetText("")
			}
			tui.app.SetFocus(tui.input)
		case tcell.KeyRune:
			if tui.inputHighlighted {
				tui.inputHighlighted = false
				tui.input.SetText("")
			}
			if !tui.input.HasFocus() {
				tui.input.SetText(tui.input.GetText() + string(event.Rune()))
			}
			tui.app.SetFocus(tui.input)
		}
		return event
	})
	// Default view just main window and input bar
	tui.AddWindow("main", mainWindow)
	tui.grid.AddItem(tui.input, 1, 0, 1, 1, 0, 0, true)
	tui.grid.SetRows(0, 1)

	return tui
}

func resizeWindow(w *window) {
	if !w.scrolling {
		_, _, _, h := w.GetInnerRect()
		w.SetMaxLines(h)
	}
}

func (t *TUI) scroll(direction string, w *window) {
	if w.scrollable {
		switch direction {
		case "up":
			// If we're at the bottom and haven't started scrolling yet then we
			// need to switch to the scroll buffer mode
			_, _, _, h := w.GetInnerRect()
			trow, _ := w.GetBufferSize()
			if !w.scrolling && h <= trow {
				// We're just starting to scroll up, so we clear the buffer, set
				// the max lines to infinite, and write the contents of the
				// scrollbuffer (w.content) to the window. This will stop new
				// writes from being written to the window until we scroll back
				// to the bottom
				w.scrolling = true
				w.SetMaxLines(0)
				w.SetBytes(w.content)
				//w.ScrollTo(trow, 0)
				//t.app.Draw(w)
				//w.ScrollToEnd() // This doesn't seem to be working...
			}
		case "down":
			// Check if we're at the bottom
			_, _, _, h := w.GetInnerRect() // height of the screen rect
			crow, _ := w.GetScrollOffset() // if we're at the top this is 0
			trows, _ := w.GetBufferSize()  // number of lines in the screen buff
			// total rows less the height of the rectangle will point at
			// the current row
			if (trows - h) == crow {
				t.scroll("end", w)
			}
		case "end":
			// We've got a bigger buffer than the height of the screen, so we
			// need to scroll to the bottom
			if w.scrolling {
				_, _, _, h := w.GetInnerRect() // height of the screen rect
				w.SetMaxLines(h)
				w.SetBytes(w.content)
			} else {
				resizeWindow(w)
			}
			w.scrolling = false
		}
	}
}

func (t *TUI) AddWindow(name string, win plugin.Window) {
	var w *window
	if cw, ok := t.windows[name]; ok {
		w = cw
		t.grid.RemoveItem(w)
	}
	nw := cview.NewTextView()
	nw.SetBorder(win.Border)
	nw.SetScrollable(win.Scrollable)
	nw.SetDynamicColors(true)
	nw.SetMaxLines(win.MaxLines)
	nw.SetWrap(false)
	nw.SetRegions(false)
	nw.SetTextAlign(cview.AlignLeft)
	nw.SetVerticalAlign(cview.AlignBottom)
	nw.SetChangedFunc(func() {
		t.app.Draw()
	})
	//wr := cview.ANSIWriter(nw)
	w = &window{
		TextView:   nw,
		content:    make([]byte, 0),
		scrollable: win.Scrollable,
	}
	t.windows[name] = w
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

func (t *TUI) SetGrid(rows []int, cols []int) {
	t.grid.SetRows(rows...)
	t.grid.SetColumns(cols...)
	t.FixInputLine(rows, cols)
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
	t.PrintBytes(name, []byte(text))
}

func (t *TUI) PrintBytes(name string, text []byte) {
	if w, ok := t.windows[name]; ok {
		if w.scrollable {
			w.content = append(w.content, text...)
			if !w.scrolling {
				// Don't write text if we're in scroll mode
				// TODO figure out way to append new lines
				t.windows[name].Write(text)
			}
		} else {
			// not a scrollable window, write text
			t.windows[name].Write(text)
		}
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
