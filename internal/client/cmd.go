package client

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func cmdInit() {
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

var CaptureCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	s := strings.TrimPrefix(matches[0], "#capture ")

	if s == "overhead" {
		ShowMain(CurrentRaw)
		Gag = true
	} else {
		ts := time.Now().Format("2006:01:02 15:04:05")
		ShowMain(fmt.Sprintf("[%s] %s", ts, CurrentRaw))
	}
}

var FuncCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	s := strings.TrimPrefix(matches[0], "#func ")
	if f, ok := functions[s]; ok { // Found the function
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

var MemStatsCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	ShowMain(fmt.Sprintf("Alloc: %d MiB", m.Alloc/1024/1024))
}
