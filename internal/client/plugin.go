package client

import (
	"os"

	"github.com/seandheath/go-mud-client/internal/tui"
	"gopkg.in/yaml.v2"
)

type PluginConfig struct {
	Actions   map[string]string     `yaml:"actions"` // Actions to be registered
	Aliases   map[string]string     `yaml:"aliases"` // Aliases to be registered
	Grid      Grid                  `yaml:"grid"`
	Windows   map[string]tui.Window `yaml:"windows"` // Windows to be registered
	Functions map[string]TriggerFunc
}

type Grid struct {
	Columns []int `yaml:"columns"`
	Rows    []int `yaml:"rows"`
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
		p.Windows = map[string]tui.Window{}
	}
	if p.Functions == nil {
		p.Functions = map[string]TriggerFunc{}
	}
	return p, nil
}

func (c *Client) LoadPlugin(name string, p *PluginConfig) {
	for re, cmd := range p.Actions {
		c.AddActionString(re, cmd)
	}
	for re, cmd := range p.Aliases {
		c.AddAliasString(re, cmd)
	}
	for n, f := range p.Functions {
		c.AddFunction(n, f)
	}
	for n, win := range p.Windows {
		c.tui.AddWindow(n, win)
	}

	if len(p.Windows) > 0 {
		c.tui.FixInputLine(p.Grid.Rows, p.Grid.Columns)
	}
	c.plugins[name] = p
}
