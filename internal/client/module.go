package plugin

import (
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/seandheath/go-mud-client/internal/client"
	"gopkg.in/yaml.v2"
)

//type Module interface {
//Load()
//}

type Window struct {
	viewport.Model
}

type Plugin struct {
	File      string            // Config file name
	Actions   map[string]string `yaml:"actions"` // Actions to be registered
	Aliases   map[string]string `yaml:"aliases"` // Aliases to be registered
	Windows   map[string]string `yaml:"windows"` // Windows to be registered
	windows   map[string]*Window
	functions map[string]client.TriggerFunc
}

// Load a module from a YAML config file
func Load(cfg string) *Plugin {
	file, err := os.ReadFile(cfg)
	if err != nil {
		LogError.Println("Unable to read module config:", err)
		return nil
	}
	m := &Module{}
	err = yaml.Unmarshal(file, m)
	if err != nil {
		LogError.Println("Failed to parse module config:", err)
		return nil
	}
	for re, cmd := range m.Actions {
		AddActionString(re, cmd)
	}
	for re, cmd := range m.Aliases {
		AddAliasString(re, cmd)
	}
	return m
}

func (m Module) LoadActions() {
	for re, cmd := range m.Actions {
		AddActionString(re, cmd)
	}
}
func (m Module) LoadAliases() {
	for re, cmd := range m.Aliases {
		AddAliasString(re, cmd)
	}
}
func (m Module) AddFunction() {
	for re, cmd := range m.Functions {
		AddFunction(re, cmd)
	}
}
