package client

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Window represents a new window in the client. The window must
// provide the height, width, and X, Y coordinates of the top left corner.
// A value of 0 for X and Y indicates the top left corner.
// A value of 0 on Width or Height represents the full width or height of the terminal.

type Window struct {
	Content string
	Vp      *viewport.Model
}
type model struct {
	input        textinput.Model
	win          map[string]*Window
	inputHistory []string
	inputIndex   int
	ViewFunc     func(map[string]*Window) string
	ResizeFunc   func(int, int, map[string]*Window) map[string]*Window
}

// AddWindow adds a window as specified by the configuration file
// You must provide a new View() function for Bubbletea
func (m *model) AddWindow(name string, width int, height int) error {
	v := viewport.New(width, height)
	m.win[name] = &Window{
		Content: "",
		Vp:      &v,
	}
	return nil
}

func initialModel() *model {
	m := &model{}

	ti := textinput.New()
	ti.CursorStyle.Blink(false)
	ti.Focus()

	m.input = ti
	m.ViewFunc = DefaultView
	m.ResizeFunc = DefaultResize

	m.win = map[string]*Window{}

	return m
}

func newKeyMap() viewport.KeyMap {
	return viewport.KeyMap{
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "½ page down"),
		),
		Up: key.NewBinding(
			key.WithKeys("ctrl+up"),
			key.WithHelp("ctrl+↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("ctrl+down", "j"),
			key.WithHelp("ctrl+↓", "down"),
		),
	}
}

func (m *model) Init() tea.Cmd { return nil }

//func UpdateWindow(w *viewport.Model, msg tea.Msg) tea.Cmd {
//model, cmd := w.Update(msg)
//w = &model
//return cmd
//}
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			val := m.input.Value()
			if val == "" && len(m.inputHistory) > 0 {
				val = m.inputHistory[len(m.inputHistory)-1]
			} else {
				m.inputHistory = append(m.inputHistory, val)
			}
			m.inputIndex = len(m.inputHistory)
			go Parse(val)
			m.input.SetValue("")
		case tea.KeyUp:
			if m.inputIndex > 0 {
				m.inputIndex -= 1
				m.input.SetValue(m.inputHistory[m.inputIndex])
			}
		case tea.KeyDown:
			if m.inputIndex < len(m.inputHistory)-1 {
				m.inputIndex += 1
				m.input.SetValue(m.inputHistory[m.inputIndex])
			} else {
				m.input.SetValue("")
				m.inputIndex = len(m.inputHistory)
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			m.win["main"].Vp.GotoBottom()
		}
	case tea.WindowSizeMsg:
		inputHeight := lipgloss.Height(m.input.View())
		m.win = m.ResizeFunc(msg.Width, msg.Height-inputHeight, m.win)
	case showText:
		if w, ok := m.win[msg.window]; ok {
			ab := w.Vp.AtBottom()
			w.Content += msg.text
			w.Vp.SetContent(w.Content)
			if ab {
				w.Vp.GotoBottom()
			}
		} else {
			go ShowMain(fmt.Sprintf("\nUnable to show text [%s] in window [%s], window not found.\n", msg.text, msg.window))
		}
	}

	v, cmd := m.win["main"].Vp.Update(msg)
	m.win["main"].Vp = &v
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	s := fmt.Sprintf("%s\n%s", m.ViewFunc(m.win), m.input.View())
	return s
}

func DefaultResize(width int, height int, ws map[string]*Window) map[string]*Window {
	if w, ok := ws["main"]; ok {
		w.Vp.Width = width
		w.Vp.Height = height
		w.Vp.GotoBottom()
	} else {
		v := viewport.New(width, height)
		ws["main"] = &Window{"", &v}
		ws["main"].Vp.KeyMap = newKeyMap()
		ws["main"].Vp.YPosition = 0 // TOP
		ws["main"].Vp.SetContent(ws["main"].Content)
		ws["main"].Vp.HighPerformanceRendering = false
		ws["main"].Vp.GotoBottom()
	}
	return ws
}

func DefaultView(ws map[string]*Window) string {
	return fmt.Sprintf("%s", ws["main"].Vp.View())
}
