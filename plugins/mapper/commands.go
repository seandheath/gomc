package mapper

import (
	"github.com/seandheath/gomc/internal/client"
	"github.com/seandheath/gomc/pkg/trigger"
)

func addCommands(c *client.Client, m *Map) {
	c.AddAliasFunc("^#map new map$", m.NewMapCmd)
	c.AddAliasFunc("^#map new area (?P<name>.+)$", m.NewAreaCmd)
	c.AddAliasFunc("^#map na (?P<name>.+)$", m.NewAreaCmd)
	c.AddAliasFunc("^#map new room (?P<move>n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)
	c.AddAliasFunc("^#map nr (?P<move>n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)
	c.AddAliasFunc("^(?P<move>n|north|e|east|s|south|w|west|u|up|d|down|lo|loo|look|map|rec|reca|recal|recall)$", m.CaptureMoveCmd)
	c.AddAliasFunc("^(?P<moves>[neswud]+)$", m.CaptureMovesCmd)
	c.AddAliasFunc("^#map start$", m.StartCmd)
	c.AddAliasFunc("^#map stop$", m.StopCmd)
	C.AddAliasFunc("^#map show (?P<window>.+)$", m.ShowCmd)
}

func (m *Map) ShowCmd(t *trigger.Trigger) {
	m.Show(t.Results["window"])
}

func (m *Map) StartCmd(t *trigger.Trigger) {
	m.mapping = true
	C.Print("\nMapping started")
}
func (m *Map) StopCmd(t *trigger.Trigger) {
	m.mapping = false
	C.Print("\nMapping stopped")
}

func (m *Map) NewMapCmd(t *trigger.Trigger) {
	M = NewMap()
	C.Print("\nMap created - add an area with #map new area <name>")
}

func (m *Map) NewAreaCmd(t *trigger.Trigger) {
	name := t.Results["name"]
	if _, ok := m.GetArea(name); ok {
		C.Print("\nArea already exists: " + name)
		return
	}
	if m.room != nil {
		// I'm already in a room/area, we need to make a transition room

	}
	a := m.NewArea(name)
	r := m.NewRoom(a, m.rmName, m.rmExitString, Coordinates{0, 0, 0})
	m.room = r
	C.Print("\nArea created: " + name)
}

func (m *Map) GetArea(name string) (*Area, bool) {
	for _, a := range m.areas {
		if a.Name == name {
			return a, true
		}
	}
	return nil, false
}

func (m *Map) NewRoomCmd(t *trigger.Trigger) {
	if m.room == nil {
		C.Print("\nMAP: I don't know where you are. I can't add a room.")
	}
}

func (m *Map) CaptureMoveCmd(t *trigger.Trigger) {
	ds := t.Results["move"]
	if dir, ok := dirmap[ds]; ok {
		m.moveStart(dir)
	} else {
		// could be recall or look?
	}
	// Pass the move command to the MUD
	C.SendNow(t.Matches[1])
}

func (m *Map) CaptureMovesCmd(t *trigger.Trigger) {
	// TODO break the moves up and add them to the move queue
}

// As you move around the MUD rooms are added automatically.
// If you run into a wall this command will remove the room
// you just created.
func (m *Map) UndoCmd(t *trigger.Trigger) {

}

// Delete the room in the provided direction
//  #map delete room <direction>
func (m *Map) DeleteCmd(t *trigger.Trigger) {

}

// Moves the player on the map without moving them in the MUD
// #map move <direction>
func (m *Map) MoveCmd(t *trigger.Trigger) {

}
