package mapper

import (
	"fmt"
	"time"
)

func (m *Map) GetPath(start *Room, finish *Room) ([]byte, int) {
	tstart := time.Now()
	m.Path = dijkstra(start, finish)
	pathcmd := getPathCommands(m.Path)
	if len(pathcmd) == 0 || m.Path == nil {
		m.Path = nil
		m.Walking = false
		C.Print("\nError finding path between rooms.\n")
		return nil, 0
	}
	C.Print(fmt.Sprintf("\nPath found in %dus.\n", time.Since(tstart).Microseconds()))
	return pathcmd, len(m.Path)
}

func dijkstra(start *Room, finish *Room) []*Room {
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
		for _, exit := range currentRoom.Exits {
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

func getPathCommands(rs []*Room) []byte {
	pathcmd := []byte{}
	for i := 0; i < len(rs)-1; i++ {
		dir := getDirection(rs[i], rs[i+1])
		if dir == "" {
			C.Print("\nError getting path commands.\n")
			return nil
		}
		if door, ok := rs[i].Doors[dir]; ok {
			pathcmd = append(pathcmd, []byte(";open "+string(dir)+"."+door.Name+";")...)
		}
		pathcmd = append(pathcmd, shortdirs[dir])
	}
	return pathcmd
}

func rebuildPath(start *Room, cameFrom map[*Room]*Room, currentRoom *Room) []*Room {
	path := []*Room{currentRoom}
	for currentRoom != start {
		//dir := getDirection(cameFrom[currentRoom], currentRoom)
		//if dir == "" {
		//C.Print("\nError rebuilding path between rooms.\n")
		//return nil, 0
		//}
		path = append([]*Room{cameFrom[currentRoom]}, path...)
		//if door, ok := cameFrom[currentRoom].Doors[dir]; ok {
		//path = append([]byte(";open "+string(dir)+"."+door.Name+";"), path...)
		//}
		currentRoom = cameFrom[currentRoom]
	}
	return path
}

func getDirection(from *Room, to *Room) Direction {
	for dir, e := range from.Exits {
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
