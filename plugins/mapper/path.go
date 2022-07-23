package mapper

import (
	"fmt"
	"time"
)

func GetPath(start *Room, finish *Room) ([]byte, int) {
	tstart := time.Now()
	path, length := dijkstra(start, finish)
	if path == nil {
		C.Print("\nError finding path between rooms.\n")
		return nil, 0
	}
	C.Print(fmt.Sprintf("\nPath found in %dus.\n", time.Since(tstart).Microseconds()))
	return path, length
}

func dijkstra(start *Room, finish *Room) ([]byte, int) {
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
	return nil, 0
}

func rebuildPath(start *Room, cameFrom map[*Room]*Room, currentRoom *Room) ([]byte, int) {
	path := []byte{}
	length := 0
	for currentRoom != start {
		dir := getDirection(cameFrom[currentRoom], currentRoom)
		if dir == "" {
			C.Print("\nError rebuilding path between rooms.\n")
			return nil, 0
		}
		path = append([]byte{shortdirs[dir]}, path...)
		if door, ok := cameFrom[currentRoom].Doors[dir]; ok {
			path = append([]byte(";open "+string(dir)+"."+door.Name+";"), path...)
		}
		currentRoom = cameFrom[currentRoom]
		length++
	}
	return path, length
}

func getDirection(from *Room, to *Room) Direction {
	for dir, e := range from.exits {
		if e == to {
			return dir
		}
	}
	return ""
}

func contains(rl []*Room, r *Room) bool {
	for _, cr := range rl {
		if cr == r {
			return true
		}
	}
	return false
}
