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
	North  Direction = "north"
	East   Direction = "east"
	South  Direction = "south"
	West   Direction = "west"
	Up     Direction = "up"
	Down   Direction = "down"
	Look   Direction = "look"
	Recall Direction = "recall"
)

type Coordinates struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

type Room struct {
	id          int                 `yaml:"id"`    // Unique ID
	name        string              `yaml:"name"`  // Room name as seen in the game
	exits       map[Direction]*Room `yaml:"exits"` // Exits out of this room
	exitString  string
	area        *Area             `yaml:"area"`        // The parent area of this room
	coordinates Coordinates       `yaml:"coordinates"` // The coordinates of this room
	tags        map[string]string `yaml:"tags"`        // List of key:value flags for this room
	//extraIdentifiers map[string]string `yaml:"identifier"`  // A list of key:value identifiers for this room
}

type Area struct {
	Name  string  `yaml:"name"`
	Rooms []*Room `yaml:"rooms"`
}

type Map struct {
	loaded bool
	room   *Room // Current room
	//area         *Area            // Current area not sure I'm gonna use this
	areas        []*Area     `yaml:"areas"` // All areas in the map
	rooms        []*Room     `yaml:"rooms"` // All rooms in the map
	nextMoves    []Direction // queued up moves
	pastMoves    []Direction // used to track last steps to determine location if we get lost
	rmExitString string      // rmExits are updated every time we see an exit string - even if we're not expecting a mov
	rmName       string      // rmName is updated every time we see a new room name - even if we're not expecting a move
	sawMove      bool
	mapping      bool // do not create or link rooms if we're not in mapping mode, puts map in read-only state
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
	"lo":    Look,
	"loo":   Look,
	"look":  Look,
	//"rec":    Recall,
	//"reca":   Recall,
	//"recal":  Recall,
	//"recall": Recall,
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

var moveDoneTrigger *trigger.Trigger

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
	C.AddFunction("moveDone", M.moveDone)
	C.AddFunction("moveFail", M.moveFail)

	moveDoneTrigger = C.AddActionFunc(`^(?P<name>.+) \[ exits: (?P<exits>\(?(north)?\)? ?\(?(east)?\)? ?\(?(south)?\)? ?\(?(west)?\)? ?\(?(up)?\)? ?\(?(down)?\)?) ?\]$`, M.moveDone)

	Config = cfg
	return Config
}
func NewMap() *Map {
	m := &Map{}
	m.areas = []*Area{}
	m.rooms = []*Room{}
	m.loaded = true
	m.mapping = true
	C.Print("\nMap created. Add an area to start mapping. Type #map new area <name>")
	return m
}

func (m *Map) NewArea(name string) *Area {
	a := &Area{}
	a.Name = name
	a.Rooms = make([]*Room, 1)
	m.AddArea(a)
	return a
}

func (m *Map) NewRoom(area *Area, name, exits string, c Coordinates) *Room {
	r := &Room{}
	r.id = m.getNewID()
	r.name = name
	r.area = area
	r.coordinates = c
	r.exits = getExits(exits)
	r.tags = make(map[string]string)
	r.exitString = exits
	area.AddRoom(r)
	m.AddRoom(r)
	return r
}

func (m *Map) AddRoom(room *Room) {
	m.rooms = append(m.rooms, room)
}

func (m *Map) AddArea(area *Area) {
	m.areas = append(m.areas, area)
}

func getExits(exits string) map[Direction]*Room {
	exitsMap := make(map[Direction]*Room)
	for _, e := range strings.Split(exits, " ") {
		if e == "" {
			continue
		}
		dir := dirmap[e]
		exitsMap[dir] = nil
	}
	return exitsMap
}

func (m *Map) SetRoom(r *Room) {
	m.room = r
}

// TODO implement mapper Load
func (m *Map) Load(path string) *Map { return nil }

// TODO implement mapper Save
func (m *Map) Save(path string) error { return nil }

