package client

import (
	"strings"
)

type ActionHandlerStruct struct{}

var ActionHandler = &ActionHandlerStruct{}

var DoCapture = false
var Gag = false
var CaptureFunc func(string)

func (a *ActionHandlerStruct) Write(p []byte) (int, error) {
	n := len(p)

	lines := strings.Split(string(p), "\r")
	for _, line := range lines {
		for _, a := range Actions {
			if a.re.MatchString(line) {
				Parse(a.cmd)
			}
		}
		if DoCapture {
			CaptureFunc(line)
			DoCapture = false
		}
		if !Gag {
			ShowMain(line)
		} else {
			ShowMain("\n")
			Gag = false
		}
	}
	return n, nil
}
