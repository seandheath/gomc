package client

import (
	"bytes"
	"errors"
	"io"
	"net"
	"os"
	"time"

	"code.rocketnine.space/tslocum/cview"
	"github.com/seandheath/gomc/pkg/trigger"
	"github.com/seandheath/gomc/pkg/util"
)

const (
	DEADLINE = time.Millisecond * 20
)

// ConnectCmd takes a string from the user and attempts to ConnectCmd to the mud server.
// If the connection is successful then a goroutine is launched to handle the connection.
func (c *Client) ConnectCmd(t *trigger.Trigger) {
	if c.conn != nil {
		c.Print("Already connected.\n")
		return
	}
	text := t.Matches[1]
	conn, err := net.Dial("tcp", text)
	if err != nil {
		c.Print("Failed to connect: " + err.Error() + "\n")
	}
	c.conn = conn
	w := cview.ANSIWriter(c)
	go func() {
		defer c.conn.Close()
		for {
			if _, err := io.Copy(w, c.conn); err != nil {
				if !errors.Is(err, os.ErrDeadlineExceeded) {
					c.Print("Connection closed: " + err.Error() + "\n")
					c.conn = nil
				} else {
					c.handleData(nil)
					c.conn.SetReadDeadline(time.Time{})
				}
			}
		}
	}()
}

func (c *Client) Write(b []byte) (int, error) {
	return c.handleData(b)
}
func (c *Client) handleData(b []byte) (int, error) {
	c.processBuffer = append(c.processBuffer, b...)

	// If the last byte isn't a newline, set a deadline to prevent blocking
	// and try and get some more data.
	if len(b) > 0 {
		if b[len(b)-1] != '\n' && b[len(b)-1] != '\r' {

			// If the timeout isn't currently set, we'll set it for DEADLINE from
			// now
			if c.timeout.IsZero() {
				c.timeout = time.Now().Add(DEADLINE)
				c.conn.SetReadDeadline(c.timeout)
			}

			// If we haven't hit the deadline yet then try and read some more data
			if time.Now().Before(c.timeout) {
				return len(b), nil
			} else {
				c.Print("\n\ntimeout\n\n")
			}

		}
	}

	// We either have a newline at the end or we hit the timeout, so we can
	// reset the deadline to just read until we get more data
	c.timeout = time.Time{}

	for _, line := range bytes.Split(c.processBuffer, []byte("\r")) {
		//tline := util.TrimEnd(line)
		c.RawLine = line
		sline := cview.StripTags(c.RawLine, true, true)
		c.TextLine = util.TrimEnd(util.SwapSemi(sline))
		ts := c.CheckTriggers(c.actions, string(c.TextLine))
		for _, t := range ts {
			t.Do()
		}
		if c.Gag {
			c.Gag = false
		} else {

			c.printBuffer = append(c.printBuffer, c.RawLine...)
		}
	}
	c.PrintBytesTo("main", c.printBuffer)
	c.processBuffer = nil
	c.printBuffer = nil
	return len(b), nil
}
