package client

import (
	"io"
	"net"
	"strings"

	"github.com/rivo/tview"
)

func (c *Client) CmdInit() {
	c.AddAlias("^#connect .*$", c.connect)
	c.AddAlias("^#capture.*$", c.capture)
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
			return
		}
	}()
}

func (c *Client) capture(text string) {
	s := strings.TrimPrefix(text, "#capture ")

	if s == "overhead" {
		c.ShowOverhead(c.CurrentRaw)
	} else {
		c.ShowChat(c.CurrentRaw)
	}
}
