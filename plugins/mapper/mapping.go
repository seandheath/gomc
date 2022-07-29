package mapper

import (
	"strings"

	"github.com/seandheath/gomc/pkg/trigger"
)

// AddRoomFromMove takes a direction
func (m *Map) AddRoomFromMove(move Direction) *Room {
	if m.Room == nil {
		// I don't know where I am so I can't link a room
		// TODO print some message here
		C.Print("\nCan't create a room, I don't know where I am. Set your current room with '#map goto <roomID>'")
		return nil
	}
	// We have a room that we are coming from, let's add a room to it
	coordinates := m.GetCoordinatesFromDir(m.Room.Coordinates, move)
	nr := m.NewRoom(m.Room.Area, m.rmName, m.rmExitString, coordinates)

	// Add exits to each room that link to each other if autolink is enabled
	if m.Autolink {
		for dir, r := range nr.Exits {
			if r == nil {
				// We have a direction but no room, check if there is a room at
				// those coordinates
				coordinates := m.GetCoordinatesFromDir(nr.Coordinates, dir)
				pr := m.GetRoomAtCoordinates(nr.Area, coordinates)
				if len(pr) == 1 {
					if pr[0] != nil {
						// Only one room at those coordinates, link that mother.
						m.linkRooms(nr, pr[0], dir)
					}
				}
			}
		}
	} else {

		// Only link the previous move
		m.linkRooms(m.Room, nr, move)
	}
	return nr
}

// linkRooms connects two rooms bidirectionally by adding exits to each room
func (m *Map) linkRooms(from *Room, to *Room, move Direction) {
	// Add the exit to the from room
	from.Exits[move] = to
	// Add the exit to the to room
	to.Exits[reverse[move]] = from

	// If the room we're coming from has a door at that exit then we need to
	// add the door to this room's door array as well
	if door, ok := to.Doors[reverse[move]]; ok {
		from.Doors[move] = door
	}
}

func GetExits(exits string) map[Direction]*Room {
	exitsMap := map[Direction]*Room{}
	for _, e := range strings.Split(exits, " ") {
		if e == "" {
			continue
		}
		dir := dirmap[e]
		exitsMap[dir] = nil

	}
	return exitsMap
}

func GetExitIDs(exits map[Direction]*Room) map[Direction]int {
	idMap := map[Direction]int{}
	for dir, r := range exits {
		if r == nil {
			idMap[dir] = 0
		} else {
			idMap[dir] = r.ID
		}
	}
	return idMap
}

// TODO implement ShiftRoom
func (m *Map) ShiftRoom(room *Room, direction Direction) *Room { return nil }

func (m *Map) MapDoor(t *trigger.Trigger) {
	// only auto-add doors if we're mapping
	if m.Mapping && m.Room != nil {
		// make sure we have a valid direction
		if dir, ok := dirmap[t.Results["dir"]]; ok {
			open := t.Results["open"]
			name := t.Results["door"]
			locked := strings.Contains(open, "lock")
			addDoor(m.Room, dir, name, locked)
			if r, ok := m.Room.Exits[dir]; ok {
				if r != nil {
					addDoor(r, reverse[dir], name, locked)
				}
			}
		}
		m.Show("map")
	}
}

func addDoor(room *Room, dir Direction, name string, locked bool) {

	if d, ok := room.Doors[dir]; ok {
		// If we have a door but it's not locked and we saw a lock string
		// we need to add a lock to the door
		if !d.Locked && locked {
			d.Locked = true
		}
	} else {
		room.Doors[dir] = &Door{
			Name:   name,
			Locked: locked,
		}
	}
}

func (m *Map) Unlink(room *Room, dir Direction, both bool) {
	if room != nil {
		if r, ok := room.Exits[dir]; ok {
			if r != nil && both {
				// There is a room there now, we'll unlink us from them
				// TODO might want to have a one-way here?
				r.Exits[reverse[dir]] = nil
			}
			room.Exits[dir] = nil
		}
	}
}