// getNewID steps through all the rooms in the map and identifies the lowest
// available ID integer. Room IDs are reused after rooms are deleted.
func (m *Map) getNewID() int {
	min := math.MaxInt
	max := 0
	for _, r := range m.rooms {
		if r.id < min {
			min = r.id + 1
		}
		if r.id >= max {
			max = r.id + 1
		}
	}
	if min < max {
		return min
	}
	return max
}

// TODO implement FindRoom
func (m *Map) FindRoom(name string, exits string) *Room { return nil }

// TODO implement ShiftRoom
func (m *Map) ShiftRoom(room *Room, direction Direction) *Room { return nil }

func (m *Map) moveStart(move Direction) {
	m.nextMoves = append(m.nextMoves, move)
	moveDoneTrigger.SetCount(len(m.nextMoves))
}

// checkMove executes when a trigger indicating a move is complete fires
// this trigger needs to be set somewhere and correspond to a string
// such as a prompt that you see after all room information scrolls
// past
func (m *Map) moveDone(t *trigger.Trigger) {
	m.rmName = t.Results["name"]
	m.rmExitString = t.Results["exits"]
	if len(m.nextMoves) > 0 {
		// We should have all the information about the new room - we can check
		// if it matches what we expect to see. If it matches then we'll move
		// the map to the new room.
		nr := m.checkMove(m.nextMoves[0])
		if nr != nil {
			m.room = nr
			m.Show("map")
		} else {
			// move failed... do something!
		}
		m.nextMoves = m.nextMoves[1:]
	}
}

// moveFail should be linked to actions that print when a move fails, such as
// when a character is asleep, resting, runs into a wall, etc...
func (m *Map) moveFail(t *trigger.Trigger) {
	// pop a room from the queue, stop walking a path, ???

}

