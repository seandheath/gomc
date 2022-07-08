package client

import (
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type window struct {
	*tview.TextView
	writer io.Writer
}

var (
	app          *tview.Application
	grid         *tview.Grid
	windows      map[string]*window
	input        *tview.InputField
	inputHistory []string
	historyIndex int
)

func AddWindow(name string, win Window) {
	var w *window
	if cw, ok := windows[name]; ok {
		w = cw
		grid.RemoveItem(cw)
	} else {
		nw := tview.NewTextView()
		nw.SetBorder(win.Border)
		nw.SetScrollable(win.Scrollable)
		nw.SetMaxLines(win.MaxLines)
		nw.SetDynamicColors(true)
		wr := tview.ANSIWriter(nw)
		w = &window{nw, wr}
		windows[name] = w
	}
	grid.AddItem(w,
		win.Row,
		win.Col,
		win.RowSpan,
		win.ColSpan,
		win.MinGridHeight,
		win.MinGridWidth,
		false,
	)
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

func uiInit() {
	windows = make(map[string]*window)
	grid = tview.NewGrid()

	input = tview.NewInputField().
		SetDoneFunc(handleInput).
		SetFieldBackgroundColor(tcell.ColorBlack)

		// Default view just main window and input bar
	AddWindow("main", defaultMainWindow)
	grid.AddItem(input, 1, 0, 1, 1, 0, 0, true).
		SetRows(0, 1)

	app = tview.NewApplication().
		EnableMouse(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				windows["main"].ScrollToEnd()
			case tcell.KeyPgUp:
				app.SetFocus(windows["main"])
			case tcell.KeyPgDn:
				app.SetFocus(windows["main"])
			case tcell.KeyUp:
				if len(inputHistory) > 0 {
					historyIndex += 1
					if historyIndex > len(inputHistory) {
						historyIndex = len(inputHistory)
					}
					input.SetText(inputHistory[len(inputHistory)-historyIndex])
				}
			case tcell.KeyDown:
				historyIndex -= 1
				if historyIndex <= 0 {
					historyIndex = 0
					input.SetText("")
				} else {
					input.SetText(inputHistory[len(inputHistory)-historyIndex])
				}
			default:
				app.SetFocus(input)
			}
			return event
		})
}

func handleInput(key tcell.Key) {
	text := input.GetText()
	switch key {
	case tcell.KeyEnter:
		historyIndex = 0
		if text == "" {
			// Redo the last command
			text = inputHistory[len(inputHistory)-1]
		}
		Parse(text)
		input.SetText("")
		inputHistory = append(inputHistory, text)
	}
}

func Show(name string, text string) {
	_, err := windows[name].writer.Write([]byte(text))
	if err != nil {
		LogError.Println(err.Error())
	}
}
