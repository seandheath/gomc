package mapper

import (
	"fmt"
	"strings"
)

type Coordinates struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

type Door struct {
	Name   string `yaml:"name"`
	Locked bool   `yaml:"locked"`
}

type Room struct {
	ID          int                 `yaml:"id"`   // Unique ID
	Name        string              `yaml:"name"` // Room name as seen in the game
	Exits       map[Direction]*Room // Need to be reconstituted on load
	Doors       map[Direction]*Door `yaml:"doors"`
	ExitIDs     map[Direction]int   `yaml:"exitIDs"`
	ExitString  string              `yaml:"exitString"`
	Area        *Area               // The parent area of this room
	Coordinates Coordinates         `yaml:"coordinates"` // The coordinates of this room
	Tags        []string            `yaml:"tags"`        // List of key:value flags for this room
	//extraIdentifiers map[string]string `yaml:"identifier"`  // A list of key:value identifiers for this room
}

func (m *Map) NewRoom(a *Area, name, exits string, c Coordinates) *Room {
	r := &Room{}
	r.ID = m.GetNewID()
	r.Name = name
	r.Area = a
	r.Coordinates = c
	r.Exits = GetExits(exits)
	r.Doors = map[Direction]*Door{}
	r.Tags = []string{}
	r.ExitString = strings.TrimSpace(exits)
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