func (m *Map) checkMove(move Direction) *Room {
	if m.room == nil {
		return nil
	}

	if move == Look {
		if m.checkRoom(m.room) {
			return m.room
		} else {
			return nil
		}
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
	coordinates := m.getCoordinatesFromDir(m.room.coordinates, move)
	nr := m.NewRoom(m.room.area, m.rmName, m.rmExitString, coordinates)

	// Add exits to each room that link to each other
	for dir, rm := range nr.exits {
		if rm == nil {
			// We have a direction but no room, check if there is a room at
			// those coordinates
			coordinates := m.getCoordinatesFromDir(nr.coordinates, dir)
			pr := m.getRoomAtCoordinates(nr.area, coordinates)
			if len(pr) == 1 {
				if pr[0] != nil {
					// Only one room at those coordinates, link that mother.
					m.linkRooms(nr, pr[0], dir)

				}
			}
		}
	}
	return nr
}

func (m *Map) newRoom() *Room {
	r := &Room{}
	r.id = m.getNewID()
	r.name = m.rmName
	r.exitString = m.rmExitString
	r.coordinates = m.room.coordinates
	m.rooms[r.id] = r
	return r
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

// linkRooms connects two rooms bidirectionally by adding exits to each room
func (m *Map) linkRooms(from *Room, to *Room, move Direction) {
	// Add the exit to the from room
	from.exits[move] = to
	// Add the exit to the to room
	to.exits[reverse[move]] = from
}

// print takes the width and height of the desired map string and returns
// a map layout centered on your current room in one string. The map key
// can be found in map.md
func (m *Map) print(width, height int) []byte {
	if m.room == nil {
		// TODO better error message
		C.Print("\nmap error: I don't know where you are, can't print a map.")
	}

	nx := width / 3  // number of rooms we can fit into the map width
	ny := height / 3 // number of rooms we can fit into the map height
	cx := nx / 2     // X coordinates of the center room
	cy := ny / 2     // Y coordinates of the center room

	if nx <= 0 || ny <= 0 {
		return []byte{}
	}

	// initialize map string array
	ma := make([][]*Room, ny)
	cl := make([][]bool, ny)
	for row := 0; row < ny; row++ {
		ma[row] = make([]*Room, nx)
		cl[row] = make([]bool, nx)
	}

	// Starting at the top left populate each room string which will consist of
	// 9 characters indicating the top, middle, and bottom row of the 3x3 character
	// room string
	for row := 0; row < ny; row++ {
		for col := 0; col < nx; col++ {
			// Gets the room at the cooridnate offset from the current room and
			// on the same Z axis
			rs := m.getRoomAtCoordinates(m.room.area, Coordinates{
				(m.room.coordinates.X - cx) + col,
				(m.room.coordinates.Y + cy) - row,
				m.room.coordinates.Z,
			})
			if len(rs) <= 0 {
				// No room
				ma[row][col] = nil
			} else if len(rs) > 1 {
				// Multiple rooms at these coordinates, not sure how to display
				// a collision yet but we'll figure it out... maybe we'll print
				// the one with the shortest path to the current room or have
				// a collision indicator on that room?
				ma[row][col] = nil
				cl[row][col] = true
			} else {
				// Only one room at those coordinates
				ma[row][col] = rs[0]
				cl[row][col] = false
			}
		}
	}

	// Now we have a 2D array of rooms that will fit into the width/height
	// provided. We need to generate a string from them. I could do all of this
	// in the above loop, but I'd like to have the populated array for debugging
	s := []byte{}
	sa := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		sa[i] = make([]byte, 3)
	}

	// Go through each row, col and collect the top, middle, bottom strings
	// for each room into sa. At the end of each row, concatenate the
	// top, middle, and bottom strings with newlines and append them to s.
	for row := 0; row < ny; row++ {
		sa[0] = []byte("   ")
		sa[1] = []byte("   ")
		sa[2] = []byte("   ")
		for col := 0; col < nx; col++ {
			rs := ma[row][col].MapStrings()
			// We have a collision so we'll mark an asterisk on top right
			if cl[row][col] {
				sa[0][2] = '*'
			}
			if row == cy && col == cx {
				// At the center room
				rs[1][1] = '#'
			}
			for subrow := 0; subrow < 3; subrow++ {
				// Each room has 3 rows
				sa[subrow] = append(sa[subrow], rs[subrow]...)
			}
		}
		// Combine the three rows of each string into our main string
		for i := 0; i < 3; i++ {
			sa[i] = append(sa[i], '\n')
			s = append(s, sa[i]...)
		}
		//s += strings.Join(sa, "\n")
	}

	return s // You should at least have a bunch of spaces and newlines
}

func (m *Map) getRoomAtCoordinates(a *Area, c Coordinates) []*Room {
	rs := make([]*Room, 0)

	// Each area has a new set of coordinates
	for _, r := range a.Rooms {
		if r != nil {
			if r.coordinates.Equals(c) {
				rs = append(rs, r)
			}
		}
	}
	if len(rs) == 0 {
		rs = append(rs, nil)
	}
	return rs
}

// Each room is represented by a 3x3 array of characters indicating exits
func (r *Room) MapStrings() [][]byte {
	rs := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		rs[i] = make([]byte, 3)
		rs[i][0] = ' ' // Set all values to space by default
		rs[i][1] = ' '
		rs[i][2] = ' '
	}
	if r != nil {
		for dir := range r.exits {
			switch dir {
			case Up:
				rs[0][0] = '^'
			case North:
				rs[0][1] = '|'
			case West:
				rs[1][0] = '-'
			case East:
				rs[1][2] = '-'
			case Down:
				rs[2][0] = 'v'
			case South:
				rs[2][1] = '|'
			}
		}
	}

	return rs
}

func (c Coordinates) Equals(o Coordinates) bool {
	return c.X == o.X && c.Y == o.Y && c.Z == o.Z
}

func (a *Area) AddRoom(r *Room) {
	a.Rooms = append(a.Rooms, r)
}

func (m *Map) Show(name string) {
	w, h := C.GetWindowSize(name)
	C.PrintBytesTo(name, m.print(w, h))
}
