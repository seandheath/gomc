package mapper

import (
	"regexp"
	"strconv"

	"github.com/seandheath/gomc/internal/client"
	"github.com/seandheath/gomc/pkg/trigger"
)

func addCommands(c *client.Client, m *Map) {
	c.AddAliasFunc("^#map reset$", m.ResetCmd)
	c.AddAliasFunc("^#map new map$", m.NewMapCmd)
	c.AddAliasFunc("^#map new area (?P<name>.+)$", m.NewAreaCmd)
	c.AddAliasFunc("^#map na (?P<name>.+)$", m.NewAreaCmd)
	c.AddAliasFunc("^#map new room (?P<move>n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)
	c.AddAliasFunc("^#map nr (?P<move>n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)
	c.AddAliasFunc("^(?P<move>north|east|south|west|up|down|lo|loo|look|map|rec|reca|recal|recall)$", m.CaptureMoveCmd)
	c.AddAliasFunc(`^(?P<speedwalk>speedwalk)? ?(?P<steps>(\d*(n|e|s|w|u|d))+)$`, m.CaptureMovesCmd)
	c.AddAliasFunc("^#map start$", m.StartCmd)
	c.AddAliasFunc("^#map stop$", m.StopCmd)
	C.AddAliasFunc("^#map show (?P<window>.+)$", m.ShowCmd)
	C.AddAliasFunc("^#map undo$", m.UndoCmd)
	C.AddAliasFunc("^#map autolink (?P<set>on|off)$", m.AutoLinkCmd)

}

func (m *Map) AutoLinkCmd(t *trigger.Trigger) {
	if t.Results["set"] == "on" {
		C.Print("\nMAP: Auto-link on")
		m.autolink = true
	} else {
		C.Print("\nMAP: AutoLink off")
		m.autolink = false
	}
}

func (m *Map) ResetCmd(t *trigger.Trigger) {
	m.Reset()
}

func (m *Map) UndoCmd(t *trigger.Trigger) {
	r := m.rooms[len(m.rooms)-1]
	m.DeleteRoom(r)
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
		m.nextMoves = append(m.nextMoves, dir)
	}
	// Pass the move command to the MUD
	C.SendNow(t.Matches[1])
}

var stepRegex = regexp.MustCompile(`(?P<num>\d*)(?P<dir>n|e|s|w|u|d)`)

func (m *Map) CaptureMovesCmd(t *trigger.Trigger) {
	// TODO break the moves up and add them to the move queue
	g := stepRegex.FindAllStringSubmatch(t.Results["steps"], -1)
	for _, step := range g {
		num := step[1]
		d := step[2]
		if num == "" {
			num = "1"
		}
		n, err := strconv.Atoi(num)
		if err != nil {
			C.Print("\nMAP: Invalid number while parsing string: " + num)
			return
		}
		for i := 0; i < n; i++ {
			if dir, ok := dirmap[d]; ok {
				m.nextMoves = append(m.nextMoves, dir)
			}
		}
	}
	if len(g) == 1 && len(g[0][0]) == 1 {
		// only have one move
		C.SendNow(t.Results["steps"])
	} else {
		// prepend speedwalk to it
		C.SendNow("speedwalk " + t.Results["steps"])
	}
}

// Delete the room in the provided direction
//  #map delete room <direction>
func (m *Map) DeleteCmd(t *trigger.Trigger) {

}

// Moves the player on the map without moving them in the MUD
// #map move <direction>
func (m *Map) MoveCmd(t *trigger.Trigger) {

}
