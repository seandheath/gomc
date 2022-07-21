package mapper

import "fmt"

type Coordinates struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

type Room struct {
	ID          int                 `yaml:"id"`   // Unique ID
	Name        string              `yaml:"name"` // Room name as seen in the game
	exits       map[Direction]*Room // Need to be reconstituted on load
	ExitIDs     map[Direction]int   `yaml:"exits"`
	ExitString  string              `yaml:"exitString"`
	area        *Area               // The parent area of this room
	Coordinates Coordinates         `yaml:"coordinates"` // The coordinates of this room
	Tags        map[string]string   `yaml:"tags"`        // List of key:value flags for this room
	//extraIdentifiers map[string]string `yaml:"identifier"`  // A list of key:value identifiers for this room
}

func (m *Map) NewRoom(a *Area, name, exits string, c Coordinates) *Room {
	r := &Room{}
	r.ID = m.GetNewID()
	r.Name = name
	r.area = a
	r.Coordinates = c
	r.exits = GetExits(exits)
	r.ExitIDs = GetExitIDs(r.exits)
	r.Tags = map[string]string{}
	r.ExitString = exits
	m.AddRoom(r)
	a.AddRoom(r)
	return r
}

func (c Coordinates) Equals(o Coordinates) bool {
	return c.X == o.X && c.Y == o.Y && c.Z == o.Z
}

func (c Coordinates) String() string {
	return fmt.Sprintf("%d,%d,%d", c.X, c.Y, c.Z)
}
