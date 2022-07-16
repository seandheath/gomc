package nodeka

import (
	"github.com/seandheath/gomc/pkg/trigger"
	"github.com/seandheath/gomc/pkg/util"
)

const OSIZE = 15

func initOmap() {
	C.AddActionFunc("^[ vi`~!@#$%^&*()-_=+\\[\\]{};:'\",.<>\\?|\\/]{36,37}$", MapLine)
	C.AddActionFunc("^$", EmptyLine)
	C.AddActionFunc("\\[ exits: ", ExitLine)
	C.AddActionFunc(`^\[Reply:`, ReplyLine)
}

var inMap = false
var lineCount = 0
var mapLine []byte

func MapLine(t *trigger.Trigger) {
	inMap = true
	if lineCount < OSIZE {
		mapLine = append(mapLine, C.RawLine...)
		lineCount += 1
	} else {
		lineCount = 0
		C.PrintBytesTo("omap", util.TrimEnd(mapLine))
		mapLine = []byte("  ")
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
		for i := 0; i <= OSIZE; i++ {
			C.PrintTo("omap", "\n")
		}
	}
}

func ReplyLine(t *trigger.Trigger) {
	inMap = false
	if lineCount != 0 {
		lineCount = 0
		for i := 0; i <= OSIZE; i++ {
			//C.PrintTo("omap", "\n")
		}
	}
}
