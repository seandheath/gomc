package nodeka

import (
	"strings"

	"github.com/seandheath/gomc/pkg/trigger"
)

const OSIZE = 15

func initOmap() {
	C.AddAction("^[ vi`~!@#$%^&*()-_=+\\[\\]{};:''\",.<>\\?|\\/]{34,37}$", MapLine)
	C.AddAction("^$", EmptyLine)
	C.AddAction("\\[ exits: ", ExitLine)
	C.AddAction(`^\[Reply:`, ReplyLine)
}

var inMap = false
var lineCount = 0
var mapLine = ""

func MapLine(t *trigger.Trigger) {
	inMap = true
	if lineCount < OSIZE {
		lineCount += 1
		mapLine += C.RawLine
	} else {
		// Final empty line
		// TODO fix for extended map
		//if t.Matches[0] == "                                    " {
		lineCount = 0
		C.PrintTo("omap", strings.TrimSuffix(mapLine, "\n"))
		mapLine = ""
		//} else {
		// something went wrong, mangled line?
		//lineCount = 0
		//mapLine = ""
		//for i := 0; i <= OSIZE; i++ {
		//C.PrintTo("omap", "\n")
		//}
		//}
	}
	C.Gag = true
}

func EmptyLine(t *trigger.Trigger) {
	if inMap {
		inMap = false
		C.Gag = true
	}
}

func ExitLine(t *trigger.Trigger) {
	if lineCount != 0 {
		lineCount = 0
		mapLine = ""
		for i := 0; i <= OSIZE; i++ {
			C.PrintTo("omap", "\n")
		}
	}
}

func ReplyLine(t *trigger.Trigger) {
	inMap = false
	if lineCount != 0 {
		lineCount = 0
		mapLine = ""
		for i := 0; i <= OSIZE; i++ {
			C.PrintTo("omap", "\n")
		}
	}
}
