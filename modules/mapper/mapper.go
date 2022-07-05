package mapper

type Direction int

const (
	North Direction = iota
	East
	South
	West
	Up
	Down
)

type Coordinates struct {
	X int
	Y int
	Z int
}

type Exit struct {
	Direction Direction
	Room      *Room
	Door      string
	Locked    bool
}

type Room struct {
	ID          int // Unique ID
	Name        string
	Exits       []Exit
	Area        *Area
	Coordinates Coordinates
	Flags       []string
}

type Area struct {
	Name  string // Must have unique name
	Rooms map[int]Room
}

type Map struct {
	CurrentRoom *Room
	CurrentArea *Area
	Areas       map[string]*Area
	Rooms       map[string]*Room
}

func (m *Map) Load() {}

func Move(direction Direction) *Room                  { return nil }
func AddRoom(name string, exits string, area *Area)   {}
func FindRoom(name string, exits string) *Room        { return nil }
func ShiftRoom(room *Room, direction Direction) *Room { return nil }
func AddArea(name string)
