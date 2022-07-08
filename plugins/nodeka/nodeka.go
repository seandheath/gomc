package nodeka

import (
	"fmt"

	"github.com/seandheath/go-mud-client/internal/client"
)

var Config *client.PluginConfig
var Client *client.Client

func Initialize(c *client.Client, file string) *client.PluginConfig {
	Client = c
	cfg, err := client.LoadConfig(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	Config = cfg

	initOmap()

	return Config
}
