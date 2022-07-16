package nodeka

import (
	"github.com/seandheath/gomc/pkg/trigger"
)

var shadowKill bool
var shadowAttack = "f "
var shadowTarget string
var shadowTrigger *trigger.Trigger

func initShadow() {
	shadowTrigger = C.AddActionFunc(`^Shadowing, you sense the way is \[ (?P<dir>north|east|south|west|up|down|right here)`, shadowDir)
	shadowTrigger.Enabled = false

	C.AddActionFunc(`^You end your shadow tactics`, shadowOff)
	C.AddAliasFunc(`^sk (?P<target>.+)$`, shadowHunt)
	C.AddAliasFunc(`^sf (?P<target>.+)$`, shadowFind)
}

func shadowHunt(t *trigger.Trigger) {
	shadowKill = true
	shadowFind(t)
}

func shadowOff(t *trigger.Trigger) {
	shadowKill = false
	shadowTrigger.Enabled = false
}

func shadowFind(t *trigger.Trigger) {
	shadowTrigger.Enabled = true
	shadowTarget = t.Results["target"]
	C.Parse("shadow " + shadowTarget)
}

func shadowDir(t *trigger.Trigger) {
	switch t.Results["dir"] {
	case "north", "south", "east", "west", "up", "down":
		C.Parse(t.Results["dir"])
	case "right here":
		if shadowKill {
			C.Parse(shadowAttack + shadowTarget)
			DeadQ.Append(func() {
				C.Parse("look")
			})
		} else {
			shadowTrigger.Enabled = false
		}
	}
}
