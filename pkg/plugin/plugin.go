package plugin

import (
	"os"

	"github.com/seandheath/gomc/pkg/trigger"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Actions        map[string]string `yaml:"actions"` // Actions to be registered
	Aliases        map[string]string `yaml:"aliases"` // Aliases to be registered
	Grid           Grid              `yaml:"grid"`
	Windows        map[string]Window `yaml:"windows"` // Windows to be registered
	Functions      map[string]trigger.Func
	CredentialFile string `yaml:"credentialfile"`
	AutoLogin      string `yaml:"autologin"`
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

// Load a plugin from a YAML config file
func ReadConfig(cfg string) (*Config, error) {
	file, err := os.ReadFile(cfg)
	if err != nil {
		return nil, err
	}
	p := &Config{}
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
		p.Functions = map[string]trigger.Func{}
	}
	return p, nil
}

// woohoo generics!
func merge[V string | trigger.Func](a map[string]V, b map[string]V) map[string]V {
	for k, v := range b {
		a[k] = v
	}
	return a
}
func Merge(a *Config, b *Config) *Config {
	merge(a.Actions, b.Actions)
	merge(a.Aliases, b.Aliases)
	merge(a.Functions, b.Functions)
	return a
}
