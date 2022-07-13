package mapper

import (
	"fmt"
	"math"

	"github.com/seandheath/go-mud-client/internal/client"
	"github.com/seandheath/go-mud-client/pkg/plugin"
	"github.com/seandheath/go-mud-client/pkg/trigger"
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
	Up
	Down
	Look
)

type Coordinates struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

type Room struct {
	id          int                 `yaml:"id"`          // Unique ID
	name        string              `yaml:"name"`        // Room name as seen in the game
	exits       map[Direction]*Room `yaml:"exits"`       // Exits out of this room
	area        *Area               `yaml:"area"`        // The parent area of this room
	coordinates Coordinates         `yaml:"coordinates"` // The coordinates of this room
	flags       map[string]string   `yaml:"flags"`       // List of key:value flags for this room
	tag         string              `yaml:"tag"`         // A unique tag for the room
}

type Area struct {
	Name  string        `yaml:"name"`
	Rooms map[int]*Room `yaml:"rooms"`
}

type Map struct {
	room      *Room // Current room
	area      *Area // Current area
	exits     map[Direction]bool
	rname     string
	Areas     map[string]*Area `yaml:"areas"` // All areas in the map
	Rooms     map[int]*Room    `yaml:"rooms"` // All rooms in the map
	moving    bool
	movequeue []Direction
	movedir   Direction
}

var dirmap map[string]Direction = map[string]Direction{
	"n":     North,
	"north": North,
	"e":     East,
	"east":  East,
	"s":     South,
	"south": South,
	"w":     West,
	"west":  West,
	"u":     Up,
	"up":    Up,
	"d":     Down,
	"down":  Down,
	"lo":    Look,
	"loo":   Look,
	"look":  Look,
	"map":   Look,
}
var reverse map[Direction]Direction = map[Direction]Direction{
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
	addCommands(C, M)
	C.AddFunction("MoveHappening", M.MoveHappening)
	C.AddFunction("MoveDone", M.MoveDone)

	Config = cfg
	return Config
}
func NewMap() *Map {
	m := &Map{}
	m.Areas = make(map[string]*Area)
	m.Rooms = make(map[int]*Room)
	m.area = &Area{}
	m.room = &Room{}
	C.Print("\nMap created. Start walking around to add rooms.")
	return m
}
func (m *Map) Load(path string) *Map  { return nil }
func (m *Map) Save(path string) error { return nil }

// getNewID steps through all the rooms in the map and identifies the lowest
// available ID integer. Room IDs are reused after rooms are deleted.
func (m *Map) getNewID() int {
	min := math.MaxInt
	max := 0
	for i, r := range m.Rooms {
		if i < min && r == nil {
			min = i
		}
		if i > max {
			max = i + 1
		}
	}
	if min < max {
		return min
	}
	return max
}

func (m *Map) AddRoom(name string, exits string, area string) {
	if m.area == nil {
		C.Print("\nCurrent area is not defined, set it with '#map area <area>'")
	}
	if m.room == nil {
		C.Print("\nCurrent room is not defined, set it with '#map goto <name|tag>'")
	}
}
func (m *Map) FindRoom(name string, exits string) *Room        { return nil }
func (m *Map) ShiftRoom(room *Room, direction Direction) *Room { return nil }
func (m *Map) AddArea(name string)
func (m *Map) Show(width int, height int) string { return "" }
func (m *Map) StartMove(direction string) {
	m.moving = true
}

func (m *Map) MoveHappening(t *trigger.Trigger) {
	if m.moving {
		// TODO Check Room
		m.rname = t.Matches[1]
		m.exits[North] = t.Matches[2] == "north"
		m.exits[East] = t.Matches[3] == "east"
		m.exits[South] = t.Matches[4] == "south"
		m.exits[West] = t.Matches[5] == "west"
		m.exits[Up] = t.Matches[6] == "up"
		m.exits[Down] = t.Matches[7] == "down"
	}
}
func (m *Map) MoveDone(t *trigger.Trigger) {}
func (m *Map) queueMove(direction Direction) {
	m.moving = true
	m.movequeue = append(m.movequeue, direction)
}
