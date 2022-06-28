package client

import (
	"net"
	"regexp"
)

type Module interface {
	Load(*Client)
}
type trigger struct {
	re  *regexp.Regexp
	cmd func(string)
}

type Client struct {
	server      string
	conn        net.Conn
	modules     map[string]Module
	actions     []trigger
	aliases     []trigger
	CurrentRaw  string
	CurrentText string
	Gag         bool
}

var (
	colorPattern  = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([lbidrus]+|\-)?)?)?\]`)
	escapePattern = regexp.MustCompile(`\[([a-zA-Z0-9_,;: \-\."#]+)\[(\[*)\]`)
)

func NewClient() *Client {
	c := &Client{
		server:      "",
		conn:        nil,
		modules:     make(map[string]Module),
		actions:     make([]trigger, 0),
		aliases:     make([]trigger, 0),
		CurrentRaw:  "",
		CurrentText: "",
		Gag:         false,
	}
	c.CmdInit()
	return c
}

func (c *Client) Run() {
	c.LaunchUI()
}

// Parse the string and send the result to the server
func (c *Client) Parse(text string) {
	if c.CheckTriggers(c.aliases, text) {
		return
	}
	if c.conn == nil {
		c.ShowMain("Not connected.\n")
		return
	} else {
		c.SendNow(text)
	}

}

func (c *Client) SendNow(text string) {
	c.ShowMain("\n" + text + "\n")
	_, err := c.conn.Write([]byte(text + "\n"))
	if err != nil {
		c.ShowMain("Error sending: " + err.Error() + "\n")
		c.conn = nil
	}
}

func (c *Client) AddAction(rs string, cmd interface{}) { c.actions = c.AddTrigger(c.actions, rs, cmd) }
func (c *Client) AddAlias(rs string, cmd interface{})  { c.aliases = c.AddTrigger(c.aliases, rs, cmd) }

// This function adds a trigger to the provided list and returns it
func (c *Client) AddTrigger(list []trigger, rs string, cmd interface{}) []trigger {
	var f func(string)
	switch cmd := cmd.(type) {
	case string:
		f = func(string) { c.Parse(cmd) }
	case func(string):
		f = cmd
	}

	re, err := regexp.Compile(rs)
	if err != nil {
		c.ShowMain("Error compiling trigger: " + err.Error() + "\n")
		return list
	}
	return append(list, trigger{re, f})
}

func (c *Client) CheckTriggers(list []trigger, text string) bool {
	matched := false
	for _, a := range list {
		if a.re.MatchString(text) {
			matched = true
			a.cmd(text)
		}
	}
	return matched
}

func (c *Client) LoadModule(name string, m Module) {
	if _, ok := c.modules[name]; !ok {
		c.modules[name] = m
	}
	c.modules[name].Load(c)
}
