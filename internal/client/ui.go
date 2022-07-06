package client

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type viewData string
type model struct {
	content      string
	ready        bool
	mainView     viewport.Model
	input        textinput.Model
	inputHistory []string
	inputIndex   int
}

func initialModel() model {
	ti := textinput.New()
	ti.CursorStyle.Blink(false)
	ti.Focus()
	return model{
		input:   ti,
		ready:   false,
		content: "", // TODO welcome banner
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
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.mainView, cmd = m.mainView.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

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
			m.mainView.GotoBottom()
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
			m.mainView.GotoBottom()
		}
	case tea.WindowSizeMsg:
		inputHeight := lipgloss.Height(m.input.View())
		if !m.ready {
			m.mainView = viewport.New(msg.Width, msg.Height-inputHeight)
			m.mainView.KeyMap = newKeyMap()
			m.mainView.YPosition = 0 // TOP
			m.mainView.SetContent(m.content)
			m.mainView.HighPerformanceRendering = false
			m.ready = true

		} else { // already have a window ready
			m.mainView.Width = msg.Width
			m.mainView.Height = msg.Height - inputHeight
		}
	case string:
		m.content += msg
		ab := m.mainView.AtBottom()
		m.mainView.SetContent(m.content)
		if ab {
			m.mainView.GotoBottom()
		}
	}

	return m, tea.Batch(cmds...)
}
func (m model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}
	return fmt.Sprintf("%s\n%s", m.mainView.View(), m.input.View())
}
