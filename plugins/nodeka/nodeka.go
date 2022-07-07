package nodeka

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/seandheath/go-mud-client/internal/client"
)

var Config *client.PluginConfig

func Initialize(file string) *client.PluginConfig {
	cfg, err := client.LoadConfig(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client.AddWindow("main", 100, 60)
	client.AddWindow("omap", 38, 18)
	client.AddWindow("chat", 0, 20)
	client.AddWindow("status", 38, 0)

	client.SetView(NodekaView)
	client.SetResize(NodekaResize)
	Config = cfg

	Config.Functions["MapLine"] = MapLine
	Config.Functions["EmptyLine"] = EmptyLine
	return Config
}

var mainStyle = lipgloss.NewStyle().
	BorderTop(true).
	BorderBottom(true).
	BorderRight(true)

func NodekaResize(width int, height int, ws map[string]*client.Window) map[string]*client.Window {
	ws["chat"].Vp.Height = 18
	ws["chat"].Vp.Width = width
	ws["status"].Vp.Width = 38
	ws["omap"].Vp.Width = 38

	statusWidth := lipgloss.Width(ws["status"].Vp.View())
	chatHeight := lipgloss.Height(ws["chat"].Vp.View())
	omapHeight := lipgloss.Height(ws["omap"].Vp.View())

	// Main Window
	if width < (100 + statusWidth) {
		ws["main"].Vp.Width = width
	} else {
		ws["main"].Vp.Width = 100
	}
	if height < (60 + chatHeight) {
		ws["main"].Vp.Height = height
	} else {
		ws["main"].Vp.Height = height - chatHeight
		ws["chat"].Vp.Width = width
		ws["status"].Vp.Height = height - chatHeight - omapHeight
	}
	return ws
}

func NodekaView(ws map[string]*client.Window) string {
	cv := ws["chat"].Vp.View()
	mv := mainStyle.Render(ws["main"].Content)
	sv := ws["status"].Vp.View()
	ov := ws["omap"].Vp.View()

	doc := strings.Builder{}
	doc.WriteString(cv)
	col := lipgloss.JoinVertical(lipgloss.Right, sv, ov)
	row := lipgloss.JoinHorizontal(lipgloss.Bottom, mv, col)
	doc.WriteString(row + "\n")
	return doc.String()
}
