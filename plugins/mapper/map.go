package mapper

type Map struct {
	Recall       int              `yaml:"recall"` // Special room for recall TODO: make this configurable for multiple rooms/triggers
	room         *Room            // Current room
	Areas        map[string]*Area `yaml:"areas"` // All areas in the map, name is key
	rooms        map[int]*Room    // All rooms in the map, ID is key
	nextMoves    []Direction      // queued up moves
	pastMoves    []Direction      // used to track last steps to determine location if we get lost
	rmExitString string           // rmExits are updated every time we see an exit string - even if we're not expecting a mov
	rmName       string           // rmName is updated every time we see a new room name - even if we're not expecting a move
	Mapping      bool             `yaml:"mapping"`  // do not create or link rooms if we're not in mapping mode, puts map in read-only state
	Debug        bool             `yaml:"debug"`    // show debug information at the top of the map output
	Autolink     bool             `yaml:"autolink"` // automatically link rooms together when they are adjacent
}

func NewMap() *Map {
	m := &Map{}
	m.Reset()
	return m
}

// Rebuild is called when a map is loaded from disk. It steps through all
// the rooms and populates the exit arrays with pointers and collects all the
// room objects into the map.rooms object
func (m *Map) Rebuild() {
	// step through each area
	for _, a := range m.Areas {
		// get all the room objects from each area
		for _, r := range a.Rooms {
			if r != nil {
				r.area = a
				// add the room to the map
				m.rooms[r.ID] = r
			}
		}
	}
	// Go through all the rooms and populate the exit arrays with pointers to
	// the room objects
	for _, cr := range m.rooms {
		cr.exits = map[Direction]*Room{}
		for dir, id := range cr.ExitIDs {
			if nr, ok := m.rooms[id]; ok {
				cr.exits[dir] = nr
			} else {
				if id == 0 {
					cr.exits[dir] = nil
				}
			}
		}
	}
}

func (m *Map) PrepareSave() {
	for _, r := range m.rooms {
		if r != nil {
			r.ExitIDs = map[Direction]int{}
			for dir, nr := range r.exits {
				if nr != nil {
					r.ExitIDs[dir] = nr.ID
				} else {
					r.ExitIDs[dir] = 0
				}
			}
		}
	}
}

func (m *Map) Reset() {
	m.Recall = 0
	m.room = nil
	m.Areas = map[string]*Area{}
	m.rooms = map[int]*Room{}
	m.nextMoves = []Direction{}
	m.pastMoves = []Direction{}
	m.rmExitString = ""
	m.rmName = ""
	m.Mapping = true
	m.Debug = true
	m.Autolink = true
	C.Print("\nMap created. Add an area to start mapping. Type #map new area <name>")
}

// Returns the room with the given ID or nil if the room doesn't exist
func (m *Map) GetRoom(id int) *Room {
	if r, ok := m.rooms[id]; ok {
		if r != nil {
			return r
		}
	}
	return nil
}

func (m *Map) DeleteRoom(r *Room) {
	m.rooms[r.ID] = nil
}

func (m *Map) AddArea(area *Area) {
	m.Areas[area.Name] = area
}

func (m *Map) AddRoom(room *Room) {
	m.rooms[room.ID] = room
}

func (m *Map) SetRoom(r *Room) {
	m.room = r
}

// GetNewID steps through all the rooms in the map and identifies the lowest
// available ID integer. Room IDs are reused after rooms are deleted.
func (m *Map) GetNewID() int {
	return getID(1, m.rooms)
}

func getID(id int, rooms map[int]*Room) int {
	for _, r := range rooms {
		if r != nil {
			if r.ID == id {
				return getID(id+1, rooms)
			}
		}
	}
	return id
}

func (m *Map) GetRoomAtCoordinates(a *Area, c Coordinates) []*Room {
	rs := []*Room{}

	// Each area has a new set of coordinates
	for _, r := range a.Rooms {
		if r != nil {
			if r.Coordinates.Equals(c) {
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

func (m *Map) Unlink(room *Room, dir Direction, both bool) {
	if room != nil {
		if r, ok := room.exits[dir]; ok {
			if r != nil && both {
				// There is a room there now, we'll unlink us from them
				// TODO might want to have a one-way here?
				r.exits[reverse[dir]] = nil
			}
			room.exits[dir] = nil
		}
	}
}
