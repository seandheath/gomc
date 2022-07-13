package nodeka

import (
	"fmt"

	"github.com/seandheath/go-mud-client/internal/client"
	"github.com/seandheath/go-mud-client/pkg/plugin"
)

var Config *plugin.Config
var Client *client.Client

func Initialize(c *client.Client, file string) *plugin.Config {
	Client = c
	cfg, err := plugin.ReadConfig(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	Config = cfg

	initOmap()

	return Config
}
