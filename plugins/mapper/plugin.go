package mapper

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/seandheath/gomc/internal/client"
	"github.com/seandheath/gomc/pkg/plugin"
	"gopkg.in/yaml.v2"
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

func SaveMap(m *Map, path string) {
	data, err := yaml.Marshal(m)
	if err != nil {
		C.Print("\nMAP: Unable to marshal map yaml for saving: " + err.Error())
		return
	}
	f, err := os.Create(path)
	if err != nil {
		C.Print("\nMAP: Unable to save map data to file: " + err.Error())
		return
	}
	defer f.Close()

	w := gzip.NewWriter(f)
	defer w.Close()

	_, err = w.Write(data)
	if err != nil {
		C.Print("\nMAP: Error writing map to file: " + err.Error())
	}
	C.Print("\nMAP: Saved map to file: " + path)
}

func (m *Map) Load(path string) {
	f, err := os.Open(path)
	if err != nil {
		C.Print("\nMAP: Unable to load map data from file: " + err.Error())
	}
	defer f.Close()

	r, err := gzip.NewReader(f)
	if err != nil {
		C.Print("\nMAP: Unable to decompress map data from file: " + err.Error())
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		C.Print("\nMAP: Unable to read map data from file: " + err.Error())
	}

	err = yaml.Unmarshal(data, m)
	if err != nil {
		C.Print("\nMAP: Unable to unmarshal map data from file contents: " + err.Error())
	}
	m.Rebuild()
	C.Print("\nMAP: Loaded " + path)
}
