package client

import (
	"io"
	"net"
	"strings"

	"github.com/rivo/tview"
)

type Writer struct{}

func CmdInit() {
	AddAlias("^#connect .*$", connect)
	AddAlias("^#capture.*$", Capture)
	AddAlias("^#func .*$", ExecuteFunction)
	AddAlias("^#match (.+)$", MatchTrigger)
}

// connect takes a string from the user and attempts to connect to the mud server.
// If the connection is successful then a goroutine is launched to handle the connection.
func connect(args []string) {
	if Conn != nil {
		ShowMain("Already connected.\n")
		return
	}
	text := strings.TrimPrefix(args[0], "#connect ")
	conn, err := net.Dial("tcp", text)
	if err != nil {
		ShowMain("Failed to connect: " + err.Error() + "\n")
	}
	Conn = conn
	go func() {
		defer Conn.Close()
		writer := &Writer{}
		w := tview.ANSIWriter(writer)
		if _, err := io.Copy(w, Conn); err != nil {
			ShowMain("Connection closed: " + err.Error() + "\n")
			Conn = nil
		}
	}()
}

func (w *Writer) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\r")
	for _, line := range lines {
		handleLine(line)
	}
	return len(p), nil
}

func handleLine(line string) {
	CurrentRaw = line
	CurrentText = stripTags(CurrentRaw)
	CheckTriggers(actions, strings.TrimSuffix(CurrentText, "\n"))
	if !Gag {
		ShowMain(CurrentRaw)
	} else {
		Gag = false
	}
}

// stripTags strips colour tags from the given string. (Region tags are not
// stripped.)
func stripTags(text string) string {
	stripped := colorPattern.ReplaceAllStringFunc(text, func(match string) string {
		if len(match) > 2 {
			return ""
		}
		return match
	})
	return escapePattern.ReplaceAllString(stripped, `[$1$2]`)
}

func Capture(args []string) {
	s := strings.TrimPrefix(args[0], "#capture ")

	if s == "overhead" {
		ShowOverhead(CurrentRaw)
		Gag = true
	} else {
		ShowChat(CurrentRaw)
	}
}

func ExecuteFunction(args []string) {
	s := strings.TrimPrefix(args[0], "#func ")
	if f, ok := fmap[s]; ok { // Found the function
		f([]string{CurrentText})
	} else {
		ShowMain("Function not found:" + args[0] + "\n")
	}
}

func MatchTrigger(args []string) {
	CheckTriggers(actions, args[1])
}
