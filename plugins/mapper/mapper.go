package mapper

import (
	"fmt"
	"math"
	"strings"

	"github.com/seandheath/gomc/internal/client"
	"github.com/seandheath/gomc/pkg/plugin"
	"github.com/seandheath/gomc/pkg/trigger"
)

type Direction string

const (
	North Direction = "north"
	East  Direction = "east"
	South Direction = "south"
	West  Direction = "west"
	Up    Direction = "up"
	Down  Direction = "down"
)

type Coordinates struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

type Room struct {
	id               int                 `yaml:"id"`    // Unique ID
	name             string              `yaml:"name"`  // Room name as seen in the game
	exits            map[Direction]*Room `yaml:"exits"` // Exits out of this room
	exitString       string
	area             *Area             `yaml:"area"`        // The parent area of this room
	coordinates      Coordinates       `yaml:"coordinates"` // The coordinates of this room
	tags             map[string]string `yaml:"tags"`        // List of key:value flags for this room
	extraIdentifiers map[string]string `yaml:"identifier"`  // A list of key:value identifiers for this room
}

type Area struct {
	Name  string        `yaml:"name"`
	Rooms map[int]*Room `yaml:"rooms"`
}

type Map struct {
	home         *Room            // Room the player recalls to
	room         *Room            // Current room
	area         *Area            // Current area
	Areas        map[string]*Area `yaml:"areas"` // All areas in the map
	Rooms        map[int]*Room    `yaml:"rooms"` // All rooms in the map
	nextMoves    []Direction
	pastmoves    []Direction
	rmExitString string // rmExits are updated every time we see an exit string - even if we're not expecting a move
	rmName       string // rmName is updated every time we see a new room name - even if we're not expecting a move
	sawMove      bool
	mapping      bool
}

var dirmap = map[string]Direction{
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
	addCommands(C, M)
	C.AddFunction("moveHappening", M.moveHappening)
	C.AddFunction("moveDone", M.moveDone)

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

func (m *Map) moveStart(move Direction) {
	m.nextMoves = append(m.nextMoves, move)
}
func (m *Map) moveHappening(t *trigger.Trigger) {
	if len(m.nextMoves) > 0 {
		m.sawMove = true
	}
}

// checkMove executes when a trigger indicating a move is complete fires
// this trigger needs to be set somewhere and correspond to a string
// such as a prompt that you see after all room information scrolls
// past
func (m *Map) moveDone(t *trigger.Trigger) {
	m.rmName = t.Results["name"]
	m.rmExitString = t.Results["exits"]
	if len(m.nextMoves) > 0 && m.sawMove {
		m.sawMove = false
		// We should have all the information about the new room - we can check
		// if it matches what we expect to see. If it matches then we'll move
		// the map to the new room.
		nr := m.checkMove(m.room, m.nextMoves[0])
		if nr != nil {
			m.room = nr
		} else {
			// move failed... do something!
		}
	}
}

func (m *Map) checkMove(from *Room, move Direction) *Room {
	if from == nil {
		return nil
	}
	// This room doesn't have an exit in that direction...
	if r, ok := m.room.exits[move]; !ok {
		return nil
	} else if r == nil {
		// We have an exit but don't know what room it points to, let's add one
		if m.mapping {
			// we're currently mapping so lets create a new room
			r = m.addRoomFromMove(move)
		}
		return r
	} else {
		if m.checkRoom(r) {
			// it's the same room
			return r
		} else {
			return nil
		}
	}
}

// checkRoom compares the provided room to the room most recently seen
// by the map. If they match, we'll return true.
func (m *Map) checkRoom(r *Room) bool {
	// Currently only checking name and exit string
	if (r.name == m.rmName) && (r.exitString == m.rmExitString) {
		return true
	}
	return false
}

// While moving around the map we'll add rooms as we go.
func (m *Map) addRoomFromMove(move Direction) *Room {
	if m.room == nil {
		// I don't know where I am so I can't link a room
		// TODO print some message here
		C.Print("\nCan't create a room, I don't know where I am. Set your current room with '#map goto <roomID>'")
		return nil
	}
	// We have a room that we are coming from, let's add a room to it
	nr := &Room{}
	//c := m.getCoordinatesFromDir(m.room.coordinates, move)

	// Get the exit strings
	es := strings.Split(m.rmExitString, " ")
	for _, e := range es {
		d := dirmap[e]
		nr.exits[d] = nil // TODO get the room at the coordinates
	}
	return nil
}
func (m *Map) getCoordinatesFromDir(c Coordinates, move Direction) Coordinates {
	switch move {
	case North: // Positive Y axis
		return Coordinates{c.X, c.Y + 1, c.Z}
	case East: // Positive X axis
		return Coordinates{c.X + 1, c.Y, c.Z}
	case South: // Negative Y axis
		return Coordinates{c.X, c.Y - 1, c.Z}
	case West: // Negative X axis
		return Coordinates{c.X - 1, c.Y, c.Z}
	case Up: // Positive Z axis
		return Coordinates{c.X, c.Y, c.Z + 1}
	case Down: // Negative Z axis
		return Coordinates{c.X, c.Y, c.Z - 1}
	}
	// Shouldn't get here
	C.Print("Error creating coordinates")
	return Coordinates{}
}
