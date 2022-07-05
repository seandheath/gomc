package client

import (
	"fmt"
	"io"
	"net"
	"regexp"
	"runtime"
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
	AddAlias("^#memstats$", MemStatsCmd)
}

// ConnectCmd takes a string from the user and attempts to ConnectCmd to the mud server.
// If the connection is successful then a goroutine is launched to handle the connection.
var ConnectCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	if Conn != nil {
		ShowMain("Already connected.\n")
		return
	}
	text := matches[1]
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

var CaptureCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	s := strings.TrimPrefix(matches[0], "#capture ")

	if s == "overhead" {
		ShowOverhead(CurrentRaw)
		Gag = true
	} else {
		ts := time.Now().Format("2006:01:02 15:04:05")
		ShowChat(fmt.Sprintf("[%s] %s", ts, CurrentRaw))
	}
}

var FuncCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	s := strings.TrimPrefix(matches[0], "#func ")
	if f, ok := fmap[s]; ok { // Found the function
		f(re, matches)
	} else {
		ShowMain("Function not found:" + matches[0] + "\n")
	}
}

var MatchCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	CheckTriggers(actions, matches[1])
}

// Nested triggers overwrite matches... need to pass into the function
var LoopCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	n, err := strconv.Atoi(matches[1])
	if err != nil {
		ShowMain("Error parsing loop number: " + err.Error() + "\n")
	}
	for i := 0; i < n; i++ {
		Parse(matches[2])
	}
}

var BaseActionCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	showtriggers(actions, "actions")
}

var AddActionCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	AddActionString(matches[1], matches[2])
}

var UnactionCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	actions = untrigger(actions, "action", matches[1])
}

var BaseAliasCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	showtriggers(aliases, "aliases")
}

var AddAliasCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	AddAliasString(matches[1], matches[2])
}

var UnaliasCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	aliases = untrigger(aliases, "alias", matches[1])
}

var MemStatsCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	ShowMain(fmt.Sprintf("Alloc: %d MiB", m.Alloc/1024/1024))
}

func showtriggers(t []Trigger, triggerType string) {
	ShowMain("## Current " + triggerType + ":\n")
	for i, a := range t {
		ShowMain(fmt.Sprintf("\n[%d]: %s", i, a.Re.String()))
	}
	ShowMain("\n")
}

func untrigger(triggerList []Trigger, triggerType string, index string) []Trigger {
	n, err := strconv.Atoi(index)
	if err != nil {
		ShowMain(fmt.Sprintf("Invalid %s number: %d\n", triggerType, n))
		return triggerList
	}
	if n >= len(actions) {
		ShowMain(fmt.Sprintf("%s not found: %d\n", triggerType, n))
		return triggerList
	}
	return append(triggerList[:n], triggerList[n+1:]...)
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
