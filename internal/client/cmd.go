package client

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/seandheath/gomc/pkg/trigger"
)

func (c *Client) cmdInit() {
	c.AddAliasFunc("^#connect (.*)$", c.ConnectCmd)
	c.AddAliasFunc("^#capture ?(.*)$", c.CaptureCmd)
	c.AddAliasFunc("^#func (.*)$", c.FuncCmd)
	c.AddAliasFunc("^#match (.+)$", c.MatchCmd)
	c.AddAliasFunc("^#(\\d+) (.+)$", c.LoopCmd)
	c.AddAliasFunc("^#action$", c.BaseActionCmd)
	c.AddAliasFunc("^#action {(.+)}{(.+)}$", c.AddActionCmd)
	c.AddAliasFunc("^#unaction (\\d+)$", c.UnactionCmd)
	c.AddAliasFunc("^#alias$", c.BaseAliasCmd)
	c.AddAliasFunc("^#alias {(.+)}{(.+)}$", c.AddAliasCmd)
	c.AddAliasFunc("^#unalias (\\d+)$", c.UnaliasCmd)
	c.AddAliasFunc("^#memstats$", c.MemStatsCmd)
}

func (c *Client) CaptureCmd(t *trigger.Trigger) {
	s := strings.TrimPrefix(t.Matches[0], "#capture ")

	if s == "overhead" {
		c.PrintTo("omap", string(c.RawLine))
		c.Gag = true
	} else {
		ts := time.Now().Format("[2006:01:02 15:04:05] ")
		c.PrintTo("chat", ts+string(c.RawLine))
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
	cmd := t.Matches[2]
	if err != nil {
		c.Print("Error parsing loop number: " + err.Error() + "\n")
	}
	for i := 0; i < n; i++ {
		c.Parse(cmd)
	}
}

func (c *Client) MemStatsCmd(t *trigger.Trigger) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	c.Print(fmt.Sprintf("Alloc: %d MiB", m.Alloc/1024/1024))
}
