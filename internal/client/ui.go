package client

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var inputStyle = func() lipgloss.Style {
	b := lipgloss.RoundedBorder()
	return lipgloss.NewStyle().BorderStyle(b)
}()

type model struct {
	viewport viewport.Model
	input    textinput.Model
	content  string
	ready    bool
}

func initialModel() model {
	m := model{
		input:   textinput.New(),
		content: "",
		ready:   false,
	}
	m.input.Focus()
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-lipgloss.Height(m.input.View()))
			m.viewport.YPosition = 0 // top of the terminal
			m.viewport.SetContent(m.content)
			m.input.Width = msg.Width
			m.ready = true
		} else {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	// Pass keyboard and mouse events to viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) inputView() string {
	return inputStyle.Render(m.input.View())
}

func (m model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}
	return fmt.Sprintf("%s\n%s", m.viewport.View(), m.input.View())
}

func (m model) ShowMain(text string) { m.Show(text) }
func (m model) Show(text string) {
	m.content += text
}
