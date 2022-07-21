package mapper

type Area struct {
	Name  string        `yaml:"name"`
	Rooms map[int]*Room `yaml:"rooms"`
}

func (m *Map) NewArea(name string) *Area {
	a := &Area{}
	a.Name = name
	a.Rooms = map[int]*Room{}
	m.AddArea(a)
	return a
}

func (a *Area) AddRoom(r *Room) {
	a.Rooms[r.ID] = r
}

func (a *Area) RemoveRoom(r *Room) {
	a.Rooms[r.ID] = nil
}
