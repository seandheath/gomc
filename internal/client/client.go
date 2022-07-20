package client

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/seandheath/gomc/internal/tui"
	"github.com/seandheath/gomc/pkg/plugin"
	"github.com/seandheath/gomc/pkg/trigger"
)

const BUFFERSIZE = 1024 * 1024

type Client struct {
	conn          net.Conn
	timeout       time.Time
	printBuffer   []byte
	processBuffer []byte
	actions       []*trigger.Trigger
	aliases       []*trigger.Trigger
	functions     map[string]trigger.Func
	plugins       map[string]*plugin.Config
	tui           *tui.TUI

	// Publicly available variables
	Gag      bool
	RawLine  []byte
	TextLine []byte
	Var      map[string]string
}

func NewClient() *Client {
	c := &Client{}
	c.conn = nil
	c.Gag = false
	c.actions = make([]*trigger.Trigger, 0)
	c.aliases = make([]*trigger.Trigger, 0)
	c.functions = map[string]trigger.Func{}
	c.plugins = map[string]*plugin.Config{}
	c.tui = tui.NewTUI(c.Parse)
	c.cmdInit()
	return c
}

func (c *Client) AddAction(t *trigger.Trigger) {
	c.actions = c.addTrigger(c.actions, t)
}
func (c *Client) AddActionFunc(rs string, cmd trigger.Func) *trigger.Trigger {
	t := trigger.NewTrigger(rs, cmd)
	c.actions = c.addTrigger(c.actions, t)
	return t
}
func (c *Client) AddActionString(rs string, cmd string) *trigger.Trigger {
	t := trigger.NewTrigger(rs, func(t *trigger.Trigger) { c.Parse(cmd) })
	c.actions = c.addTrigger(c.actions, t)
	return t
}
func (c *Client) AddAlias(t *trigger.Trigger) {
	c.aliases = c.addTrigger(c.aliases, t)
}
func (c *Client) AddAliasFunc(rs string, cmd trigger.Func) *trigger.Trigger {
	t := trigger.NewTrigger(rs, cmd)
	c.aliases = c.addTrigger(c.aliases, t)
	return t
}
func (c *Client) AddAliasString(rs string, cmd string) *trigger.Trigger {
	t := trigger.NewTrigger(rs, func(t *trigger.Trigger) { c.Parse(cmd) })
	c.aliases = c.addTrigger(c.aliases, t)
	return t
}

// Parse the string and send the result to the server
func (c *Client) Parse(text string) {
	if c.CheckTriggers(c.aliases, text) { // Check for aliases / commands
		return
	} else if c.conn == nil { // Not connected yet
		c.Print("Not connected.\n")
		return
	} else if strings.Contains(text, ";") { // Allow splitting commands by ;
		s := strings.Split(text, ";")
		for _, t := range s {
			c.Parse(t)
		}
	} else {
		c.SendNow(text)
	}
}

func (c *Client) SendNow(text string) {
	if c.conn == nil {
		c.Print("Not connected.\n")
		return
	}
	c.Print(text + "\n")
	_, err := c.conn.Write([]byte(text + "\n"))
	if err != nil {
		c.Print("Error sending: " + err.Error() + "\n")
		c.conn = nil
	}
}

func (c *Client) Run() {
	c.tui.Run()
}

// AddFunction maps a string to a function so that you can call the function
// from the mud with #function <name>
func (c *Client) AddFunction(name string, f func(t *trigger.Trigger)) {
	c.functions[name] = f
}

// Print prints text on the main screen
func (c *Client) Print(text string) {
	c.PrintTo("main", text)
}

func (c *Client) PrintBytesTo(window string, b []byte) {
	c.tui.PrintBytes(window, b)
}

// PrintTo prints text to the specified window
func (c *Client) PrintTo(window string, text string) {
	c.tui.Print(window, text)
}

func (c *Client) LoadPlugin(name string, p *plugin.Config) {
	for re, cmd := range p.Actions {
		c.AddActionString(re, cmd)
	}
	for re, cmd := range p.Aliases {
		c.AddAliasString(re, cmd)
	}
	for n, f := range p.Functions {
		c.AddFunction(n, f)
	}
	for n, w := range p.Windows {
		c.tui.AddWindow(n, w)
	}

	c.tui.SetGrid(p.Grid.Rows, p.Grid.Columns)

	c.plugins[name] = p
}

func (c *Client) CheckTriggers(list []*trigger.Trigger, text string) bool {
	matched := false
	for _, t := range list {
		if t.Enabled {
			t.Matches = t.FindStringSubmatch(string(text))
			if len(t.Matches) > 0 {
				matched = true
				if len(t.SubexpNames()) > 0 {
					for i, m := range t.Matches {
						if i > 0 {
							t.Results[t.SubexpNames()[i]] = m
						}
					}
				}
				t.Do()
			}
		}
	}
	return matched
}

func (c *Client) getFunc(cmd string) trigger.Func {

	return nil
}

// This function adds a trigger to the provided list and returns the new trigger
func (c *Client) addTrigger(list []*trigger.Trigger, t *trigger.Trigger) []*trigger.Trigger {
	return append(list, t)
}

func (c *Client) BaseActionCmd(t *trigger.Trigger) {
	c.showtriggers(c.actions, "actions")
}

func (c *Client) AddActionCmd(t *trigger.Trigger) {
	c.AddActionString(t.Matches[1], t.Matches[2])
}

func (c *Client) UnactionCmd(t *trigger.Trigger) {
	c.actions = c.untrigger(c.actions, "action", t.Matches[1])
}

func (c *Client) BaseAliasCmd(t *trigger.Trigger) {
	c.showtriggers(c.aliases, "aliases")
}

func (c *Client) AddAliasCmd(t *trigger.Trigger) {
	c.AddAliasString(t.Matches[1], t.Matches[2])
}

func (c *Client) UnaliasCmd(t *trigger.Trigger) {
	c.aliases = c.untrigger(c.aliases, "alias", t.Matches[1])
}

func (c *Client) showtriggers(t []*trigger.Trigger, ttype string) {
	c.Print("## Current " + ttype + ":\n")
	for i, a := range t {
		c.Print(fmt.Sprintf("\n[%d]: %s", i, a.String()))
	}
	c.Print("\n")
}

func (c *Client) untrigger(triggerList []*trigger.Trigger, triggerType string, index string) []*trigger.Trigger {
	n, err := strconv.Atoi(index)
	if err != nil {
		c.Print(fmt.Sprintf("Invalid %s number: %d\n", triggerType, n))
		return triggerList
	}
	if n >= len(c.actions) {
		c.Print(fmt.Sprintf("%s not found: %d\n", triggerType, n))
		return triggerList
	}
	return append(triggerList[:n], triggerList[n+1:]...)
}

func (c *Client) GetWindowSize(name string) (int, int) {
	if win, ok := c.tui.Windows[name]; !ok {
		c.Print("\nNo window named '" + name + "' found.\n")
		return 0, 0
	} else {
		_, _, w, h := win.GetInnerRect()
		return w, h
	}
}
