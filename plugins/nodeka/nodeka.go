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
	return Config
}
