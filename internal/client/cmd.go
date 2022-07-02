package client

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/rivo/tview"
)

type Writer struct{}

func CmdInit() {
	AddAlias("^#connect (.*)$", ConnectCmd)
	AddAlias("^#capture ?(.*)$", CaptureCmd)
	AddAlias("^#func (.*)$", FuncCmd)
	AddAlias("^#match (.+)$", MatchCmd)
	AddAlias("^#(\\d+) (.+)$", LoopCmd)
	AddAlias("^#action$", BaseActionCmd)
	AddAlias("^#action {(.+)}{(.+)}$", AddActionCmd)
	AddAlias("^#unaction (\\d+)$", UnactionCmd)
	AddAlias("^#alias$", BaseAliasCmd)
	AddAlias("^#alias {(.+)}{(.+)}$", AddAliasCmd)
	AddAlias("^#unalias (\\d+)$", UnaliasCmd)
}

func BaseActionCmd() {
	showtriggers(actions, "actions")
}
func AddActionCmd() {
	AddAction(CurrentMatches[1], CurrentMatches[2])
}
func UnactionCmd() {
	untrigger(actions, "action")
}
func BaseAliasCmd() {
	showtriggers(aliases, "aliases")
}
func AddAliasCmd() {
	AddAlias(CurrentMatches[1], CurrentMatches[2])
}
func UnaliasCmd() {
	untrigger(aliases, "alias")
}

func showtriggers(t []Trigger, triggerType string) {
	ShowMain("## Current " + triggerType + ":\n")
	for i, a := range t {
		ShowMain(fmt.Sprintf("\n[%d]: %s", i, a.Re.String()))
	}
	ShowMain("\n")
}

func untrigger(triggerList []Trigger, triggerType string) {
	n, err := strconv.Atoi(CurrentMatches[1])
	if err != nil {
		ShowMain(fmt.Sprintf("Invalid %s number: %s\n", triggerType, CurrentMatches[1]))
		return
	}
	if n >= len(actions) {
		ShowMain(fmt.Sprintf("%s not found: %s\n", triggerType, CurrentMatches[1]))
		return
	}
	tmp := append(triggerList[:n], triggerList[n+1:]...)
	triggerList = tmp
}

// Nested triggers overwrite matches... need to pass into the function
func LoopCmd() {
	n, err := strconv.Atoi(CurrentMatches[1])
	if err != nil {
		ShowMain("Error parsing loop number: " + err.Error() + "\n")
	}
	for i := 0; i < n; i++ {
		Parse(CurrentMatches[2])
	}
}

// ConnectCmd takes a string from the user and attempts to ConnectCmd to the mud server.
// If the connection is successful then a goroutine is launched to handle the connection.
func ConnectCmd() {
	if Conn != nil {
		ShowMain("Already connected.\n")
		return
	}
	text := strings.TrimPrefix(CurrentMatches[0], "#connect ")
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
		line = strings.ReplaceAll(line, ";", ":") // Stops trigger abuse // TODO make config for this
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

func CaptureCmd() {
	s := strings.TrimPrefix(CurrentMatches[0], "#capture ")

	if s == "overhead" {
		ShowOverhead(CurrentRaw)
		Gag = true
	} else {
		ts := time.Now().Format("2006:01:02 15:04:05")
		ShowChat(fmt.Sprintf("[%s] %s", ts, CurrentRaw))
	}
}

func FuncCmd() {
	s := strings.TrimPrefix(CurrentMatches[0], "#func ")
	if f, ok := fmap[s]; ok { // Found the function
		f()
	} else {
		ShowMain("Function not found:" + CurrentMatches[0] + "\n")
	}
}

func MatchCmd() {
	CheckTriggers(actions, CurrentMatches[1])
}
