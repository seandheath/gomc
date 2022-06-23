package client

import (
	"io"
	"net"

	"github.com/rivo/tview"
)

type Cmd struct {
	Name string
	Run  func(string)
	Help string
}

var Commands map[string]Cmd = make(map[string]Cmd)

func init() {
	Commands["#connect"] = Cmd{"#connect", connect, "#connect <host>:<port>"}
	Commands["#capture"] = Cmd{"#capture", capture, "#capture [chat|overhead]"}
}

func AddCommand(cmd Cmd) {
	if _, ok := Commands[cmd.Name]; ok {
		ShowMain("Command already exists.\n")
	} else {
		Commands[cmd.Name] = cmd
	}
}

func connect(text string) {
	conn, err := net.Dial("tcp", text)
	if err != nil {
		ShowMain("Failed to connect: " + err.Error() + "\n")
	}
	Connection = conn
	Variables["raw"] = ""
	Variables["text"] = ""
	go func() {
		defer Connection.Close()
		w := tview.ANSIWriter(ActionHandler)
		if _, err := io.Copy(w, Connection); err != nil {
			ShowMain("Connection closed: " + err.Error() + "\n")
			return
		}
	}()
}

func capture(text string) {
	DoCapture = true
	Gag = true
	if text == "overhead" {
		CaptureFunc = ShowOverhead
	} else {
		CaptureFunc = ShowChat
	}
}
