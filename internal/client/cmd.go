package client

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/seandheath/go-mud-client/pkg/trigger"
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

func (c *Client) CaptureCmd(t *trigger.Trigger) {
	s := strings.TrimPrefix(t.Matches[0], "#capture ")

	if s == "overhead" {
		c.PrintTo("omap", c.RawLine)
		c.Gag = true
	} else {
		ts := time.Now().Format("2006:01:02 15:04:05")
		c.PrintTo("chat", fmt.Sprintf("[%s] %s\n", ts, strings.TrimSuffix(c.RawLine, "\n")))
	}
}

func (c *Client) FuncCmd(t *trigger.Trigger) {
	s := strings.TrimPrefix(t.Matches[0], "#func ")
	if f, ok := c.functions[s]; ok { // Found the function
		f(t)
	} else {
		c.Print("Function not found:" + t.Matches[0] + "\n")
	}
}

func (c *Client) MatchCmd(t *trigger.Trigger) {
	c.CheckTriggers(c.actions, t.Matches[1])
}

// Nested triggers overwrite matches... need to pass into the function
func (c *Client) LoopCmd(t *trigger.Trigger) {
	n, err := strconv.Atoi(t.Matches[1])
	if err != nil {
		c.Print("Error parsing loop number: " + err.Error() + "\n")
	}
	for i := 0; i < n; i++ {
		c.Parse(t.Matches[2])
	}
}

func (c *Client) MemStatsCmd(t *trigger.Trigger) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	c.Print(fmt.Sprintf("Alloc: %d MiB", m.Alloc/1024/1024))
}
