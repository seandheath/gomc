package client

import (
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

type Module interface {
	Load()
}
type Trigger struct {
	Re  *regexp.Regexp
	Cmd func()
}

var (
	Server         string
	Conn           net.Conn
	modules        map[string]Module
	actions        []Trigger
	aliases        []Trigger
	fmap           map[string]func()
	CurrentRaw     string
	CurrentText    string
	CurrentMatches []string
	CurrentTrigger Trigger
	Gag            bool
	LogError       *log.Logger
	LogInfo        *log.Logger
)

var (
	colorPattern  = regexp.MustCompile(`\[([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([a-zA-Z]+|#[0-9a-zA-Z]{6}|\-)?(:([lbidrus]+|\-)?)?)?\]`)
	escapePattern = regexp.MustCompile(`\[([a-zA-Z0-9_,;: \-\."#]+)\[(\[*)\]`)
)

func init() {
	Server = ""
	Conn = nil
	modules = make(map[string]Module)
	actions = make([]Trigger, 0)
	aliases = make([]Trigger, 0)
	fmap = make(map[string]func())
	CurrentRaw = ""
	CurrentText = ""
	Gag = false
	LogError = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogInfo = log.New(os.Stderr, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	CmdInit()
}

// Parse the string and send the result to the server
func Parse(text string) {
	if CheckTriggers(aliases, text) { // Check for aliases / commands
		return
	} else if Conn == nil { // Not connected yet
		ShowMain("Not connected.\n")
		return
	} else if strings.Contains(text, ";") { // Allow splitting commands by ;
		s := strings.Split(text, ";")
		for _, t := range s {
			Parse(t)
		}
	} else {
		SendNow(text)
	}
}

func SendNow(text string) {
	ShowMain(text + "\n")
	_, err := Conn.Write([]byte(text + "\n"))
	if err != nil {
		ShowMain("Error sending: " + err.Error() + "\n")
		Conn = nil
	}
}

func AddAction(rs string, cmd interface{}) { actions = addTrigger(actions, rs, cmd) }
func AddAlias(rs string, cmd interface{})  { aliases = addTrigger(aliases, rs, cmd) }

// This function adds a trigger to the provided list and returns it
func addTrigger(list []Trigger, rs string, cmd interface{}) []Trigger {
	var f func()
	switch ct := cmd.(type) {
	case string:
		if strings.HasPrefix(ct, "#func") { // Check for a function call - saves parsing later and directly registers the function now
			tok := strings.Split(ct, " ")
			if len(tok) != 2 { // Should show #function <name> and <name> should be registered
				LogError.Println("Invalid function call: " + ct + "\n")
			} else {
				if _, ok := fmap[tok[1]]; !ok { // The function isn't registered yet
					LogError.Println("Error: Function not found in function table, did you register it? : " + ct + "\n")
				} else {
					f = fmap[tok[1]]
				}
			}
		} else { // Just sending some text to the mud
			f = func() {
				Parse(ct)
			}
		}
	case func():
		f = ct // Calls a function
	default:
		LogError.Println("Error: Invalid trigger type for match: " + rs + "\n")
	}

	re, err := regexp.Compile(rs)
	if err != nil {
		ShowMain("Error compiling trigger: " + err.Error() + "\n")
		return list
	}
	return append(list, Trigger{re, f})
}

func CheckTriggers(list []Trigger, text string) bool {
	matched := false
	for _, t := range list {
		m := t.Re.FindStringSubmatch(text)
		if len(m) > 0 {
			matched = true
			CurrentTrigger = t
			CurrentMatches = m
			CurrentTrigger.Cmd()
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
func RegisterFunction(name string, f func()) {
	fmap[name] = f
}
