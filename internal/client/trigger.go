package client

import (
	"fmt"
	"regexp"
	"strconv"
)

type TriggerFunc func(*regexp.Regexp, []string)
type Trigger struct {
	Re  *regexp.Regexp
	Cmd TriggerFunc
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
