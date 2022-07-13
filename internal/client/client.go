package client

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/seandheath/go-mud-client/internal/tui"
	"github.com/seandheath/go-mud-client/pkg/plugin"
	"github.com/seandheath/go-mud-client/pkg/trigger"
)

const BUFFERSIZE = 1024

type Client struct {
	conn      net.Conn
	buffer    []byte
	Gag       bool
	RawLine   string
	TextLine  string
	LogError  *log.Logger
	LogInfo   *log.Logger
	actions   []trigger.Trigger
	aliases   []trigger.Trigger
	functions map[string]trigger.Func
	plugins   map[string]*plugin.Config
	tui       *tui.TUI
}

func NewClient() *Client {
	c := &Client{}
	c.conn = nil
	c.buffer = make([]byte, BUFFERSIZE)
	c.Gag = false
	c.RawLine = "raw"
	c.TextLine = "text"
	c.LogError = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	c.LogInfo = log.New(os.Stderr, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	c.actions = []trigger.Trigger{}
	c.aliases = []trigger.Trigger{}
	c.functions = map[string]trigger.Func{}
	c.plugins = map[string]*plugin.Config{}
	c.tui = tui.NewTUI()
	c.tui.Parse = c.Parse
	c.cmdInit()
	return c
}

func (c *Client) AddAction(rs string, cmd trigger.Func) { c.actions = c.addTrigger(c.actions, rs, cmd) }
func (c *Client) AddActionString(rs string, cmd string) {
	c.actions = c.addTriggerString(c.actions, rs, cmd)
}
func (c *Client) AddAlias(rs string, cmd trigger.Func) { c.aliases = c.addTrigger(c.aliases, rs, cmd) }
func (c *Client) AddAliasString(rs string, cmd string) {
	c.aliases = c.addTriggerString(c.aliases, rs, cmd)
}

// Parse the string and send the result to the server
func (c *Client) Parse(text string) {
	if c.CheckTriggers(c.aliases, text) { // Check for aliases / commands
		return
	} else if c.conn == nil { // Not connected yet
		c.ShowMain("Not connected.\n")
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
	c.ShowMain(text + "\n")
	_, err := c.conn.Write([]byte(text + "\n"))
	if err != nil {
		c.ShowMain("Error sending: " + err.Error() + "\n")
		c.conn = nil
	}
}

func (c *Client) Run() {
	c.tui.Run()
}

// AddFunction maps a string to a function so that you can call the function
// from the mud with #function <name>
func (c *Client) AddFunction(name string, f func(t *trigger.Match)) {
	c.functions[name] = f
}

type showText struct {
	window string
	text   string
}

func (c *Client) ShowMain(text string) {
	c.Show("main", text)
}

func (c *Client) Show(window string, text string) {
	c.tui.Show(window, text)
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
	for n, win := range p.Windows {
		c.tui.AddWindow(n, win)
	}

	if len(p.Windows) > 0 {
		c.tui.FixInputLine(p.Grid.Rows, p.Grid.Columns)
	}
	c.plugins[name] = p
}

func (c *Client) CheckTriggers(list []trigger.Trigger, text string) bool {
	matched := false
	for _, t := range list {
		m := t.Re.FindStringSubmatch(text)
		if len(m) > 0 {
			matched = true
			t.Cmd(&trigger.Match{m, &t})
		}
	}
	return matched
}

func (c *Client) addTriggerString(list []trigger.Trigger, rs string, cmd string) []trigger.Trigger {
	f := func(t *trigger.Match) {
		c.Parse(cmd)
	}
	return c.addTrigger(list, rs, f)
}

func (c *Client) getFunc(cmd string) trigger.Func {

	return nil
}

// This function adds a trigger to the provided list and returns it
func (c *Client) addTrigger(list []trigger.Trigger, rs string, cmd trigger.Func) []trigger.Trigger {
	re, err := regexp.Compile(rs)
	if err != nil {
		c.ShowMain("Error compiling trigger: " + err.Error() + "\n")
		return list
	}
	return append(list, trigger.Trigger{re, cmd})
}

func (c *Client) BaseActionCmd(t *trigger.Match) {
	c.showtriggers(c.actions, "actions")
}

func (c *Client) AddActionCmd(t *trigger.Match) {
	c.AddActionString(t.Matches[1], t.Matches[2])
}

func (c *Client) UnactionCmd(t *trigger.Match) {
	c.actions = c.untrigger(c.actions, "action", t.Matches[1])
}

func (c *Client) BaseAliasCmd(t *trigger.Match) {
	c.showtriggers(c.aliases, "aliases")
}

func (c *Client) AddAliasCmd(t *trigger.Match) {
	c.AddAliasString(t.Matches[1], t.Matches[2])
}

func (c *Client) UnaliasCmd(t *trigger.Match) {
	c.aliases = c.untrigger(c.aliases, "alias", t.Matches[1])
}

func (c *Client) showtriggers(t []trigger.Trigger, ttype string) {
	c.ShowMain("## Current " + ttype + ":\n")
	for i, a := range t {
		c.ShowMain(fmt.Sprintf("\n[%d]: %s", i, a.Re.String()))
	}
	c.ShowMain("\n")
}

func (c *Client) untrigger(triggerList []trigger.Trigger, triggerType string, index string) []trigger.Trigger {
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
