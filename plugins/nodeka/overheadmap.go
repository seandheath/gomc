package nodeka

import (
	"github.com/seandheath/go-mud-client/internal/client"
)

var inMap = false

func initOmap() {
	Client.AddAction("^[ vi`~!@#$%^&*()-_=+\\[\\]{};:''\",.<>\\?|\\/]{34,37}$", MapLine)
	Client.AddAction("^$", EmptyLine)
}

// Handle seeing an overhead map line
func MapLine(t *client.TriggerMatch) {
	inMap = true
	Client.Show("omap", Client.RawLine)
	Client.Gag = true
}

// Handle seeing an empty line
func EmptyLine(t *client.TriggerMatch) {
	if inMap {
		inMap = false
		Client.Gag = true
	}
}
