package nodeka

import (
	"log"

	"github.com/olebedev/config"
	"github.com/seandheath/go-mud-client/internal/client"
)

var (
	isLoaded = false
	Client   *client.Client
	fmap     map[string]func(string) = map[string]func(string){
		"MapLine":   MapLine,
		"EmptyLine": EmptyLine,
	}
)

type Module struct{}

func (m *Module) Load(c *client.Client) {
	if isLoaded {
		return
	}
	Client = c
	cfg, err := config.ParseYamlFile("modules/nodeka/nodeka.yaml")
	if err != nil {
		Client.ShowMain("Error loading nodeka config: " + err.Error() + "\n")
		return
	}
	actions, err := cfg.Map("actions")
	if err != nil {
		log.Print("Error loading nodeka config: " + err.Error() + "\n")
		return
	} else {
		for k, v := range actions {
			Client.AddAction(k, v)
		}
	}
	aliases, err := cfg.Map("aliases")
	if err != nil {
		log.Print("Error loading nodeka config: " + err.Error() + "\n")
		return
	} else {
		for k, v := range aliases {
			Client.AddAlias(k, v)
		}
	}

	fmap["MapLine"] = MapLine
	fmap["EmptyLine"] = EmptyLine

	for k, v := range fmap {
		Client.RegisterFunction(k, v)
	}
	isLoaded = true
}
