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

type Queue struct {
	*Trigger
	queue []func()
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

func NewQueue(re string) *Queue {
	q := &Queue{}
	q.Trigger = NewTrigger(re, q.Do)
	q.Enabled = false
	return q
}

func (q *Queue) Prepend(f func()) {
	q.Enabled = true
	q.queue = append([]func(){f}, q.queue...)
}
func (q *Queue) Append(f func()) {
	q.Enabled = true
	q.queue = append(q.queue, f)
}
func (q *Queue) Do(t *Trigger) {
	q.Enabled = false
	for _, f := range q.queue {
		f()
	}
	q.queue = []func(){}
}
