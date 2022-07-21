package mapper

import (
	"regexp"

	"github.com/seandheath/gomc/pkg/trigger"
)

var stripParen = regexp.MustCompile(`\(?\)?`)

// checkMove executes when a trigger indicating a move is complete fires
// this trigger needs to be set somewhere and correspond to a string
// such as a prompt that you see after all room information scrolls
// past
func (m *Map) MoveDone(t *trigger.Trigger) {
	m.rmName = t.Results["name"]
	m.rmExitString = stripParen.ReplaceAllString(t.Results["exits"], "")
	if len(m.nextMoves) > 0 {
		// We should have all the information about the new room - we can check
		// if it matches what we expect to see. If it matches then we'll move
		// the map to the new room.
		nr := m.checkMove(m.nextMoves[0])
		if nr != nil {
			m.room = nr
		} else {
			// Couldn't find the room... uh oh
			pr := m.FindRoom(m.rmName, m.rmExitString)
			if len(pr) <= 0 {
				// No room matches this description, we're lost.
				m.room = nil
			} else if len(pr) == 1 {
				// There is only one room that matches this description, so
				// we'll just set ourselves to it
				m.room = pr[0]
			} else {
				// There are multiple rooms that match this description, so
				// we'll need to ask the user which room they want to move to.
				// For now we'll just choose the first one
				m.room = pr[0]
			}
		}
		m.nextMoves = m.nextMoves[1:]
		m.Show("map")
	}
}

func (m *Map) FindRoom(name string, exits string) []*Room {
	rl := []*Room{}
	for _, r := range m.rooms {
		if r.Name == name && r.ExitString == exits {
			rl = append(rl, r)
		}
	}
	return rl
}

func (m *Map) checkMove(move Direction) *Room {
	if m.room == nil {
		return nil
	}

	if move == Look {
		if m.checkRoom(m.room) {
			return m.room
		} else {
			return nil
		}
	}

	if move == Recall {
		if m.rooms[m.Recall] != nil {
			return m.rooms[m.Recall]
		}
	}

	// This room doesn't have an exit in that direction...
	if r, ok := m.room.exits[move]; !ok {
		return nil
	} else {
		if r != nil {
			// We've already got a room at that exit, return it
			return r
		} else {
			// We've got an exit but don't know what room is there
			// possible room based on coordinates
			c := m.GetCoordinatesFromDir(m.room.Coordinates, move)
			prc := m.GetRoomAtCoordinates(m.room.area, c)
			if len(prc) == 1 {
				// There is one room at the coordinates specified, let's check if we're in it
				if m.checkRoom(prc[0]) {
					// We're in the room at those coords, add the link
					m.linkRooms(m.room, prc[0], move)
					return prc[0]
				}
			} else if len(prc) > 1 {
				// More than one room at those coordinates, we'll need to check which one we're in
				// returns -1 if no matches
				r := m.checkRooms(prc)
				if r != nil {
					// We got a match
					m.linkRooms(m.room, r, move)
					return r
				}
			}
			// None of the possible rooms matched
			if m.Mapping {
				// Add a new room
				return m.AddRoomFromMove(move)
			}
		}
	}
	// Couldn't find a room or make one, we're lost
	return nil
}

// Could possibly fail if there are collisions of rooms with the same name and
// exits, but we'll cross that bridge when we come to it
func (m *Map) checkRooms(l []*Room) *Room {
	for _, r := range l {
		if m.checkRoom(r) {
			return r
		}
	}
	return nil
}

// checkRoom compares the provided room to the room most recently seen
// by the map. If they match, we'll return true.
func (m *Map) checkRoom(r *Room) bool {
	// Currently only checking name and exit string
	if (r.Name == m.rmName) && (r.ExitString == m.rmExitString) {
		return true
	}
	return false
}

// MoveFail should be linked to actions that print when a move fails, such as
// when a character is asleep, resting, runs into a wall, etc...
func (m *Map) MoveFail(t *trigger.Trigger) {
	m.nextMoves = m.nextMoves[1:]
}

// MoveRecall moves the user to the `recall` room saved in the map. The `recall`
// room can be set using the `#map set recall <id>` command.
func (m *Map) MoveRecall(t *trigger.Trigger) {
	if _, ok := m.rooms[m.Recall]; ok {
		m.room = m.rooms[m.Recall]
	}
}

// Clears the pending moves
func (m *Map) MoveClear(t *trigger.Trigger) {
	m.nextMoves = []Direction{}
}
