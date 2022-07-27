package mapper

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/seandheath/gomc/pkg/trigger"
)

func addCommands(m *Map) {
	// Map commands
	C.AddAlias("^#map reset$", m.ResetCmd)
	C.AddAlias("^#map new map$", m.NewMapCmd)
	C.AddAlias("^#map show (?P<window>.+)$", m.ShowCmd)
	C.AddAlias("^#map save ?(?P<path>.+)?$", m.SaveCmd)
	C.AddAlias("^#map load ?(?P<path>.+)?$", m.LoadCmd)

	// Area commands
	C.AddAlias("^#map new area (?P<name>.+)$", m.NewAreaCmd)
	C.AddAlias("^#map na (?P<name>.+)$", m.NewAreaCmd)

	// Room commands
	C.AddAlias("^#map new room (?P<move>n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)
	C.AddAlias("^#map nr (?P<move>n|north|e|east|s|south|w|west|u|up|d|down)$", m.NewRoomCmd)
	C.AddAlias(`^#map shift (?P<dir>(n|e|s|w|u|d))$`, m.ShiftCmd)

	// Mapping commands
	C.AddAlias("^#map undo$", m.UndoCmd)
	C.AddAlias("^#map start$", m.StartCmd)
	C.AddAlias("^#map stop$", m.StopCmd)
	C.AddAlias("^#map autolink (?P<set>on|off)$", m.AutoLinkCmd)
	C.AddAlias(`^#map link (?P<dir>(n|e|s|w|u|d)) (?P<id>\d+)$`, m.LinkDirCmd)
	C.AddAlias(`^#map linkone (?P<dir>(n|e|s|w|u|d)) (?P<id>\d+)$`, m.LinkOneDirCmd)
	C.AddAlias(`^#map rmlink (?P<dir>(n|e|s|w|u|d))$`, m.UnlinkDirCmd)      // uni-directional
	C.AddAlias(`^#map rmlinks (?P<dir>(n|e|s|w|u|d))$`, m.UnlinkDirBothCmd) // bi-directional
	C.AddAlias(`^#map rmexit (?P<dir>(n|e|s|w|u|d))$`, m.RmExitCmd)
	C.AddAlias(`^#map rmroom (?P<id>\d+)$`, m.RmRoomCmd)               // bi-directional
	C.AddAlias(`^#map rmroom (?P<dir>(n|e|s|w|u|d))$`, m.RmRoomDirCmd) // bi-directional
	C.AddAlias(`^#map tag add (?P<tag>.+)$`, m.TagAddCmd)
	C.AddAlias(`^#map tag delete (?P<tag>.+)$`, m.TagDeleteCmd)
	C.AddAlias(`^#map tag show$`, m.TagShowCmd)
	C.AddAlias(`^#map clean$`, m.CleanCmd)

	// Move commands
	C.AddAlias("^(?P<move>north|east|south|west|up|down|lo|loo|look|map|rec|reca|recal|recall)$", m.CaptureMoveCmd)
	C.AddAlias(`^(?P<speedwalk>speedwalk)? ?(?P<steps>(\d*(n|e|s|w|u|d))+)$`, m.CaptureMovesCmd)
	C.AddAlias(`^#map go (?P<id>\d+)$`, m.GoIDCmd)
	C.AddAlias(`^#go (?P<id>\d+)$`, m.GoIDCmd)
	C.AddAlias(`^#map go (?P<tag>\w+)$`, m.GoTagCmd)
	C.AddAlias(`^#go (?P<tag>\w+)$`, m.GoTagCmd)
}

func (m *Map) GoIDCmd(t *trigger.Trigger) {
	if m.room == nil {
		C.Print("\nMAP: I don't know where you are, so I can't go there.")
		return
	}
	start := m.room
	id, err := strconv.Atoi(t.Results["id"])
	if err != nil {
		C.Print("\nMAP: Invalid ID: " + t.Results["id"])
		return
	}
	finish := m.GetRoom(id)
	path, len := GetPath(start, finish)
	C.Parse(string(path))
	C.Print(fmt.Sprintf("\nPath found, length: %d, steps: %s\n\n", len, string(path)))
}
func (m *Map) GoTagCmd(t *trigger.Trigger) {
	if m.room == nil {
		C.Print("\nMAP: I don't know where you are, so I can't go there.")
		return
	}
	rms := m.GetRoomsByTag(t.Results["tag"])
	if rms == nil {
		C.Print("\nMAP: No rooms found with tag: " + t.Results["tag"])
		return
	}
	path := []byte{}
	length := math.MaxInt

	for _, rm := range rms {
		p, l := GetPath(m.room, rm)
		if l < length {
			path = p
			length = l
		}
	}
	if len(path) > 0 {
		C.Parse(string(path))
		C.Print(fmt.Sprintf("\nPath found, length: %d, steps: %s\n\n", length, string(path)))
	} else {
		C.Print("\nMAP: No path found.")
	}
}

func (m *Map) TagShowCmd(t *trigger.Trigger) {
	if m.room != nil {
		C.Print("\nTags: " + strings.Join(m.room.Tags, ", "))
	}
}
func (m *Map) TagAddCmd(t *trigger.Trigger) {
	if m.room != nil {
		m.room.Tags = append(m.room.Tags, t.Results["tag"])
		C.Print("\nTag added: " + t.Results["tag"])
	}
}

func (m *Map) TagDeleteCmd(t *trigger.Trigger) {
	if m.room != nil {
		for i, tag := range m.room.Tags {
			if tag == t.Results["tag"] {
				m.room.Tags = append(m.room.Tags[:i], m.room.Tags[i+1:]...)
				C.Print("\nTag deleted: " + t.Results["tag"])
				return
			}
		}
		C.Print("\nTag not found: " + t.Results["tag"])
	}
}

func (m *Map) RmExitCmd(t *trigger.Trigger) {
	if dir, ok := dirmap[t.Results["dir"]]; ok {
		delete(m.room.exits, dir)
	}
	m.Show("map")
}

func (m *Map) RmRoomDirCmd(t *trigger.Trigger) {
	r := GetRoomFromDir(m.room, t.Results["dir"])
	m.RmRoom(r)

}

func GetRoomFromDir(rm *Room, dir string) *Room {
	if dir, ok := dirmap[dir]; ok {
		return rm.exits[dir]
	}
	return nil
}

func (m *Map) RmRoomCmd(t *trigger.Trigger) {
	if m.room == nil {
		return
	}
	id, err := strconv.Atoi(t.Results["id"])
	if err != nil {
		C.Print("\nMAP: Invalid ID: " + t.Results["id"])
		return
	}
	tr := m.GetRoom(id)
	m.RmRoom(tr)
}

func (m *Map) RmRoom(rm *Room) {
	if rm != nil {
		for _, r0 := range m.rooms {
			for dir, r1 := range r0.exits {
				if rm == r1 {
					// Remove links to the target room
					m.Unlink(r0, dir, false)
				}
			}
		}
		delete(m.rooms, rm.ID)
		delete(rm.area.Rooms, rm.ID)
	}
}

func (m *Map) UnlinkDirCmd(t *trigger.Trigger) {
	if dir, ok := dirmap[t.Results["dir"]]; ok {
		m.Unlink(m.room, dir, false)
	}
	m.Show("map")
}

func (m *Map) UnlinkDirBothCmd(t *trigger.Trigger) {
	if dir, ok := dirmap[t.Results["dir"]]; ok {
		m.Unlink(m.room, dir, true)
	}
	m.Show("map")
}

func (m *Map) LinkOneDirCmd(t *trigger.Trigger) {
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
			C.Print(fmt.Sprintf("\nLinked %d to the %s", id, dir))
		}
	}
	m.Show("map")
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
	m.Show("map")
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

func (m *Map) ShiftCmd(t *trigger.Trigger) {
	if m.room == nil {
		C.Print("\nMAP: I don't know where you are. I can't shift.")
		return
	}
	if dir, ok := dirmap[t.Results["dir"]]; ok {
		nc := m.GetCoordinatesFromDir(m.room.Coordinates, dir)
		m.room.Coordinates = nc
		C.Print("\nShifted room " + string(dir) + ".\n")
	}
}

func (m *Map) CleanCmd(t *trigger.Trigger) {
	for _, r := range m.rooms {
		if r != nil {
			r.ExitString = strings.TrimSpace(r.ExitString)
		}
	}
}
