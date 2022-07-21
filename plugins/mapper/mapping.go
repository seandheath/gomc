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
	coordinates := m.GetCoordinatesFromDir(m.room.Coordinates, move)
	nr := m.NewRoom(m.room.area, m.rmName, m.rmExitString, coordinates)

	// Add exits to each room that link to each other if autolink is enabled
	if m.Autolink {
		for dir, r := range nr.exits {
			if r == nil {
				// We have a direction but no room, check if there is a room at
				// those coordinates
				coordinates := m.GetCoordinatesFromDir(nr.Coordinates, dir)
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
	from.ExitIDs[move] = to.ID
	// Add the exit to the to room
	to.exits[reverse[move]] = from
	to.ExitIDs[reverse[move]] = from.ID
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
