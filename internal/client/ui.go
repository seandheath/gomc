package client

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	App            *tview.Application
	Grid           *tview.Grid
	MainWindow     *tview.TextView
	ChatWindow     *tview.TextView
	OverheadWindow *tview.TextView
	Input          *tview.InputField
	inputHistory   []string
	lineCount      int64
)

func Launch() {
	inputHistory := make([]string, 1)
	inputHistory[0] = ""

	//ScrollBuffer := make([]string, 0)

	App = tview.NewApplication().
		EnableMouse(true)
	MainWindow = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			App.Draw()
		})
	ChatWindow = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			App.Draw()
		})
	OverheadWindow = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false).
		SetMaxLines(18).
		SetChangedFunc(func() {
			App.Draw()
		})
	Input = tview.NewInputField().
		SetDoneFunc(handleInput)
	Grid = tview.NewGrid().
		SetColumns(0, 40).
		SetRows(18, 0, 1).
		SetBorders(true).
		AddItem(ChatWindow, 0, 0, 1, 1, 18, 0, false).
		AddItem(OverheadWindow, 0, 1, 1, 1, 18, 40, false).
		AddItem(MainWindow, 1, 0, 1, 2, 0, 0, false).
		AddItem(Input, 2, 0, 1, 2, 1, 0, true)

	if err := App.SetRoot(Grid, true).SetFocus(Input).Run(); err != nil {
		log.Fatal(err)
	}
}

func handleInput(key tcell.Key) {
	text := Input.GetText()
	switch key {
	case tcell.KeyEnter:
		if text == "" {
			// Redo the last command
			text = inputHistory[len(inputHistory)-1]
		}
		Parse(text)
		Input.SetText("")
		inputHistory = append(inputHistory, text)
	}
}

func ShowMain(text string)     { Show(text, MainWindow) }
func ShowChat(text string)     { Show(text, ChatWindow) }
func ShowOverhead(text string) { Show(text, OverheadWindow) }
func Show(text string, w *tview.TextView) {
	w.Write([]byte(text))
}
