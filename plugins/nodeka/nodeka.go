package nodeka

import (
	"fmt"

	"github.com/seandheath/go-mud-client/internal/client"
	"github.com/seandheath/go-mud-client/pkg/plugin"
	"github.com/seandheath/go-mud-client/pkg/trigger"
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

	C.AddAction(`(\[ (?<landed>\d+) of (?<total>\d+) \].+)?(tickl|graz|scratch|bruis|sting|wound|shend|scath|pummel|pummell|batter|splinter|disfigur|fractur|lacerat|RUPTUR|MUTILAT|DEHISC|MAIM|DISMEMBER|SUNDER|CREMAT|EVISCERAT|RAVAG|IMMOLAT|LIQUIFY|LIQUIFI|VAPORIZ|ATOMIZ|OBLITERAT|ETHEREALIZ|ERADICAT)(s|S|e|E|es|ES|ed|ED|ing|ING)? \((?<damage>\d+)\) `, DamageLine)

	return Config
}

func DamageLine(t *trigger.Match) {

}
