package mapper

import (
	"bytes"
	"fmt"
)

// print takes the width and height of the desired map string and returns
// a map layout centered on your current room in one string. The map key
// can be found in map.md
func (m *Map) print(width, height int) []byte {
	s := bytes.Buffer{}
	if m.room == nil {
		// TODO better error message
		s.WriteString("Map location unknown")
		for i := 1; i < height; i++ {
			s.WriteByte('\n')
		}
		return s.Bytes()
	}

	if m.Debug {
		s.WriteString("Area: " + m.room.area.Name + "\n")
		s.WriteString("Name: " + m.room.Name + "\n")
		s.WriteString(fmt.Sprintf("ID: %d\n", m.room.ID))
		s.WriteString("Coordinates: " + m.room.Coordinates.String() + "\n")
		s.WriteString(fmt.Sprintf("Path: %d\n", len(m.nextMoves)))
		height -= 5
	}

	nx := width / 3  // number of rooms we can fit into the map width
	ny := height / 3 // number of rooms we can fit into the map height
	cx := nx / 2     // X coordinates of the center room
	cy := ny / 2     // Y coordinates of the center room

	if nx <= 0 || ny <= 0 {
		return []byte{}
	}

	// initialize map string array
	ma := make([][]*Room, ny)
	cl := make([][]bool, ny)
	for row := 0; row < ny; row++ {
		ma[row] = make([]*Room, nx)
		cl[row] = make([]bool, nx)
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
				ma[row][col] = nil
			} else if len(rs) > 1 {
				// Multiple rooms at these coordinates, not sure how to display
				// a collision yet but we'll figure it out... maybe we'll print
				// the one with the shortest path to the current room or have
				// a collision indicator on that room?
				ma[row][col] = nil
				cl[row][col] = true
			} else {
				// Only one room at those coordinates
				ma[row][col] = rs[0]
				cl[row][col] = false
			}
		}
	}

	// Now we have a 2D array of rooms that will fit into the width/height
	// provided. We need to generate a string from them. I could do all of this
	// in the above loop, but I'd like to have the populated array for debugging
	sa := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		sa[i] = make([]byte, 3)
	}

	// Go through each row, col and collect the top, middle, bottom strings
	// for each room into sa. At the end of each row, concatenate the
	// top, middle, and bottom strings with newlines and append them to s.
	for row := 0; row < ny; row++ {
		sa[0] = []byte("   ")
		sa[1] = []byte("   ")
		sa[2] = []byte("   ")
		for col := 0; col < nx; col++ {
			rs := ma[row][col].MapStrings()
			// We have a collision so we'll mark an asterisk on top right
			if cl[row][col] {
				sa[0][2] = '*'
			}
			if row == cy && col == cx {
				// At the center room
				rs[1][1] = '#'
			}
			for subrow := 0; subrow < 3; subrow++ {
				// Each room has 3 rows
				sa[subrow] = append(sa[subrow], rs[subrow]...)
			}
		}
		// Combine the three rows of each string into our main string
		for i := 0; i < 3; i++ {
			s.Write(sa[i])
			s.WriteByte('\n')
		}
		//s += strings.Join(sa, "\n")
	}

	return s.Bytes() // You should at least have a bunch of spaces and newlines
}

func getExitChar(r *Room, dir Direction) byte {
	switch dir {
	case North, East, South, West:
		nr := r.exits[dir]
		if nr != nil {
			if nr.area != r.area {
				return '?'
			}
		}
	}
	switch dir {
	case Up:
		return '^'
	case North:
		return '|'
	case East:
		return '-'
	case South:
		return '|'
	case West:
		return '-'
	case Down:
		return 'v'
	}
	return ' '
}

// Each room is represented by a 3x3 array of characters indicating exits
func (r *Room) MapStrings() [][]byte {
	rs := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		rs[i] = make([]byte, 3)
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
	C.PrintBytesTo(name, m.print(w, h))
}
