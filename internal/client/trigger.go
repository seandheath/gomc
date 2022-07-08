package client

import (
	"fmt"
	"regexp"
	"strconv"
)

type TriggerMatch struct {
	Matches []string
	Trigger *Trigger
}
type TriggerFunc func(*TriggerMatch)
type Trigger struct {
	Re  *regexp.Regexp
	Cmd TriggerFunc
}

func (c *Client) CheckTriggers(list []Trigger, text string) bool {
	matched := false
	for _, t := range list {
		m := t.Re.FindStringSubmatch(text)
		if len(m) > 0 {
			matched = true
			t.Cmd(&TriggerMatch{m, &t})
		}
	}
	return matched
}

func (c *Client) addTriggerString(list []Trigger, rs string, cmd string) []Trigger {
	f := func(t *TriggerMatch) {
		c.Parse(cmd)
	}
	return c.addTrigger(list, rs, f)
}

func (c *Client) getFunc(cmd string) TriggerFunc {

	return nil
}

// This function adds a trigger to the provided list and returns it
func (c *Client) addTrigger(list []Trigger, rs string, cmd TriggerFunc) []Trigger {
	re, err := regexp.Compile(rs)
	if err != nil {
		c.ShowMain("Error compiling trigger: " + err.Error() + "\n")
		return list
	}
	return append(list, Trigger{re, cmd})
}

func (c *Client) BaseActionCmd(t *TriggerMatch) {
	c.showtriggers(c.actions, "actions")
}

func (c *Client) AddActionCmd(t *TriggerMatch) {
	c.AddActionString(t.Matches[1], t.Matches[2])
}

func (c *Client) UnactionCmd(t *TriggerMatch) {
	c.actions = c.untrigger(c.actions, "action", t.Matches[1])
}

func (c *Client) BaseAliasCmd(t *TriggerMatch) {
	c.showtriggers(c.aliases, "aliases")
}

func (c *Client) AddAliasCmd(t *TriggerMatch) {
	c.AddAliasString(t.Matches[1], t.Matches[2])
}

func (c *Client) UnaliasCmd(t *TriggerMatch) {
	c.aliases = c.untrigger(c.aliases, "alias", t.Matches[1])
}

func (c *Client) showtriggers(t []Trigger, triggerType string) {
	c.ShowMain("## Current " + triggerType + ":\n")
	for i, a := range t {
		c.ShowMain(fmt.Sprintf("\n[%d]: %s", i, a.Re.String()))
	}
	c.ShowMain("\n")
}

func (c *Client) untrigger(triggerList []Trigger, triggerType string, index string) []Trigger {
	n, err := strconv.Atoi(index)
	if err != nil {
		c.ShowMain(fmt.Sprintf("Invalid %s number: %d\n", triggerType, n))
		return triggerList
	}
	if n >= len(c.actions) {
		c.ShowMain(fmt.Sprintf("%s not found: %d\n", triggerType, n))
		return triggerList
	}
	return append(triggerList[:n], triggerList[n+1:]...)
}
