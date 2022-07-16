package nodeka

import "github.com/seandheath/gomc/pkg/trigger"

func initFootpad() {
	initShadow()
	C.AddAliasFunc(`^f (?P<target>.+)$`, crit)
	C.AddActionFunc(`^You are unable to locate a weak spot`, critFail)
}

var critTarget = ""

func crit(t *trigger.Trigger) {
	critTarget = t.Results["target"]
	C.Parse("critical attack " + critTarget)
}
func critFail(t *trigger.Trigger) {
	C.Parse("critical attack " + critTarget)
}
