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
	C.AddAlias(`^#bot start (?P<area>.+)$`, BotStart)
}

var visited = []*mapper.Room{}
var rooms = []*mapper.Room{}

func BotStart(t *trigger.Trigger) {
}
