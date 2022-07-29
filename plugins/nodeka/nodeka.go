package nodeka

import (
	"fmt"
	"strconv"

	"github.com/seandheath/gomc/internal/client"
	"github.com/seandheath/gomc/pkg/plugin"
	"github.com/seandheath/gomc/pkg/queue"
	"github.com/seandheath/gomc/pkg/trigger"
)

type Position string

const (
	Sleep Position = "sleep"
	Rest  Position = "rest"
	Stand Position = "stand"
	Fight Position = "fight"
)

type Character struct {
	Exp             int
	Align           int
	Gold            int
	Lag             int
	CurrentHP       int
	MaxHP           int
	CurrentMana     int
	MaxMana         int
	CurrentSpirit   int
	MaxSpirit       int
	CurrentEnd      int
	MaxEnd          int
	PKFlag          bool
	Tank            string
	TankCondition   string
	Target          string
	TargetCondition string
	Position        Position
}

var Config *plugin.Config
var C *client.Client
var ReplyQ *queue.Queue
var DeadQ *queue.Queue
var My Character

func Init(c *client.Client, file string) *plugin.Config {
	My = Character{}

	C = c
	cfg, err := plugin.ReadConfig(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	initOmap()
	initFootpad()
	initBot()
	ReplyQ = queue.NewQueue()
	//DeadQ = trigger.NewQueue(`is dead!$`)
	//C.AddActionTrigger(ReplyQ.Trigger)
	//C.AddActionTrigger(DeadQ.Trigger)

	abc := initAbilities()
	Config = plugin.Merge(cfg, abc)

	// Sum damage up and show it at the beginning of the line
	C.AddActionTrigger(trigger.NewTrigger(`(\[ (?P<landed>\d+) of (?P<total>\d+) \].+)?(tickl|graz|scratch|bruis|sting|wound|shend|scath|pummel|pummell|batter|splinter|disfigur|fractur|lacerat|RUPTUR|MUTILAT|DEHISC|MAIM|DISMEMBER|SUNDER|CREMAT|EVISCERAT|RAVAG|IMMOLAT|LIQUIFY|LIQUIFI|VAPORIZ|ATOMIZ|OBLITERAT|ETHEREALIZ|ERADICAT)(s|S|e|E|es|ES|ed|ED|ing|ING)? \((?P<damage>\d+)\) `, DamageLine))
	C.AddActionTrigger(trigger.NewTrigger(`^The closed (?P<door>.+) block\(s\) your passage (?P<direction>.+)\.$`, OpenDoor))
	C.AddFunction("ReplyPrompt", ReplyPrompt)
	C.AddFunction("PoolPrompt", PoolPrompt)
	C.AddFunction("CombatPrompt", CombatPrompt)
	C.AddFunction("AbilityFailed", AbilityFailed)

	return Config
}

func OpenDoor(t *trigger.Trigger) {
	C.Parse("open " + t.Results["door"])
	C.Parse(t.Results["direction"])
}

func DamageLine(t *trigger.Trigger) {
	var td int
	d, err := strconv.Atoi(t.Results["damage"])
	if err != nil {
		return
	}
	if t.Matches[2] == "" {
		td = d
	} else {
		l, err := strconv.Atoi(t.Results["landed"])
		if err != nil {
			return
		}
		// TODO: log accuracy
		//t, err := strconv.Atoi(t.Results["total"])
		//if err != nil {
		//return
		//}
		td = l * d
	}
	// Add the damage string to the beginning of the line
	C.RawLine = append([]byte(fmt.Sprintf("[ %d ] ", td)), C.RawLine...)
}

func getResult(pool string, results map[string]string) int {
	if val, ok := results[pool]; ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0
		}
		return i
	}
	return 0
}
func PoolPrompt(t *trigger.Trigger) {
	My.CurrentHP = getResult("chp", t.Results)
	My.MaxHP = getResult("mhp", t.Results)
	My.CurrentMana = getResult("cm", t.Results)
	My.MaxMana = getResult("mm", t.Results)
	My.CurrentSpirit = getResult("cs", t.Results)
	My.MaxSpirit = getResult("ms", t.Results)
	My.CurrentEnd = getResult("ce", t.Results)
	My.MaxEnd = getResult("me", t.Results)
}

func ReplyPrompt(t *trigger.Trigger) {
	if My.Position == Fight {
		// Just finished fighting
		My.Position = Stand
		Attempt = nil
	}
	My.Exp = getResult("exp", t.Results)
	My.Gold = getResult("gold", t.Results)
	My.Align = getResult("align", t.Results)
	My.Lag = getResult("lag", t.Results)
	if _, ok := t.Results["pk"]; ok {
		My.PKFlag = true
	} else {
		My.PKFlag = false
	}
	ReplyQ.Do()
}

func CombatPrompt(t *trigger.Trigger) {
	My.Position = Fight
	My.Tank = t.Results["tank"]
	My.TankCondition = t.Results["tankcond"]
	My.Target = t.Results["target"]
	My.TargetCondition = t.Results["targetcond"]
	My.Lag = getResult("lag", t.Results)
	if My.Lag < 2000 {
		DoCombo()
	}
}
