package nodeka

import (
	"fmt"
	"log"
	"os"

	"github.com/seandheath/gomc/pkg/plugin"
	"github.com/seandheath/gomc/pkg/trigger"
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
	Triggers   []*trigger.Trigger
}

type Abilities struct {
	Buffs     []string            `yaml:"buffs"`
	Combos    map[string][]string `yaml:"combos"`
	Abilities map[string]*Ability `yaml:"abilities"`
}

var activations = map[string]string{} // Map of activation strings to ability name
var activePreventions = map[string]bool{}
var abilities = map[string]*Ability{} // Map of ability names to ability structs
var combo = []string{}
var buffs = []string{}
var Attempt *Ability = nil

func initAbilities() *plugin.Config {
	cfg, err := plugin.ReadConfig("plugins/nodeka/autobuff.yaml")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	b, err := os.ReadFile("plugins/nodeka/abilities.yaml")
	if err != nil {
		fmt.Println(err)
	}
	ab := Abilities{}

	err = yaml.Unmarshal(b, &ab)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	abilities = ab.Abilities
	combo = ab.Combos["default"]
	buffs = ab.Buffs

	// Go through our wanted buffs and create actions for activation strings
	// Also map the activation strings to the buff names
	for name, ab := range abilities {
		for _, activation := range ab.Activation {
			activations[activation] = name                    // Map the activation string
			t := trigger.NewTrigger(activation, AbilityFired) // Create a trigger for the activation string
			C.AddActionTrigger(t)                             // Create the action
			ab.Triggers = append(ab.Triggers, t)
		}
	}
	C.AddAction("^You are no longer affected by: (.+)\\.$", AbilityDown)
	C.AddAction("^You cannot perform (.+) abilities again yet", PreventUsed)
	C.AddAction("^You may again perform (.+) abilities", PreventAvailable)
	C.AddAction(`^You lose your focus\.$`, AbilityFailed)
	C.AddAlias("^spel$", DoBuffs)
	C.AddAlias("^abon$", AutobuffOn)
	C.AddAlias("^aboff$", AutobuffOff)

	return cfg
}

var autoBuffOn = true

func AutobuffOn(t *trigger.Trigger) {
	autoBuffOn = true
}
func AutobuffOff(t *trigger.Trigger) {
	autoBuffOn = false
}

var AbilityFired trigger.Func = func(t *trigger.Trigger) {
	if name, ok := activations[t.String()]; ok { // Get the ab name from the activation string map
		if ab, ok := abilities[name]; ok { // Get the buff from our buff list
			ab.IsActive = true // Set it to active
			if ab.Prevention != "" {
				activePreventions[ab.Prevention] = true
			}
			if Attempt == ab {
				// We were attempting this buff, we saw it so we clear attempt
				Attempt = nil
				//DoCombo()
			}
		}
	}
}

// AbilityDown handles when a buff drops, preparing it to be cast again.
var AbilityDown trigger.Func = func(t *trigger.Trigger) {
	if buff, ok := abilities[t.Matches[1]]; ok {
		buff.IsActive = false
	}
	//if autoBuffOn {
	//ReplyQ.Prepend(func() { DoBuffs(nil) })
	//}
}

var PreventUsed trigger.Func = func(t *trigger.Trigger) {
	AbilityFailed(nil)
	activePreventions[t.Matches[1]] = true
}

var PreventAvailable trigger.Func = func(t *trigger.Trigger) {
	activePreventions[t.Matches[1]] = false
	//if autoBuffOn {
	//ReplyQ.Prepend(func() { DoBuffs(nil) })
	//}
}

func DoBuffs(t *trigger.Trigger) {
	for _, an := range buffs {
		if ab, ok := abilities[an]; ok {
			if !ab.IsActive && !isPrevented(ab) {
				DoAbility(an, ab)
			}
		} else {
			C.Print("Ability not found: " + an)
		}
	}
}

func isPrevented(buff *Ability) bool {
	return activePreventions[buff.Prevention]
}

func AbilityFailed(t *trigger.Trigger) {
	if Attempt != nil {
		Attempt = nil
		DoCombo()
	}
}

func DoAbility(name string, ab *Ability) *Ability {
	// TODO: Do alignment and pool check
	// TODO: Check prefer invoke etc...
	if ab.Execute == "" {
		C.Parse(name)
	} else {
		C.Parse(ab.Execute)
	}
	return ab
}

func DoCombo() {
	if Attempt == nil {
		for _, name := range combo {
			if ab, ok := abilities[name]; ok {
				if !isPrevented(ab) {
					Attempt = ab
					DoAbility(name, ab)
					return
				}
			} else {
				C.Print("Ability not found: " + name)
			}
		}
	}
}
