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
	historyIndex   int
)

func (c *Client) LaunchUI() {
	app = tview.NewApplication().
		EnableMouse(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyESC:
				mainWindow.ScrollToEnd()
			case tcell.KeyPgUp:
				app.SetFocus(mainWindow)
			case tcell.KeyPgDn:
				app.SetFocus(mainWindow)
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
		SetMaxLines(16).
		SetChangedFunc(func() {
			app.Draw()
		})
	input = tview.NewInputField().
		SetDoneFunc(c.handleInput)
	grid = tview.NewGrid().
		SetColumns(0, 40).
		SetRows(16, 0, 1).
		SetBorders(true).
		AddItem(chatWindow, 0, 0, 1, 1, 16, 0, false).
		AddItem(overheadWindow, 0, 1, 1, 1, 16, 40, false).
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
		historyIndex = 0
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
