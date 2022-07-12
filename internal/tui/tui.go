package tui

import (
	"io"
	"log"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const SLEEP_INTERVAL = time.Millisecond * 10

type window struct {
	*tview.TextView
	writer      io.Writer
	content     []string
	bufferIndex int
	bufferSize  int
}

type TUI struct {
	app              *tview.Application
	grid             *tview.Grid
	windows          map[string]*window
	input            *tview.InputField
	inputHistory     []string
	inputHighlighted bool
	historyIndex     int
	dataReady        bool
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

var mainWindow = Window{
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

func NewTUI() *TUI {
	tui := &TUI{}
	tui.windows = make(map[string]*window)
	tui.grid = tview.NewGrid()
	tui.input = tview.NewInputField().
		SetDoneFunc(tui.handleInput).
		SetFieldBackgroundColor(tcell.ColorBlack)
	tui.inputHistory = []string{""}

	tui.app = tview.NewApplication().
		EnableMouse(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				scrollToEnd(tui.windows["main"])
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
				}
				tui.app.SetFocus(tui.input)
			}
			return event
		})
	// Default view just main window and input bar
	tui.AddWindow("main", mainWindow)
	tui.grid.AddItem(tui.input, 2, 0, 1, 1, 0, 0, true).
		SetRows(0, 1)

	return tui
}

func setScrollBuffer(w *window, scrollTarget int) {
	// Make a new string to fill the view buffer
	ns := w.content[w.bufferIndex : w.bufferIndex+w.bufferSize]

	// Write the string to the window
	w.Clear()
	w.writer.Write([]byte(strings.Join(ns, "\n")))
	w.ScrollTo(scrollTarget, 0) // re-scroll to the provided offset
}

func scrollToEnd(w *window) {
	if w.bufferIndex > 0 {
		w.bufferIndex = 0
		setScrollBuffer(w, 0)
	} else {
		w.ScrollToEnd()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func scroll(w *window, action tview.MouseAction, event *tcell.EventMouse) {
	row, _ := w.GetScrollOffset()
	switch action {
	case tview.MouseScrollUp:
		// We're at the top of the view buffer
		if row <= 0 {
			_, _, _, height := w.GetInnerRect()
			// Increment the bufferIndex by half the buffer size
			// If that would carry us past the end of the buffer,
			// then set the index to the last chunk of the buffer
			newIndex := max(
				(w.bufferIndex + w.bufferSize),  // Incremented buffer value
				(len(w.content) - w.bufferSize), // Last chunk of the buffer
			)

			// What line in the new buffer are we going to scroll to
			//scrollTarget := (newIndex - w.bufferIndex) + height
			w.bufferIndex = newIndex - height

			setScrollBuffer(w, 0)
		}
	case tview.MouseScrollDown:
		if w.bufferIndex > 0 {
			// We're scrolling

		}
		_, _, _, height := w.GetInnerRect()
		rowIndex := row - height
		if rowIndex >= 0 {
			if w.bufferIndex > 0 {
				// We're scrolling
				newIndex := min(
					0,                              // bottom of the buffer
					w.bufferIndex-(w.bufferSize/2), // decrement by half the buffer size
				)
				scrollTarget := w.bufferIndex - newIndex
				w.bufferIndex = newIndex
				setScrollBuffer(w, scrollTarget)
			}
		}
	}
}

func (t *TUI) AddWindow(name string, win Window) {
	var w *window
	if cw, ok := t.windows[name]; ok {
		w = cw
		t.grid.RemoveItem(w)
	} else {
		nw := tview.NewTextView()
		nw.SetBorder(win.Border)
		nw.SetScrollable(win.Scrollable)
		nw.SetDynamicColors(true)
		nw.SetMaxLines(win.MaxLines)
		nw.SetChangedFunc(func() {
			t.app.Draw()
		})
		wr := tview.ANSIWriter(nw)
		w = &window{nw, wr, make([]string, 0), 0, 500}
		if win.Scrollable {
			w.bufferSize = win.MaxLines
			nw.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
				scroll(w, action, event)
				return action, event
			})
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
		go t.Parse(text)
		t.inputHighlighted = true
		if t.inputHistory[len(t.inputHistory)-1] != text {
			t.inputHistory = append(t.inputHistory, text)
		}
	}
}

func (t *TUI) Show(name string, text string) {
	t.dataReady = true
	if w, ok := t.windows[name]; ok {
		w.content = append(w.content, strings.Split(text, "\n")...)
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
	if err := t.app.SetRoot(t.grid, true).SetFocus(t.input).Run(); err != nil {
		log.Fatal(err)
	}
}
