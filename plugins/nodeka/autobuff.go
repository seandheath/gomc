package nodeka

import (
	"fmt"
	"log"
	"os"

	"github.com/seandheath/go-mud-client/pkg/plugin"
	"github.com/seandheath/go-mud-client/pkg/trigger"
	"gopkg.in/yaml.v2"
)

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

func initAutobuff() *plugin.Config {
	cfg, err := plugin.ReadConfig("plugins/nodeka/autobuff.yaml")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	b, err := os.ReadFile("plugins/nodeka/abilities.yaml")
	if err != nil {
		C.LogError.Println(err)
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
			activations[activation] = name  // Map the activation string
			C.AddAction(activation, BuffUp) // Create the action
		}
	}
	C.AddAction("^You are no longer affected by: (.+)\\.$", BuffDown)
	C.AddAction("^You cannot perform (.+) abilities again yet", PreventUsed)
	C.AddAction("^You may again perform (.+) abilities", PreventAvailable)
	C.AddAction("^Your botanswer is: autobuff done", func(t *trigger.Trigger) { attempting = false })
	C.AddAlias("^spel$", CheckBuffs)
	C.AddAlias("^abon$", AutobuffOn)
	C.AddAlias("^aboff$", AutobuffOff)
	return cfg
}

var autoBuffOn = false

func AutobuffOn(t *trigger.Trigger) {
	autoBuffOn = true
}
func AutobuffOff(t *trigger.Trigger) {
	autoBuffOn = false
}

var BuffUp trigger.Func = func(t *trigger.Trigger) {
	if name, ok := activations[t.String()]; ok { // Get the buff name from the activation string map
		if buff, ok := abilities[name]; ok { // Get the buff from our buff list
			buff.IsActive = true // Set it to active
			if buff.Prevention != "" {
				activePreventions[buff.Prevention] = true
			}
		}
	}
}

// BuffDown handles when a buff drops, preparing it to be cast again.
var BuffDown trigger.Func = func(t *trigger.Trigger) {
	if buff, ok := abilities[t.Matches[1]]; ok {
		buff.IsActive = false
	}
	if autoBuffOn {
		CheckBuffs(nil)
	}
}

var PreventUsed trigger.Func = func(t *trigger.Trigger) {
	activePreventions[t.Matches[1]] = true
}

var PreventAvailable trigger.Func = func(t *trigger.Trigger) {
	if autoBuffOn {
		CheckBuffs(nil)
	}
	activePreventions[t.Matches[1]] = false
}
var attempting = false
var CheckBuffs trigger.Func = func(t *trigger.Trigger) {
	for name, buff := range abilities {
		if !buff.IsActive && !isPrevented(buff) {
			if !attempting {
				DoBuff(name, buff)
				C.SendNow("botanswer autobuff done")
				attempting = true
			}
		}
	}
	if attempting {
	}
}

func isPrevented(buff *Ability) bool {
	return activePreventions[buff.Prevention]
}

func DoBuff(name string, buff *Ability) {
	// TODO: Do alignment and pool check
	// TODO: Check prefer invoke etc...
	if buff.Execute == "" {
		C.Parse(name)
	} else {
		C.Parse(buff.Execute)
	}
}
