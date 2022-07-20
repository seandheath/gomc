package mapper

import "strings"

// AddRoomFromMove takes a direction
func (m *Map) AddRoomFromMove(move Direction) *Room {
	if m.room == nil {
		// I don't know where I am so I can't link a room
		// TODO print some message here
		C.Print("\nCan't create a room, I don't know where I am. Set your current room with '#map goto <roomID>'")
		return nil
	}
	// We have a room that we are coming from, let's add a room to it
	coordinates := m.GetCoordinatesFromDir(m.room.coordinates, move)
	nr := m.NewRoom(m.room.area, m.rmName, m.rmExitString, coordinates)

	// Add exits to each room that link to each other if autolink is enabled
	if m.autolink {
		for dir, rm := range nr.exits {
			if rm == nil {
				// We have a direction but no room, check if there is a room at
				// those coordinates
				coordinates := m.GetCoordinatesFromDir(nr.coordinates, dir)
				pr := m.GetRoomAtCoordinates(nr.area, coordinates)
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
		m.linkRooms(m.room, nr, move)
	}
	return nr
}

// linkRooms connects two rooms bidirectionally by adding exits to each room
func (m *Map) linkRooms(from *Room, to *Room, move Direction) {
	// Add the exit to the from room
	from.exits[move] = to
	// Add the exit to the to room
	to.exits[reverse[move]] = from
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

// TODO implement ShiftRoom
func (m *Map) ShiftRoom(room *Room, direction Direction) *Room { return nil }
