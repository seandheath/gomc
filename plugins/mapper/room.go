package mapper

import "fmt"

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

func (m *Map) NewRoom(area *Area, name, exits string, c Coordinates) *Room {
	r := &Room{}
	r.id = m.GetNewID()
	r.name = name
	r.area = area
	r.coordinates = c
	r.exits = GetExits(exits)
	r.tags = map[string]string{}
	r.exitString = exits
	area.AddRoom(r)
	m.AddRoom(r)
	return r
}

func (c Coordinates) Equals(o Coordinates) bool {
	return c.X == o.X && c.Y == o.Y && c.Z == o.Z
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%d,%d,%d", c.X, c.Y, c.Z)
}
