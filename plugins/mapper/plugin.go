package mapper

import (
	"fmt"

	"github.com/seandheath/gomc/internal/client"
	"github.com/seandheath/gomc/pkg/plugin"
)

type Direction string

const (
	North  Direction = "north"
	East   Direction = "east"
	South  Direction = "south"
	West   Direction = "west"
	Up     Direction = "up"
	Down   Direction = "down"
	Look   Direction = "look"
	Recall Direction = "recall"
)

var dirmap = map[string]Direction{
	"n":      North,
	"north":  North,
	"e":      East,
	"east":   East,
	"s":      South,
	"south":  South,
	"w":      West,
	"west":   West,
	"u":      Up,
	"up":     Up,
	"d":      Down,
	"down":   Down,
	"lo":     Look,
	"loo":    Look,
	"look":   Look,
	"rec":    Recall,
	"reca":   Recall,
	"recal":  Recall,
	"recall": Recall,
}

var reverse = map[Direction]Direction{
	North: South,
	East:  West,
	South: North,
	West:  East,
	Up:    Down,
	Down:  Up,
}

var C *client.Client
var M *Map
var Config *plugin.Config

func Init(cli *client.Client, file string) *plugin.Config {
	C = cli
	cfg, err := plugin.ReadConfig(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	M = NewMap()

	// TODO load map from yaml
	addCommands(C, M)
	C.AddFunction("MoveDone", M.MoveDone)
	C.AddFunction("MoveFail", M.MoveFail)
	C.AddFunction("MoveRecall", M.MoveRecall)
	C.AddFunction("MoveClear", M.MoveClear)

	C.AddActionFunc(`^(?P<name>.+) \[ exits: (?P<exits>\(?(north)?\)? ?\(?(east)?\)? ?\(?(south)?\)? ?\(?(west)?\)? ?\(?(up)?\)? ?\(?(down)?\)?) ?\]$`, M.MoveDone)

	Config = cfg
	return Config
}
