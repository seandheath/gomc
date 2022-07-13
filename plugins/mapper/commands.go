package mapper

import (
	"github.com/seandheath/go-mud-client/internal/client"
	"github.com/seandheath/go-mud-client/pkg/trigger"
)

func addCommands(c *client.Client, m *Map) {
	C.AddAlias("#map new area (.+)$", m.NewAreaCmd)
	C.AddAlias("#map na (.+)$", m.NewAreaCmd)
	C.AddAlias("#map new room (n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)
	C.AddAlias("#map nr (n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)
	C.AddAlias("^(n|north|e|east|s|south|w|west|u|up|d|down|lo|loo|look|map)$", m.CaptureMoveCmd)
}

func (m *Map) NewAreaCmd(t *trigger.Match) {
	name := t.Matches[1]
	if _, ok := m.Areas[name]; ok {
		C.Print("Area already exists: " + name)
	} else {
		m.Areas[name] = &Area{Name: name, Rooms: make(map[int]*Room)}
	}
}

func (m *Map) NewRoomCmd(t *trigger.Match) {
	//name := t.Matches[1]
	//exits := t.Matches[2]
}

func (m *Map) CaptureMoveCmd(t *trigger.Match) {
	if dir, ok := dirmap[t.Matches[1]]; ok {
		m.move(dir)
	}
	// Pass the move command to the MUD
	C.SendNow(t.Matches[1])
}

// As you move around the MUD rooms are added automatically.
// If you run into a wall this command will remove the room
// you just created.
func (m *Map) UndoCmd(t *trigger.Match) {

}

// Delete the room in the provided direction
//  #map delete room <direction>
func (m *Map) DeleteCmd(t *trigger.Match) {

}

// Moves the player on the map without moving them in the MUD
// #map move <direction>
func (m *Map) MoveCmd(t *trigger.Match) {

}
