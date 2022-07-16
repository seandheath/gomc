package nodeka

import (
	"fmt"
	"strconv"

	"github.com/seandheath/gomc/internal/client"
	"github.com/seandheath/gomc/pkg/plugin"
	"github.com/seandheath/gomc/pkg/trigger"
)

var Config *plugin.Config
var C *client.Client

func Init(c *client.Client, file string) *plugin.Config {
	C = c
	cfg, err := plugin.ReadConfig(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	initOmap()
	abc := initAutobuff()
	Config = plugin.Merge(cfg, abc)

	// Sum damage up and show it at the beginning of the line
	C.AddAction(`(\[ (?P<landed>\d+) of (?P<total>\d+) \].+)?(tickl|graz|scratch|bruis|sting|wound|shend|scath|pummel|pummell|batter|splinter|disfigur|fractur|lacerat|RUPTUR|MUTILAT|DEHISC|MAIM|DISMEMBER|SUNDER|CREMAT|EVISCERAT|RAVAG|IMMOLAT|LIQUIFY|LIQUIFI|VAPORIZ|ATOMIZ|OBLITERAT|ETHEREALIZ|ERADICAT)(s|S|e|E|es|ES|ed|ED|ing|ING)? \((?P<damage>\d+)\) `, DamageLine)
	C.AddAction(`^The closed (?P<door>.+) block\(s\) your passage (?P<direction>.+)\.$`, OpenDoor)

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
