package mapper

import (
	"fmt"
)

// print takes the width and height of the desired map string and returns
// a map layout centered on your current room in one string. The map key
// can be found in map.md
func (m *Map) print(width, height int, unicode bool) []rune {
	s := []rune{}
	if m.room == nil {
		// TODO better error message
		s = append(s, []rune("Map location unknown")...)
		for i := 1; i < height; i++ {
			s = append(s, '\n')
		}
		return s
	}

	if m.Debug {
		s = append(s, []rune(fmt.Sprintf("%s:%d:%s\n", m.room.area.Name, m.room.ID, m.room.Name))...)
		//s = append(s, []rune("Name: "+m.room.Name+"\n")...)
		//s = append(s, []rune(fmt.Sprintf("ID: %d\n", m.room.ID))...)
		//s = append(s, []rune("Coordinates: "+m.room.Coordinates.String()+"\n")...)
		//s = append(s, []rune(fmt.Sprintf("Path: %d\n", len(m.nextMoves)))...)
		height -= 1
	}

	roomsize := 3
	if unicode {
		roomsize = 1
	}
	nx := width / roomsize  // number of rooms we can fit into the map width
	ny := height / roomsize // number of rooms we can fit into the map height
	cx := nx / 2            // X coordinates of the center room
	cy := ny / 2            // Y coordinates of the center room

	if nx <= 0 || ny <= 0 {
		return []rune{}
	}

	// initialize map string array
	roomArray := make([][]*Room, ny)
	collisionArray := make([][]bool, ny)
	for row := 0; row < ny; row++ {
		roomArray[row] = make([]*Room, nx)
		collisionArray[row] = make([]bool, nx)
	}

	// Starting at the top left populate each room string which will consist of
	// 9 characters indicating the top, middle, and bottom row of the 3x3 character
	// room string
	for row := 0; row < ny; row++ {
		for col := 0; col < nx; col++ {
			// Gets the room at the cooridnate offset from the current room and
			// on the same Z axis
			rs := m.GetRoomAtCoordinates(m.room.area, Coordinates{
				(m.room.Coordinates.X - cx) + col,
				(m.room.Coordinates.Y + cy) - row,
				m.room.Coordinates.Z,
			})
			if len(rs) <= 0 {
				// No room
				roomArray[row][col] = nil
			} else if len(rs) > 1 {
				// Multiple rooms at these coordinates, not sure how to display
				// a collision yet but we'll figure it out... maybe we'll print
				// the one with the shortest path to the current room or have
				// a collision indicator on that room?
				roomArray[row][col] = nil
				collisionArray[row][col] = true
			} else {
				// Only one room at those coordinates
				roomArray[row][col] = rs[0]
				collisionArray[row][col] = false
			}
		}
	}

	// Go through each row, col and collect the top, middle, bottom strings
	// for each room into sa. At the end of each row, concatenate the
	// top, middle, and bottom strings with newlines and append them to s.
	for row := 0; row < ny; row++ {
		center := false
		if row == cy {
			center = true
		}

		var r []rune
		if unicode {
			r = getUnicodeRow(roomArray[row], center)
		} else {
			r = getAsciiRow(roomArray[row], collisionArray[row], center)
		}
		s = append(s, r...)
	}
	return s
}

func getUnicodeRow(row []*Room, center bool) []rune {
	rs := []rune{}
	for col := 0; col < len(row); col++ {
		b := UnicodeRoom(row[col])
		rs = append(rs, b)

	}
	return rs
}

func getAsciiRow(row []*Room, collision []bool, center bool) []rune {
	// Now we have a 2D array of rooms that will fit into the width/height
	// provided. We need to generate a string from them. I could do all of this
	// in the above loop, but I'd like to have the populated array for debugging
	sa := make([][]rune, 3)
	sa[0] = []rune{}
	sa[1] = []rune{}
	sa[2] = []rune{}
	for col := 0; col < len(row); col++ {
		rs := AsciiRoom(row[col])
		// We have a collision so we'll mark an asterisk on top right
		if collision[col] {
			rs[0][0] = '*'
		}
		if center && col == len(row)/2 {
			// At the center room
			rs[1][1] = '#'
		}
		for subrow := 0; subrow < 3; subrow++ {
			// Each room has 3 rows
			sa[subrow] = append(sa[subrow], rs[subrow]...)
		}
	}
	s := []rune{}
	for i := 0; i < 3; i++ {
		s = append(s, sa[i]...)
		s = append(s, '\n')
	}
	return s
}

func getExitChar(r *Room, dir Direction) rune {
	switch dir {
	case North, East, South, West:
		nr := r.exits[dir]
		if nr != nil {
			if nr.area != r.area {
				return '?'
			}
		}
	}
	if _, ok := r.Doors[dir]; ok {
		return '+'
	}
	switch dir {
	case Up:
		return '^'
	case Down:
		return 'v'
	case North, South:
		return '|'
	case East, West:
		return '-'
	}
	return ' '
}

func UnicodeRoom(r *Room) rune {
	if r == nil {
		return ' '
	}
	n := checkmap(North, r.exits)
	e := checkmap(East, r.exits)
	s := checkmap(South, r.exits)
	w := checkmap(West, r.exits)

	if n && e && s && w {
		return '\u253C'
	}
	if n && e && s && !w {
		return '\u251C'
	}
	if n && e && !s && w {
		return '\u2534'
	}
	if n && e && !s && !w {
		return '\u2514'
	}
	if n && !e && s && w {
		return '\u2524'
	}
	if n && !e && s && !w {
		return '\u2502'
	}
	if n && !e && !s && !w {
		return '\u2575'
	}
	if !n && e && s && w {
		return '\u252C'
	}
	if !n && e && s && !w {
		return '\u250C'
	}
	if !n && e && !s && w {
		return '\u2500'
	}
	if !n && e && !s && !w {
		return '\u2576'
	}
	if !n && !e && s && w {
		return '\u2510'
	}
	if !n && !e && s && !w {
		return '\u2577'
	}
	if !n && !e && !s && w {
		return '\u2574'
	}
	return ' '
}

func checkmap(d Direction, m map[Direction]*Room) bool {
	_, ok := m[d]
	return ok
}

// Each room is represented by a 3x3 array of characters indicating exits
func AsciiRoom(r *Room) [][]rune {
	rs := make([][]rune, 3)
	for i := 0; i < 3; i++ {
		rs[i] = make([]rune, 3)
		rs[i][0] = ' ' // Set all values to space by default
		rs[i][1] = ' '
		rs[i][2] = ' '
	}
	if r != nil {
		for dir := range r.exits {
			switch dir {
			case Up:
				rs[0][2] = getExitChar(r, Up)
			case North:
				rs[0][1] = getExitChar(r, North)
			case West:
				rs[1][0] = getExitChar(r, West)
			case East:
				rs[1][2] = getExitChar(r, East)
			case Down:
				rs[2][0] = getExitChar(r, Down)
			case South:
				rs[2][1] = getExitChar(r, South)
			}
		}
		rs[1][1] = 'o'
	}

	return rs
}

func (m *Map) Show(name string) {
	w, h := C.GetWindowSize(name)
	s := m.print(w, h, false)
	C.PrintTo(name, string(s))
}
