package client

import (
	"io"
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
		w := tview.ANSIWriter(c)
		if _, err := io.Copy(w, c.conn); err != nil {
			c.ShowMain("Connection closed: " + err.Error() + "\n")
			c.conn = nil
		}
	}()
}

func (c *Client) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\r")
	for _, line := range lines {
		c.handleLine(line)
	}
	return len(p), nil
}

func (c *Client) handleLine(line string) {
	c.CurrentRaw = line
	c.CurrentText = stripTags(c.CurrentRaw)
	c.CheckTriggers(c.actions, strings.TrimSuffix(c.CurrentText, "\n"))
	if !c.Gag {
		c.ShowMain(c.CurrentRaw)
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
		c.ShowOverhead(c.CurrentRaw)
		c.Gag = true
	} else {
		c.ShowChat(c.CurrentRaw)
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
