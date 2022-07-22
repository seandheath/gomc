package mapper

import (
	"fmt"
	"time"
)

func GetPath(start *Room, finish *Room) []byte {
	tstart := time.Now()
	path := dijkstra(start, finish)
	if path == nil {
		C.Print("\nError finding path between rooms.\n")
	}
	C.Print(fmt.Sprintf("\nPath found in %dus.\n", time.Since(tstart).Microseconds()))
	return path
}

func dijkstra(start *Room, finish *Room) []byte {
	openSet := []*Room{}
	cameFrom := map[*Room]*Room{}
	gScore := map[*Room]int{}

	openSet = append(openSet, start)
	cameFrom[start] = nil

	for len(openSet) > 0 {
		currentRoom := openSet[0]
		if currentRoom == finish {
			return rebuildPath(start, cameFrom, currentRoom)
		}
		openSet = openSet[1:]
		for _, exit := range currentRoom.exits {
			if exit != nil {
				score := gScore[currentRoom] + 1
				if lastScore, ok := gScore[exit]; ok {
					if score < lastScore {
						gScore[exit] = score
						cameFrom[exit] = currentRoom
					}
				} else {
					// never scored this exit before
					gScore[exit] = score
					cameFrom[exit] = currentRoom
				}
				if !contains(openSet, exit) {
					openSet = append(openSet, exit)
				}
			}
		}
	}
	C.Print("\nUnable to find a path between the provided rooms.\n")
	return nil
}

func rebuildPath(start *Room, cameFrom map[*Room]*Room, currentRoom *Room) []byte {
	path := []byte{}
	for currentRoom != start {
		dir := getDirection(cameFrom[currentRoom], currentRoom)
		if dir == "" {
			C.Print("\nError rebuilding path between rooms.\n")
			return nil
		}
		path = append([]byte{shortdirs[dir]}, path...)
		if door, ok := cameFrom[currentRoom].Doors[dir]; ok {
			path = append([]byte(";open "+string(dir)+"."+door.Name+";"), path...)
		}
		currentRoom = cameFrom[currentRoom]
	}
	return path
}

func getDirection(from *Room, to *Room) Direction {
	for dir, e := range from.exits {
		if e == to {
			return dir
		}
	}
	return ""
}

// Only works for rooms in the same area
func astar(start *Room, finish *Room) []byte {
	openSet := []*Room{}
	cameFrom := map[*Room]*Room{}
	gMap := map[*Room]int{}
	//fMap := map[*Room]int{}

	openSet = append(openSet, start)

	for len(openSet) > 0 {
		currentRoom := openSet[0]
		if currentRoom == finish {
			return rebuildPath(start, cameFrom, currentRoom)
		}
		openSet = openSet[1:]
		for _, exit := range currentRoom.exits {
			gscore := gMap[currentRoom] + 1
			if gscore < gMap[exit] {
				cameFrom[exit] = currentRoom
				gMap[exit] = gscore
				//fMap[exit] = gscore + heuristic(exit, finish)
			}
			if !contains(openSet, exit) {
				openSet = append(openSet, exit)
			}
		}
	}
	C.Print("\nUnable to find a path between the provided rooms.\n")
	return nil

}

func contains(rl []*Room, r *Room) bool {
	for _, cr := range rl {
		if cr == r {
			return true
		}
	}
	return false
}

func getRoomsWithAreaExits(a *Area) []*Room {
	exits := []*Room{}
	for _, r := range a.Rooms {
		if hasAreaExit(r) {
			exits = append(exits, r)
		}
	}
	return exits
}

func hasAreaExit(r *Room) bool {
	for _, e := range r.exits {
		if e != nil {
			if e.area != r.area {
				return true
			}
		}
	}
	return false
}
