package client

import (
	"log"
	"net"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	uiProgram   *tea.Program
	uiModel     *model
	Conn        net.Conn
	CurrentRaw  string
	CurrentText string
	Gag         bool
	LogError    *log.Logger
	LogInfo     *log.Logger
	actions     []Trigger
	aliases     []Trigger
	functions   map[string]TriggerFunc
	plugins     map[string]*PluginConfig
)

func init() {
	Conn = nil
	CurrentRaw = ""
	CurrentText = ""
	Gag = false
	LogError = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogInfo = log.New(os.Stderr, "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	actions = []Trigger{}
	aliases = []Trigger{}
	functions = map[string]TriggerFunc{}
	plugins = map[string]*PluginConfig{}
	uiModel = initialModel()
	cmdInit()
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

func Run() {
	uiProgram = tea.NewProgram(
		uiModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if err := uiProgram.Start(); err != nil {
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
		ShowMain("Error compiling trigger: " + err.Error() + "\n")
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

func LoadPlugin(name string, p *PluginConfig) {
	for re, cmd := range p.Actions {
		AddActionString(re, cmd)
	}
	for re, cmd := range p.Aliases {
		AddAliasString(re, cmd)
	}
	for n, f := range p.Functions {
		AddFunction(n, f)
	}
	plugins[name] = p
}

// AddFunction maps a string to a function so that you can call the function
// from the mud with #function <name>
func AddFunction(name string, f func(*regexp.Regexp, []string)) {
	functions[name] = f
}

func AddWindow(name string, width int, height int) {
	uiModel.AddWindow(name, width, height)
}

func SetView(f func(ws map[string]*Window) string) error {
	uiModel.ViewFunc = f
	return nil
}

func SetResize(f func(int, int, map[string]*Window) map[string]*Window) {
	uiModel.ResizeFunc = f
}

type showText struct {
	window string
	text   string
}

func ShowMain(text string) {
	Show("main", text)
}
func Show(window string, text string) {
	uiProgram.Send(showText{window, text})
}
