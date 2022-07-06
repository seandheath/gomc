package client

import (
	"os"

	"gopkg.in/yaml.v2"
)

type PluginInterface interface {
	Initialize(*PluginConfig) error
}
type PluginConfig struct {
	PluginInterface
	Actions   map[string]string `yaml:"actions"` // Actions to be registered
	Aliases   map[string]string `yaml:"aliases"` // Aliases to be registered
	Windows   map[string]string `yaml:"windows"` // Windows to be registered
	Functions map[string]TriggerFunc
}

// Load a module from a YAML config file
func LoadConfig(cfg string) (*PluginConfig, error) {
	file, err := os.ReadFile(cfg)
	if err != nil {
		return nil, err
	}
	p := &PluginConfig{}
	err = yaml.Unmarshal(file, p)
	if err != nil {
		return nil, err
	}
	if p.Actions == nil {
		p.Actions = map[string]string{}
	}
	if p.Aliases == nil {
		p.Aliases = map[string]string{}
	}
	if p.Windows == nil {
		p.Windows = map[string]string{}
	}
	if p.Functions == nil {
		p.Functions = map[string]TriggerFunc{}
	}
	return p, nil
}
