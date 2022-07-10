package client

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/seandheath/go-mud-client/internal/tui"
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
	actions   []Trigger
	aliases   []Trigger
	functions map[string]TriggerFunc
	plugins   map[string]*PluginConfig
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
	c.actions = []Trigger{}
	c.aliases = []Trigger{}
	c.functions = map[string]TriggerFunc{}
	c.plugins = map[string]*PluginConfig{}
	c.tui = tui.NewTUI()
	c.tui.Parse = c.Parse
	c.cmdInit()
	return c
}

func (c *Client) AddAction(rs string, cmd TriggerFunc) { c.actions = c.addTrigger(c.actions, rs, cmd) }
func (c *Client) AddActionString(rs string, cmd string) {
	c.actions = c.addTriggerString(c.actions, rs, cmd)
}
func (c *Client) AddAlias(rs string, cmd TriggerFunc) { c.aliases = c.addTrigger(c.aliases, rs, cmd) }
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
func (c *Client) AddFunction(name string, f func(t *TriggerMatch)) {
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
