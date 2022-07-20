package mapper

type Area struct {
	Name  string  `yaml:"name"`
	Rooms []*Room `yaml:"rooms"`
}

func (m *Map) NewArea(name string) *Area {
	a := &Area{}
	a.Name = name
	a.Rooms = []*Room{}
	m.AddArea(a)
	return a
}

func (a *Area) AddRoom(r *Room) {
	a.Rooms = append(a.Rooms, r)
}

func (a *Area) RemoveRoom(r *Room) {
	for i, room := range a.Rooms {
		if room == r {
			a.Rooms = append(a.Rooms[:i], a.Rooms[i+1:]...)
			break
		}
	}
}
