package client

import (
	"log"
	"net"
	"os"
	"regexp"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type TriggerFunc func(*regexp.Regexp, []string)
type Module interface {
	Load()
}
type Trigger struct {
	Re  *regexp.Regexp
	Cmd TriggerFunc
}

var (
	myUI        *tea.Program
	Model       model
	Server      string
	Conn        net.Conn
	modules     map[string]Module
	actions     []Trigger
	aliases     []Trigger
	fmap        map[string]TriggerFunc
	CurrentRaw  string
	CurrentText string
	Gag         bool
	LogError    *log.Logger
	LogInfo     *log.Logger
	stats       runtime.MemStats
)

func init() {
	Server = ""
	Conn = nil
	modules = make(map[string]Module)
	actions = make([]Trigger, 0)
	aliases = make([]Trigger, 0)
	fmap = make(map[string]TriggerFunc)
	CurrentRaw = ""
	CurrentText = ""
	Gag = false
	LogError = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogInfo = log.New(os.Stderr, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Model = initialModel()
	CmdInit()
}

// Parse the string and send the result to the server
func Parse(text string) {
	if CheckTriggers(aliases, text) { // Check for aliases / commands
		return
	} else if Conn == nil { // Not connected yet
		//Show("Not connected.\n")
		Show("Not connected.\n")
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
	Show(text)
	//Show(text + "\n")
	_, err := Conn.Write([]byte(text + "\n"))
	if err != nil {
		//Show("Error sending: " + err.Error() + "\n")
		Show("Error sending: " + err.Error() + "\n")
		Conn = nil
	}
}

func LaunchUI() {
	myUI = tea.NewProgram(Model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if err := myUI.Start(); err != nil {
		LogError.Fatal("Error starting program: ", err)
		os.Exit(1)
	}
}

func AddAction(rs string, cmd TriggerFunc)  { actions = addTrigger(actions, rs, cmd) }
func AddActionString(rs string, cmd string) { actions = addTriggerString(actions, rs, cmd) }
func AddAlias(rs string, cmd TriggerFunc)   { aliases = addTrigger(aliases, rs, cmd) }
func AddAliasString(rs string, cmd string)  { aliases = addTriggerString(aliases, rs, cmd) }

func addTriggerString(list []Trigger, rs string, cmd string) []Trigger {
	f := func(*regexp.Regexp, []string) {
		Parse(cmd)
	}
	return addTrigger(list, rs, f)
}

// This function adds a trigger to the provided list and returns it
func addTrigger(list []Trigger, rs string, cmd TriggerFunc) []Trigger {
	re, err := regexp.Compile(rs)
	if err != nil {
		Show("Error compiling trigger: " + err.Error() + "\n")
		return list
	}
	return append(list, Trigger{re, cmd})
}

func CheckTriggers(list []Trigger, text string) bool {
	matched := false
	for _, t := range list {
		m := t.Re.FindStringSubmatch(text)
		if len(m) > 0 {
			matched = true
			t.Cmd(t.Re, m)
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
func RegisterFunction(name string, f func(*regexp.Regexp, []string)) {
	fmap[name] = f
}

func Show(text string) {
	myUI.Send(text)
}
