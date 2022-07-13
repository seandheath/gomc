package nodeka

import (
	"strings"

	"github.com/seandheath/go-mud-client/internal/client"
)

func initOmap() {
	Client.AddAction("^[ vi`~!@#$%^&*()-_=+\\[\\]{};:''\",.<>\\?|\\/]{34,37}$", MapLine)
	Client.AddAction("^$", EmptyLine)
	Client.AddAction("\\[ exits: ", ExitLine)
}

var inMap = false
var lineCount = 0
var mapLine = ""

func MapLine(t *client.TriggerMatch) {
	inMap = true
	if lineCount > 14 {
		// Final empty line
		// TODO fix for extended map
		if t.Matches[0] == "                                    " {
			lineCount = 0
			Client.Show("omap", strings.TrimSuffix(mapLine, "\n"))
			mapLine = ""
		} else {
			// something went wrong, mangled line?
			lineCount = 0
			mapLine = ""
			for i := 0; i < 16; i++ {
				Client.Show("omap", "\n")
			}
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

func ExitLine(t *client.TriggerMatch) {
	if lineCount != 0 {
		lineCount = 0
		mapLine = ""
		for i := 0; i < 16; i++ {
			Client.Show("omap", "\n")
		}
	}
}
