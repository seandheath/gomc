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

type window struct {
	content string
	vp      *viewport.Model
}
type model struct {
	ready        bool
	input        textinput.Model
	win          map[string]*window
	inputHistory []string
	inputIndex   int
}

func initialModel() model {
	ti := textinput.New()
	ti.CursorStyle.Blink(false)
	ti.Focus()
	return model{
		input: ti,
		win:   map[string]*window{},
	}
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

func (m model) Init() tea.Cmd { return nil }
func UpdateWindow(w *viewport.Model, msg tea.Msg) tea.Cmd {
	model, cmd := w.Update(msg)
	w = &model
	return cmd
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
			m.win["main"].vp.GotoBottom()
		}
	case tea.WindowSizeMsg:
		inputHeight := lipgloss.Height(m.input.View())
		if !m.ready { // Initial setup
			v := viewport.New(msg.Width, msg.Height-inputHeight)
			m.win["main"] = &window{"", &v}
			m.win["main"].vp.KeyMap = newKeyMap()
			m.win["main"].vp.YPosition = 0 // TOP
			m.win["main"].vp.SetContent(m.win["main"].content)
			m.win["main"].vp.HighPerformanceRendering = false
			m.ready = true
			m.win["main"].vp.GotoBottom()

		} else { // already have a window ready
			m.win["main"].vp.Width = msg.Width
			m.win["main"].vp.Height = msg.Height - inputHeight
			m.win["main"].vp.GotoBottom()
		}
	case showText:
		if w, ok := m.win[msg.window]; ok {
			ab := w.vp.AtBottom()
			w.content += msg.text
			w.vp.SetContent(w.content)
			if ab {
				w.vp.GotoBottom()
			}
		} else {
			go ShowMain(fmt.Sprintf("\nUnable to show text [%s] in window [%s], window not found.\n", msg.text, msg.window))
		}
	}

	v, cmd := m.win["main"].vp.Update(msg)
	m.win["main"].vp = &v
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
func (m model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}
	return fmt.Sprintf("%s\n%s", m.win["main"].vp.View(), m.input.View())
}
