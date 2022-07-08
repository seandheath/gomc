package client

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func (c *Client) cmdInit() {
	c.AddAlias("^#connect (.*)$", c.ConnectCmd)
	c.AddAlias("^#capture ?(.*)$", c.CaptureCmd)
	c.AddAlias("^#func (.*)$", c.FuncCmd)
	c.AddAlias("^#match (.+)$", c.MatchCmd)
	c.AddAlias("^#(\\d+) (.+)$", c.LoopCmd)
	c.AddAlias("^#action$", c.BaseActionCmd)
	c.AddAlias("^#action {(.+)}{(.+)}$", c.AddActionCmd)
	c.AddAlias("^#unaction (\\d+)$", c.UnactionCmd)
	c.AddAlias("^#alias$", c.BaseAliasCmd)
	c.AddAlias("^#alias {(.+)}{(.+)}$", c.AddAliasCmd)
	c.AddAlias("^#unalias (\\d+)$", c.UnaliasCmd)
	c.AddAlias("^#memstats$", c.MemStatsCmd)
}

func (c *Client) CaptureCmd(t *TriggerMatch) {
	s := strings.TrimPrefix(t.Matches[0], "#capture ")

	if s == "overhead" {
		c.Show("omap", c.RawLine)
		c.Gag = true
	} else {
		ts := time.Now().Format("2006:01:02 15:04:05")
		c.Show("chat", fmt.Sprintf("[%s] %s\n", ts, strings.TrimSuffix(c.RawLine, "\n")))
	}
}

func (c *Client) FuncCmd(t *TriggerMatch) {
	s := strings.TrimPrefix(t.Matches[0], "#func ")
	if f, ok := c.functions[s]; ok { // Found the function
		f(t)
	} else {
		c.ShowMain("Function not found:" + t.Matches[0] + "\n")
	}
}

func (c *Client) MatchCmd(t *TriggerMatch) {
	c.CheckTriggers(c.actions, t.Matches[1])
}

// Nested triggers overwrite matches... need to pass into the function
func (c *Client) LoopCmd(t *TriggerMatch) {
	n, err := strconv.Atoi(t.Matches[1])
	if err != nil {
		c.ShowMain("Error parsing loop number: " + err.Error() + "\n")
	}
	for i := 0; i < n; i++ {
		c.Parse(t.Matches[2])
	}
}

func (c *Client) MemStatsCmd(t *TriggerMatch) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	c.ShowMain(fmt.Sprintf("Alloc: %d MiB", m.Alloc/1024/1024))
}
