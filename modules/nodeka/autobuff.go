package nodeka

import (
	"github.com/seandheath/go-mud-client/internal/client"
)

type cost struct {
	Health    int
	Mana      int
	Spirit    int
	Endurance int
}
type Ability struct {
	Name        string   // Name of the ability REQUIRED
	Cost        cost     // Ability cost in mana, spirit, or endurance
	Prevention  string   // prevention string
	MinAlign    int      // Must have at least this alignment to use
	MaxAlign    int      // Must have less than this alignment to use
	IsActive    bool     // Is the ability active?
	IsPrevented bool     // Is the ability prevented?
	Activation  []string // List of regexes that match activation strings
}

var AllAbilities = map[string]Ability{}     // List of all abilities
var Buffs = map[string]Ability{}            // List of buffs to manage
var Attacks = map[string]Ability{}          // List of attacks to manage
var activationStrings = map[string]string{} // Map of activation strings to ability name

// BuffLoad initializes all the actions for buffs
func BuffLoad() {
	// Go through our wanted buffs and create actions for activation strings
	// Also map the activation strings to the buff names
	for name, buff := range Buffs {
		for _, activation := range buff.Activation {
			activationStrings[activation] = name // Map the activation string
			client.AddAction(activation, BuffUp) // Create the action
		}
	}
	client.AddAction("^You are no longer affected by: (.+)\\.$", BuffDown)
}

// BuffDown handles when a buff drops, preparing it to be cast again.
func BuffDown(args []string) {
	if buff, ok := Buffs[args[1]]; ok {
		buff.IsActive = false
	}
}

func BuffUp(args []string) {
	if name, ok := activationStrings[args[0]]; ok { // Get the buff name from the activation string map
		if buff, ok := Buffs[name]; ok { // Get the buff from our buff list
			buff.IsActive = true // Set it to active
		}
	}
}

func NewAbility(n string) *Ability {
	return &Ability{
		Name:       n,
		Prevention: "",
		Cost:       cost{0, 0, 0, 0}, // health, mana, spirit, endurance
		MinAlign:   -1001,            // No restriction
		MaxAlign:   1001,             // No restriction
	}
}
