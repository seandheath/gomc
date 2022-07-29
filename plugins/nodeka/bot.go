package nodeka

import (
	"github.com/seandheath/gomc/pkg/trigger"
	"github.com/seandheath/gomc/plugins/mapper"
)

type Mob struct {
	Single   string `yaml:"single"`
	Multiple string `yaml:"multiple"`
}

type BotConfig struct {
	Area string          `yaml:"area"`
	Mobs map[string]*Mob `yaml:"mobs"`
}

func initBot() {
	C.AddAlias(`^#bot start (?P<tag>.+)$`, BotStart)
}

var botPath []*mapper.Room

func BotStart(t *trigger.Trigger) {
	// Get the current map
	m := mapper.GetMap()
	if m == nil || m.Room == nil {
		C.Print("You are not in a map.")
		return
	}
	// Get the tag to start the bot at
	tag := t.Results["tag"]
	// We'll automatically pick the first room in the tag list
	tagList := m.GetRoomsByTag(tag)
	// We have at least one room with the tag
	if tagList[0] == nil {
		C.Print("\nMAP: No rooms found with tag: " + tag)
		return
	}
	// Get the minimum spanning tree of the area starting at the tagged room
	botPath = getMST(tagList[0])
	if len(botPath) == 0 {
		C.Print("\nMAP: Unable to generate MST for area\n")
		return
	}

	// Get the path from current location to the first room in the MST
	pathcmd, len := m.GetPath(m.Room, botPath[0])
	if len == 0 {
		// We're already there!
		Begin()
	} else {
		// We need to move to the first room in the MST
		C.Print("\nMAP: Moving to first room in MST\n")
		m.PathQ.Append(func() { Begin() })
		m.Walking = true
		C.Parse(string(pathcmd))

	}
}

func Begin() {
	C.Parse("dance")
}

func getMST(tgt *mapper.Room) []*mapper.Room {
	openSet := []*mapper.Room{tgt}
	visited := []*mapper.Room{}
	for len(openSet) > 0 {
		for _, ex := range openSet[0].Exits {
			if ex != nil && !contains(visited, ex) && !contains(openSet, ex) && ex.Area == tgt.Area {
				openSet = append(openSet, ex)
			}
		}
		// Remove the room we just visited
		visited = append(visited, openSet[0])
		openSet = openSet[1:]
	}
	return visited
}

func contains(set []*mapper.Room, tgt *mapper.Room) bool {
	for _, r := range set {
		if r == tgt {
			return true
		}
	}
	return false
}
