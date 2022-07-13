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
	X int `yaml:"x"`
	Y int `yaml:"y"`
	Z int `yaml:"z"`
}

type Room struct {
	ID          int                 `yaml:"id"`
	Name        string              `yaml:"name"`
	Exits       map[Direction]*Room `yaml:"exits"`
	Area        *Area               `yaml:"area"`
	Coordinates Coordinates         `yaml:"coordinates"`
	Tags        []string            `yaml:"tags"`
}

type Area struct {
	Name  string        `yaml:"name"`
	Rooms map[int]*Room `yaml:"rooms"`
}

type Map struct {
	CurrentRoom *Room
	CurrentArea *Area
	Areas       map[string]Area `yaml:"areas"`
	Rooms       map[string]Room `yaml:"rooms"`
}

var reverse map[Direction]Direction = map[Direction]Direction{
	North: South,
	East:  West,
	South: North,
	West:  East,
	Up:    Down,
	Down:  Up,
}

func (m *Map) Load(path string) *Map                  { return nil }
func (m *Map) Save(path string) error                 { return nil }
func Move(direction Direction) *Room                  { return nil }
func AddRoom(name string, exits string, area string)  {}
func FindRoom(name string, exits string) *Room        { return nil }
func ShiftRoom(room *Room, direction Direction) *Room { return nil }
func AddArea(name string)
func Show(width int, height int) string { return "" }
