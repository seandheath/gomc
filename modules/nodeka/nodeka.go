package nodeka

import (
	"log"

	"github.com/olebedev/config"
	"github.com/seandheath/go-mud-client/internal/client"
)

var isLoaded = false

type Module struct {
	Client *client.Client
}

func (m *Module) Load(c *client.Client) {
	if isLoaded {
		return
	}
	m.Client = c
	cfg, err := config.ParseYamlFile("modules/nodeka/nodeka.yaml")
	if err != nil {
		m.Client.ShowMain("Error loading nodeka config: " + err.Error() + "\n")
		return
	}
	actions, err := cfg.Map("actions")
	if err != nil {
		log.Print("Error loading nodeka config: " + err.Error() + "\n")
		return
	} else {
		for k, v := range actions {
			m.Client.AddAction(k, v)
		}
	}
	aliases, err := cfg.Map("aliases")
	if err != nil {
		log.Print("Error loading nodeka config: " + err.Error() + "\n")
		return
	} else {
		for k, v := range aliases {
			m.Client.AddAlias(k, v)
		}
	}
	isLoaded = true
}
