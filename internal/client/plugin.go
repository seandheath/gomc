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
	Grid      Grid              `yaml:"grid"`
	Windows   map[string]Window `yaml:"windows"` // Windows to be registered
	Functions map[string]TriggerFunc
}

type Window struct {
	Row           int  `yaml:"row"`
	Col           int  `yaml:"col"`
	RowSpan       int  `yaml:"rowspan"`
	ColSpan       int  `yaml:"colspan"`
	MinGridHeight int  `yaml:"mingridheight"`
	MinGridWidth  int  `yaml:"mingridwidth"`
	Border        bool `yaml:"border"`
	Scrollable    bool `yaml:"scrollable"`
	MaxLines      int  `yaml:"maxlines"`
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
		p.Windows = map[string]Window{}
	}
	if p.Functions == nil {
		p.Functions = map[string]TriggerFunc{}
	}
	return p, nil
}

func LoadPlugin(name string, p *PluginConfig) {
	for re, cmd := range p.Actions {
		AddActionString(re, cmd)
	}
	for re, cmd := range p.Aliases {
		AddAliasString(re, cmd)
	}
	for n, f := range p.Functions {
		AddFunction(n, f)
	}
	for n, win := range p.Windows {
		AddWindow(n, win)
	}

	if len(p.Windows) > 0 {
		// Reset input line
		grid.RemoveItem(input)
		grid.AddItem(input,
			len(p.Grid.Rows),    // row
			0,                   // col
			1,                   // rowSpan
			len(p.Grid.Columns), // colSpan
			0,                   // minGridHeight
			0,                   // minGridWidth
			true,                // focus
		)
		if len(p.Grid.Columns) > 0 {
			grid.SetColumns(p.Grid.Columns...)
		}
		if len(p.Grid.Rows) > 0 {
			// Add a row for the input line
			grid.SetRows(append(p.Grid.Rows, 1)...)
		}

	}
	plugins[name] = p
}
