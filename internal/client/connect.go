package client

import (
	"bufio"
	"errors"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/seandheath/go-mud-client/pkg/trigger"
)

const (
	DEADLINE = time.Millisecond * 5
)

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var ansiRegexp = regexp.MustCompile(ansi)

// ConnectCmd takes a string from the user and attempts to ConnectCmd to the mud server.
// If the connection is successful then a goroutine is launched to handle the connection.
func (c *Client) ConnectCmd(t *trigger.Match) {
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
	s := bufio.NewScanner(c.conn)
	s.Split(bufio.ScanLines)
	go func() {
		defer c.conn.Close()
		newData := false
		newLine := ""
		last := ""
		for {
			lines, err := c.readLines(last)
			last = ""
			if errors.Is(err, os.ErrDeadlineExceeded) {
				// Deadline exceeded, we have all the data
				// Reset the deadline
				c.conn.SetReadDeadline(time.Time{})
				newData = false
			} else if err != nil {
				// Some other connection error
				c.conn = nil
				c.Print("Disconnected: " + err.Error())
				return
			} else {
				// New data and deadline not exceeded
				newData = true
			}

			// We have lines
			if lines != nil {
				// We got new data, no timeout yet, and we don't know if it's all the data
				// We'll check the last string for a newline and pass it back to readlines to
				// prepend the next round of data
				if newData && !strings.HasSuffix(lines[len(lines)-1], "\n") {
					last = lines[len(lines)-1]
					lines = lines[:len(lines)-1]
				}

				// Reset our new line to print and handle all the lines
				newLine = ""
				for _, line := range lines {
					newLine += c.handleLine(line)
				}
				if newLine != "" {
					c.Print(newLine)
				}
			}

			// We got data, set deadline and read again
			if newData {
				c.conn.SetReadDeadline(time.Now().Add(DEADLINE))
			}
		}
	}()
}

func (c *Client) handleLine(line string) string {
	c.RawLine = line
	c.TextLine = strip(c.RawLine)
	c.CheckTriggers(c.actions, strings.TrimSuffix(c.TextLine, "\n"))
	if c.Gag {
		c.Gag = false
		return ""
	}
	return c.RawLine
}
func (c *Client) readLines(last string) ([]string, error) {

	n, err := c.conn.Read(c.buffer)

	// If something goes wrong kick it back
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return nil, err
	}
	lines := strings.Split(string(c.buffer[:n]), "\r")
	lines[0] = last + lines[0]

	// No data
	return lines, err
}

func strip(str string) string {
	return ansiRegexp.ReplaceAllString(str, "")
}
