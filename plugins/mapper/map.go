package mapper

type Map struct {
	loaded bool
	recall *Room // Special room for recall TODO: make this configurable for multiple rooms/triggers
	room   *Room // Current room
	//area         *Area            // Current area not sure I'm gonna use this
	areas        []*Area     `yaml:"areas"` // All areas in the map
	rooms        []*Room     `yaml:"rooms"` // All rooms in the map
	nextMoves    []Direction // queued up moves
	pastMoves    []Direction // used to track last steps to determine location if we get lost
	rmExitString string      // rmExits are updated every time we see an exit string - even if we're not expecting a mov
	rmName       string      // rmName is updated every time we see a new room name - even if we're not expecting a move
	mapping      bool        // do not create or link rooms if we're not in mapping mode, puts map in read-only state
	debug        bool        // show debug information at the top of the map output
	autolink     bool        // automatically link rooms together when they are adjacent
}

func NewMap() *Map {
	m := &Map{}
	m.Reset()
	return m
}

func (m *Map) Reset() *Map {
	m.areas = []*Area{}
	m.rooms = []*Room{}
	m.loaded = true
	m.mapping = true
	m.debug = true
	m.autolink = true
	C.Print("\nMap created. Add an area to start mapping. Type #map new area <name>")
	return m
}

func (m *Map) DeleteRoom(r *Room) {
	r.area.RemoveRoom(r)
	for i, room := range m.rooms {
		if room == r {
			m.rooms = append(m.rooms[:i], m.rooms[i+1:]...)
			break
		}
	}
	r = nil
}

func (m *Map) AddArea(area *Area) {
	m.areas = append(m.areas, area)
}

func (m *Map) AddRoom(room *Room) {
	m.rooms = append(m.rooms, room)
}

func (m *Map) SetRoom(r *Room) {
	m.room = r
}

// GetNewID steps through all the rooms in the map and identifies the lowest
// available ID integer. Room IDs are reused after rooms are deleted.
func (m *Map) GetNewID() int {
	return getID(0, m.rooms)
}

func getID(id int, rooms []*Room) int {
	for _, r := range rooms {
		if r.id == id {
			return getID(id+1, rooms)
		}
	}
	return id
}

func (m *Map) GetRoomAtCoordinates(a *Area, c Coordinates) []*Room {
	rs := []*Room{}

	// Each area has a new set of coordinates
	for _, r := range a.Rooms {
		if r != nil {
			if r.coordinates.Equals(c) {
				rs = append(rs, r)
			}
		}
	}
	return rs
}

func (m *Map) GetCoordinatesFromDir(c Coordinates, move Direction) Coordinates {
	switch move {
	case North: // Positive Y axis
		return Coordinates{c.X, c.Y + 1, c.Z}
	case East: // Positive X axis
		return Coordinates{c.X + 1, c.Y, c.Z}
	case South: // Negative Y axis
		return Coordinates{c.X, c.Y - 1, c.Z}
	case West: // Negative X axis
		return Coordinates{c.X - 1, c.Y, c.Z}
	case Up: // Positive Z axis
		return Coordinates{c.X, c.Y, c.Z + 1}
	case Down: // Negative Z axis
		return Coordinates{c.X, c.Y, c.Z - 1}
	}
	// Shouldn't get here
	C.Print("Error creating coordinates")
	return Coordinates{}
}
