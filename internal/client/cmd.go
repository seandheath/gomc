package client

import (
	"bufio"
	"net"
	"strings"

	"github.com/rivo/tview"
)

func (c *Client) CmdInit() {
	c.AddAlias("^#connect .*$", c.connect)
	c.AddAlias("^#capture.*$", c.Capture)
	c.AddAlias("^#func .*$", c.ExecuteFunction)
}

// connect takes a string from the user and attempts to connect to the mud server.
// If the connection is successful then a goroutine is launched to handle the connection.
func (c *Client) connect(text string) {
	if c.conn != nil {
		c.ShowMain("Already connected.\n")
		return
	}
	text = strings.TrimPrefix(text, "#connect ")
	conn, err := net.Dial("tcp", text)
	if err != nil {
		c.ShowMain("Failed to connect: " + err.Error() + "\n")
	}
	c.conn = conn
	go func() {
		defer c.conn.Close()
		r := bufio.NewReader(c.conn)
		for {
			line, err := r.ReadBytes('\r')
			if err != nil {
				c.ShowMain("Connection closed.\n")
				c.conn = nil
				return
			}
			c.handleLine(line[:len(line)-2]) // Removes the \n\r at the end of each line
		}
	}()
}

func (c *Client) handleLine(bytes []byte) {
	line := string(bytes)
	c.CurrentRaw = tview.TranslateANSI(line)
	c.CurrentText = stripTags(c.CurrentRaw)
	c.CheckTriggers(c.actions, c.CurrentText)
	if !c.Gag {
		c.ShowMain(c.CurrentRaw + "\n")
	} else {
		c.Gag = false
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

func (c *Client) Capture(text string) {
	s := strings.TrimPrefix(text, "#capture ")

	if s == "overhead" {
		c.ShowOverhead(c.CurrentRaw + "\n")
		c.Gag = true
	} else {
		c.ShowChat(c.CurrentRaw + "\n")
	}
}

func (c *Client) ExecuteFunction(text string) {
	s := strings.TrimPrefix(text, "#func ")
	if f, ok := c.fmap[s]; ok {
		f(c.CurrentText)
	} else {
		c.ShowMain("Function not found:" + text + "\n")
	}
}
