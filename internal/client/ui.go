package client

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app            *tview.Application
	grid           *tview.Grid
	mainWindow     *tview.TextView
	chatWindow     *tview.TextView
	overheadWindow *tview.TextView
	input          *tview.InputField
	inputHistory   []string
)

func (c *Client) LaunchUI() {
	inputHistory := make([]string, 1)
	inputHistory[0] = ""

	//ScrollBuffer := make([]string, 0)

	app = tview.NewApplication().
		EnableMouse(true)
	mainWindow = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	chatWindow = tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	overheadWindow = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(false).
		SetMaxLines(18).
		SetChangedFunc(func() {
			app.Draw()
		})
	input = tview.NewInputField().
		SetDoneFunc(c.handleInput)
	grid = tview.NewGrid().
		SetColumns(0, 40).
		SetRows(18, 0, 1).
		SetBorders(true).
		AddItem(chatWindow, 0, 0, 1, 1, 18, 0, false).
		AddItem(overheadWindow, 0, 1, 1, 1, 18, 40, false).
		AddItem(mainWindow, 1, 0, 1, 2, 0, 0, false).
		AddItem(input, 2, 0, 1, 2, 1, 0, true)

	if err := app.SetRoot(grid, true).SetFocus(input).Run(); err != nil {
		log.Fatal(err)
	}
}

func (c *Client) handleInput(key tcell.Key) {
	text := input.GetText()
	switch key {
	case tcell.KeyEnter:
		if text == "" {
			// Redo the last command
			text = inputHistory[len(inputHistory)-1]
		}
		c.Parse(text)
		input.SetText("")
		inputHistory = append(inputHistory, text)
	}
}

func (c *Client) ShowMain(text string)     { c.Show(text, mainWindow) }
func (c *Client) ShowChat(text string)     { c.Show(text, chatWindow) }
func (c *Client) ShowOverhead(text string) { c.Show(text, overheadWindow) }
func (c *Client) Show(text string, w *tview.TextView) {
	w.Write([]byte(text))
}
