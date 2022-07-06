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
var ConnectCmd TriggerFunc = func(re *regexp.Regexp, matches []string) {
	if Conn != nil {
		Show("Already connected.\n")
		return
	}
	text := matches[1]
	conn, err := net.Dial("tcp", text)
	if err != nil {
		Show("Failed to connect: " + err.Error() + "\n")
	}
	Conn = conn
	go func() {
		defer Conn.Close()
		buffer := make([]byte, bufio.MaxScanTokenSize)
		scanner := bufio.NewScanner(Conn)
		scanner.Split(split)
		scanner.Buffer(buffer, bufio.MaxScanTokenSize)
		for scanner.Scan() {
			CurrentRaw := scanner.Text()
			CurrentText := strip(CurrentRaw)
			CheckTriggers(actions, strings.TrimSuffix(CurrentText, "\n"))
			Show(CurrentRaw)
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
		if timeout := lastRead.Add(time.Millisecond * 100); timeout.After(time.Now()) {
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
