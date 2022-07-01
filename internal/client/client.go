package client

import (
	"log"
	"net"
	"regexp"
	"strings"
)

type Module interface {
	Load()
}
type trigger struct {
	re  *regexp.Regexp
	cmd func([]string)
}

var (
	Server      string
	Conn        net.Conn
	modules     map[string]Module
	actions     []trigger
	aliases     []trigger
	fmap        map[string]func([]string)
	CurrentRaw  string
	CurrentText string
	Gag         bool
)

var (
	colorPattern  = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([lbidrus]+|\-)?)?)?\]`)
	escapePattern = regexp.MustCompile(`\[([a-zA-Z0-9_,;: \-\."#]+)\[(\[*)\]`)
)

func init() {
	Server = ""
	Conn = nil
	modules = make(map[string]Module)
	actions = make([]trigger, 0)
	aliases = make([]trigger, 0)
	fmap = make(map[string]func([]string))
	CurrentRaw = ""
	CurrentText = ""
	Gag = false
	CmdInit()
}

// Parse the string and send the result to the server
func Parse(text string) {
	if CheckTriggers(aliases, text) {
		return
	}
	if Conn == nil {
		ShowMain("Not connected.\n")
		return
	} else {
		SendNow(text)
	}

}

func SendNow(text string) {
	ShowMain("\n" + text + "\n")
	_, err := Conn.Write([]byte(text + "\n"))
	if err != nil {
		ShowMain("Error sending: " + err.Error() + "\n")
		Conn = nil
	}
}

func AddAction(rs string, cmd interface{}) { actions = addTrigger(actions, rs, cmd) }
func AddAlias(rs string, cmd interface{})  { aliases = addTrigger(aliases, rs, cmd) }

// This function adds a trigger to the provided list and returns it
func addTrigger(list []trigger, rs string, cmd interface{}) []trigger {
	var f func([]string)
	switch ct := cmd.(type) {
	case string:
		if strings.HasPrefix(ct, "#func") { // Check for a function call - saves parsing later and directly registers the function now
			tok := strings.Split(ct, " ")
			if len(tok) != 2 { // Should show #function <name> and <name> should be registered
				log.Fatal("Error: Invalid function call: " + ct + "\n")
			} else {
				if _, ok := fmap[tok[1]]; !ok { // The function isn't registered yet
					log.Fatal("Error: Function not found in function table, did you register it? : " + ct + "\n")
				} else {
					f = fmap[tok[1]]
				}
			}
		} else { // Just sending some text to the mud
			f = func([]string) {
				Parse(ct)
			}
		}
	case func([]string):
		f = ct // Calls a function
	default:
		log.Fatal("Error: Invalid trigger type for match: " + rs + "\n")
	}

	re, err := regexp.Compile(rs)
	if err != nil {
		ShowMain("Error compiling trigger: " + err.Error() + "\n")
		return list
	}
	return append(list, trigger{re, f})
}

func CheckTriggers(list []trigger, text string) bool {
	matched := false
	for _, a := range list {
		m := a.re.FindStringSubmatch(text)
		if len(m) > 0 {
			matched = true
			a.cmd(m)
		}
	}
	return matched
}

func LoadModule(name string, m Module) {
	if _, ok := modules[name]; !ok {
		modules[name] = m
	}
	modules[name].Load()
}

// RegisterFunction maps a string to a function so that you can call the function
// from the mud with #function <name>
func RegisterFunction(name string, f func([]string)) {
	fmap[name] = f
}
