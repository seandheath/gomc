package client

import (
	"log"
	"net"
	"regexp"
	"strings"
)

type matcher struct {
	re  *regexp.Regexp
	cmd string
}

var (
	Actions    []matcher
	Aliases    []matcher
	Variables  map[string]string
	Connection net.Conn
	Server     string
)

func Run() {
	Variables = make(map[string]string)
	Actions = make([]matcher, 0)
	Aliases = make([]matcher, 0)
	Connection = nil
	Server = ""
	Launch()
}

// Parse the string and send the result to the server
func Parse(text string) {
	if strings.HasPrefix(text, "#") {
		cmd := strings.Split(text, " ")[0]
		if _, ok := Commands[cmd]; ok {
			Commands[cmd].Run(strings.TrimSpace(strings.TrimPrefix(text, cmd)))
			return
		} else {
			ShowMain("Unknown command: " + cmd + "\n")
			return
		}
	} else {
		for _, a := range Aliases {
			if a.re.MatchString(text) {
				text = a.cmd
				break
			}
		}
	}
	if Connection == nil {
		ShowMain("Not connected.\n")
		return
	} else {
		SendNow(text)
	}

}

func SendNow(text string) {
	_, err := Connection.Write([]byte(text + "\n"))
	if err != nil {
		log.Fatal("failed SendNow: ", err)
	}
}
