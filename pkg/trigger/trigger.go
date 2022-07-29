package trigger

import (
	"regexp"
)

type Func func(*Trigger)
type Trigger struct {
	*regexp.Regexp

	Matches  []string
	Results  map[string]string
	cmd      Func
	counting bool // If counting is enabled then the trigger will execute Count times and then disable
	count    int
	Enabled  bool
}

func NewTrigger(re string, cmd Func) *Trigger {
	t := &Trigger{}
	t.Regexp = regexp.MustCompile(re)
	t.cmd = cmd
	t.Matches = []string{}
	t.Results = map[string]string{}
	t.Enabled = true
	return t
}

func (t *Trigger) Do() {
	if t.counting {
		if t.count > 0 {
			t.count--
		} else {
			t.Enabled = false
		}
	}
	if t.Enabled {
		t.cmd(t)
	}
}

func (t *Trigger) SetCount(c int) {
	t.counting = true
	t.Enabled = true
	t.count = c
}

func (t *Trigger) Increment() {
	t.counting = true
	t.Enabled = true
	t.count++
}
