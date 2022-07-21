package mapper

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/seandheath/gomc/pkg/trigger"
)

func addCommands(m *Map) {
	// Map commands
	C.AddAliasFunc("^#map reset$", m.ResetCmd)
	C.AddAliasFunc("^#map new map$", m.NewMapCmd)
	C.AddAliasFunc("^#map show (?P<window>.+)$", m.ShowCmd)
	C.AddAliasFunc("^#map save (?P<path>.+)$", m.SaveCmd)
	C.AddAliasFunc("^#map load (?P<path>.+)$", m.LoadCmd)

	// Area commands
	C.AddAliasFunc("^#map new area (?P<name>.+)$", m.NewAreaCmd)
	C.AddAliasFunc("^#map na (?P<name>.+)$", m.NewAreaCmd)

	// Room commands
	C.AddAliasFunc("^#map new room (?P<move>n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)
	C.AddAliasFunc("^#map nr (?P<move>n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)

	// Mapping commands
	C.AddAliasFunc("^#map undo$", m.UndoCmd)
	C.AddAliasFunc("^#map start$", m.StartCmd)
	C.AddAliasFunc("^#map stop$", m.StopCmd)
	C.AddAliasFunc("^#map autolink (?P<set>on|off)$", m.AutoLinkCmd)
	C.AddAliasFunc(`^#map link (?P<dir>(n|e|s|w|u|d)) (?P<id>\d+)$`, m.LinkDirCmd)
	C.AddAliasFunc(`^#map rmlink (?P<dir>(n|e|s|w|u|d))$`, m.UnlinkDirCmd) // bi-directional
	C.AddAliasFunc(`^#map rmexit (?P<dir>(n|e|s|w|u|d))$`, m.RmExitCmd)    // single

	// Move commands
	C.AddAliasFunc("^(?P<move>north|east|south|west|up|down|lo|loo|look|map|rec|reca|recal|recall)$", m.CaptureMoveCmd)
	C.AddAliasFunc(`^(?P<speedwalk>speedwalk)? ?(?P<steps>(\d*(n|e|s|w|u|d))+)$`, m.CaptureMovesCmd)
}

func (m *Map) RmExitCmd(t *trigger.Trigger) {
	if dir, ok := dirmap[t.Results["dir"]]; ok {
		m.Unlink(m.room, dir, false)
	}
}

func (m *Map) UnlinkDirCmd(t *trigger.Trigger) {
	if dir, ok := dirmap[t.Results["dir"]]; ok {
		m.Unlink(m.room, dir, true)
	}
}

func (m *Map) LinkDirCmd(t *trigger.Trigger) {
	if dir, ok := dirmap[t.Results["dir"]]; ok {
		if m.room != nil {
			id, err := strconv.Atoi(t.Results["id"])
			if err != nil {
				C.Print("\nMAP: Failed to parse id: " + t.Results["id"])
				return
			}
			nr := m.GetRoom(id)
			if nr == nil {
				C.Print(fmt.Sprintf("\nMAP: Unable to find room with ID: %d", id))
				return
			}
			m.room.exits[dir] = nr
			nr.exits[reverse[dir]] = m.room
			C.Print(fmt.Sprintf("\nLinked %d to the %s", id, dir))
		}
	}
	m.Show("map")
}

func (m *Map) LoadCmd(t *trigger.Trigger) {
	path := t.Results["path"]
	m.Load(path)
}

func (m *Map) SaveCmd(t *trigger.Trigger) {
	path := t.Results["path"]
	SaveMap(m, path)
}

func (m *Map) AutoLinkCmd(t *trigger.Trigger) {
	if t.Results["set"] == "on" {
		C.Print("\nMAP: Auto-link on")
		m.Autolink = true
	} else {
		C.Print("\nMAP: AutoLink off")
		m.Autolink = false
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
	m.Mapping = true
	C.Print("\nMapping started")
}
func (m *Map) StopCmd(t *trigger.Trigger) {
	m.Mapping = false
	C.Print("\nMapping stopped")
}

func (m *Map) NewMapCmd(t *trigger.Trigger) {
	M = NewMap()
	C.Print("\nMap created - add an area with #map new area <name>")
}

func (m *Map) NewAreaCmd(t *trigger.Trigger) {
	name := t.Results["name"]
	if _, ok := m.Areas[name]; ok {
		C.Print("\nArea already exists: " + name)
		return
	}
	a := m.NewArea(name)

	// Move our current room to the new area
	if m.room != nil {
		m.room.area.RemoveRoom(m.room)
		m.room.area = a
		a.AddRoom(m.room)
		m.room.Coordinates = Coordinates{0, 0, 0} // Every area has it's own coordinate set
	} else {
		// Room is nil so we don't have any rooms yet, make one
		r := m.NewRoom(a, m.rmName, m.rmExitString, Coordinates{0, 0, 0})
		m.room = r
	}

	C.Print("\nArea created: " + name)
	m.Show("map")
}

func (m *Map) NewRoomCmd(t *trigger.Trigger) {
	if m.room == nil {
		C.Print("\nMAP: I don't know where you are. I can't add a room.")
		return
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
