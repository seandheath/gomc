package autobuff

import (
	"fmt"
	"os"
	"regexp"

	"github.com/seandheath/go-mud-client/internal/client"
	"gopkg.in/yaml.v2"
)

var Config *client.PluginConfig

type Ability struct {
	Mana       int      `yaml:"mana"`       // Mana cost of ability
	Spirit     int      `yaml:"spirit"`     // Spirit cost of ability
	Endurance  int      `yaml:"endurance"`  // Endurance cost of ability
	Prevention string   `yaml:"prevention"` // prevention string
	MinAlign   int      `yaml:"minalign"`   // Must have at least this alignment to use
	MaxAlign   int      `yaml:"maxalign"`   // Must have less than this alignment to use
	IsActive   bool     // Is the ability active?
	Activation []string `yaml:"activation"` // List of regexes that match activation strings
	Execute    string   `yaml:"execute"`    // String to execute the ability
}

type Abilities struct {
	Abilities map[string]*Ability `yaml:"abilities"`
}

var activations = map[string]string{} // Map of activation strings to ability name
var activePreventions = map[string]bool{}
var abilities = map[string]*Ability{} // Map of ability names to ability structs

func Initialize(file string) *client.PluginConfig {
	cfg, err := client.LoadConfig(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	Config = cfg

	b, err := os.ReadFile("plugins/autobuff/abilities.yaml")
	if err != nil {
		client.LogError.Println(err)
	}
	ab := Abilities{}

	err = yaml.Unmarshal(b, &ab)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	abilities = ab.Abilities

	// Go through our wanted buffs and create actions for activation strings
	// Also map the activation strings to the buff names
	for name, buff := range abilities {
		for _, activation := range buff.Activation {
			activations[activation] = name       // Map the activation string
			client.AddAction(activation, BuffUp) // Create the action
		}
	}
	client.AddAction("^You are no longer affected by: (.+)\\.$", BuffDown)
	client.AddAction("^You cannot perform (.+) abilities again yet", PreventUsed)
	client.AddAction("^You may again perform (.+) abilities", PreventAvailable)
	client.AddAlias("^spel$", CheckBuffs)
	return Config
}

var BuffUp client.TriggerFunc = func(re *regexp.Regexp, matches []string) {
	if name, ok := activations[re.String()]; ok { // Get the buff name from the activation string map
		if buff, ok := abilities[name]; ok { // Get the buff from our buff list
			buff.IsActive = true // Set it to active
			if buff.Prevention != "" {
				activePreventions[buff.Prevention] = true
			}
		}
	}
}

// BuffDown handles when a buff drops, preparing it to be cast again.
var BuffDown client.TriggerFunc = func(re *regexp.Regexp, matches []string) {
	if buff, ok := abilities[matches[1]]; ok {
		buff.IsActive = false
	}
}

var PreventUsed client.TriggerFunc = func(re *regexp.Regexp, matches []string) {
	activePreventions[matches[1]] = true
}

var PreventAvailable client.TriggerFunc = func(re *regexp.Regexp, matches []string) {
	activePreventions[matches[1]] = false
}
var CheckBuffs client.TriggerFunc = func(re *regexp.Regexp, matches []string) {
	for name, buff := range abilities {
		if !buff.IsActive && !isPrevented(buff) {
			DoBuff(name, buff)
		}
	}
}

func isPrevented(buff *Ability) bool {
	return activePreventions[buff.Prevention]
}

func DoBuff(name string, buff *Ability) {
	// TODO: Do alignment and pool check
	// TODO: Check prefer invoke etc...
	//if buff.Endurance > 0 { // It's a skill
	//client.Parse(name)
	//} else {
	//if buff.Mana > 0 { // Cast it
	//client.Parse("cast '" + name + "'")
	//} else if buff.Spirit > 0 { // invoke it
	//client.Parse("invoke '" + name + "'")
	//}
	//}
	if buff.Execute == "" {
		client.Parse(name)
	} else {
		client.Parse(buff.Execute)
	}
}
