package nodeka

import (
	"github.com/seandheath/go-mud-client/internal/client"
)

var inMap = false

func initOmap() {
	Client.AddAction("^[ vi`~!@#$%^&*()-_=+\\[\\]{};:''\",.<>\\?|\\/]{34,37}$", MapLine)
	Client.AddAction("^ {36}$", SpaceLine)
	Client.AddAction("^$", EmptyLine)
}

// Handle seeing an overhead map line
var lineCount = 0

func MapLine(t *client.TriggerMatch) {
	inMap = true
	Client.Show("omap", Client.RawLine)
	lineCount += 1
	Client.Gag = true
}

// Handle seeing an empty line
func EmptyLine(t *client.TriggerMatch) {
	if inMap {
		inMap = false
		Client.Gag = true
	}
}

func SpaceLine(t *client.TriggerMatch) {
	// End of map
	if lineCount == 16 {
		Client.Gag = true
		lineCount = 0
	}
}
