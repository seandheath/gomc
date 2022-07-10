package nodeka

import (
	"strings"

	"github.com/seandheath/go-mud-client/internal/client"
)

func initOmap() {
	Client.AddAction("^[ vi`~!@#$%^&*()-_=+\\[\\]{};:''\",.<>\\?|\\/]{34,37}$", MapLine)
	Client.AddAction("^$", EmptyLine)
}

var inMap = false
var lineCount = 0
var mapLine = ""

func MapLine(t *client.TriggerMatch) {
	inMap = true
	if lineCount > 14 {
		// Final empty line
		if t.Matches[0] == "                                    " {
			lineCount = 0
			Client.Show("omap", strings.TrimSuffix(mapLine, "\n"))
		} else {
			// something went wrong, mangled line?
			lineCount = 0
			mapLine = ""
		}
	} else {
		lineCount += 1
		mapLine += Client.RawLine
	}
	Client.Gag = true
}

func EmptyLine(t *client.TriggerMatch) {
	if inMap {
		inMap = false
		Client.Gag = true
	}
}
