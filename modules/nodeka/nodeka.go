package nodeka

import (
	"log"

	"github.com/olebedev/config"
	"github.com/seandheath/go-mud-client/internal/client"
)

var (
	Triggers                           = map[string]func([]string){}
	isLoaded                           = false
	fmap     map[string]func([]string) = map[string]func([]string){
		"MapLine":   MapLine,
		"EmptyLine": EmptyLine,
	}
)

type Module struct{}

func (m *Module) Load() {
	if isLoaded {
		return
	}
	cfg, err := config.ParseYamlFile("modules/nodeka/nodeka.yaml")
	if err != nil {
		client.ShowMain("Error loading nodeka config: " + err.Error() + "\n")
		return
	}

	// Have to register functions first
	fmap["MapLine"] = MapLine
	fmap["EmptyLine"] = EmptyLine

	for k, v := range fmap {
		client.RegisterFunction(k, v)
	}

	BuffLoad()

	actions, err := cfg.Map("actions")
	if err != nil {
		log.Print("Error loading nodeka config: " + err.Error() + "\n")
		return
	} else {
		for k, v := range actions {
			client.AddAction(k, v)
		}
	}

	aliases, err := cfg.Map("aliases")
	if err != nil {
		log.Print("Error loading nodeka config: " + err.Error() + "\n")
		return
	} else {
		for k, v := range aliases {
			client.AddAlias(k, v)
		}
	}

	isLoaded = true
}
