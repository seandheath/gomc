package client

import (
	"bufio"
	"bytes"
	"net"
	"regexp"
	"strings"
	"time"
)

var lastRead time.Time = time.Now()

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var ansiRegexp = regexp.MustCompile(ansi)

// ConnectCmd takes a string from the user and attempts to ConnectCmd to the mud server.
// If the connection is successful then a goroutine is launched to handle the connection.
func (c *Client) ConnectCmd(t *TriggerMatch) {
	if c.Conn != nil {
		c.ShowMain("Already connected.\n")
		return
	}
	text := t.Matches[1]
	conn, err := net.Dial("tcp", text)
	if err != nil {
		c.ShowMain("Failed to connect: " + err.Error() + "\n")
	}
	c.Conn = conn
	go func() {
		defer c.Conn.Close()
		//buffer := make([]byte, bufio.MaxScanTokenSize)
		scanner := bufio.NewScanner(c.Conn)
		scanner.Split(split)
		//scanner.Buffer(buffer, bufio.MaxScanTokenSize)
		for scanner.Scan() {
			c.RawLine = scanner.Text()
			c.TextLine = strip(c.RawLine)
			c.CheckTriggers(c.actions, strings.TrimSuffix(c.TextLine, "\n"))
			if !c.Gag {
				c.ShowMain(c.RawLine)
			} else {
				c.Gag = false
			}
		}
	}()
}

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if len(data) > 0 {
		if i := bytes.IndexByte(data, '\r'); i >= 0 {
			lastRead = time.Now()
			return i + 1, data[0:i], nil
		}
		if timeout := lastRead.Add(time.Microsecond * 100); timeout.After(time.Now()) {
			lastRead = time.Now()
			return len(data), data, nil
		}
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func strip(str string) string {
	return ansiRegexp.ReplaceAllString(str, "")
}
