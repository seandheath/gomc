package nodeka

import (
	"fmt"

	"github.com/seandheath/go-mud-client/internal/client"
)

var Config *client.PluginConfig

func Initialize(file string) *client.PluginConfig {
	cfg, err := client.LoadConfig(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	Config = cfg

	Config.Functions["MapLine"] = MapLine
	Config.Functions["EmptyLine"] = EmptyLine

	//windows := make(map[string]*tview.TextView)

	//windows["main"] = tview.NewTextView().
	//SetDynamicColors(true)
	//windows["chat"] = tview.NewTextView().
	//SetDynamicColors(true)
	//windows["overhead"] = tview.NewTextView().
	//SetDynamicColors(true).
	//SetScrollable(false).
	//SetMaxLines(16)
	//input = tview.NewInputField().
	//SetDoneFunc(handleInput)
	//grid = tview.NewGrid().
	//SetColumns(0, 40).
	//SetRows(16, 0, 1).
	//SetBorders(true).
	//AddItem(windows["chat"], 0, 0, 1, 1, 16, 0, false).
	//AddItem(windows["overhead"], 0, 1, 1, 1, 16, 40, false).
	//AddItem(windows["main"], 1, 0, 1, 2, 0, 0, false).
	//AddItem(input, 2, 0, 1, 2, 1, 0, true)
	return Config
}
