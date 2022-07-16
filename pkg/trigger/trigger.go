package trigger

import (
	"regexp"
)

type Func func(*Trigger)
type Trigger struct {
	*regexp.Regexp
	Cmd      Func
	Matches  []string
	Results  map[string]string
	counting bool
	Count    int
	Enabled  bool
}

type Queue struct {
	*Trigger
	queue []func()
}

func NewTrigger(re string, cmd Func) *Trigger {
	return &Trigger{
		Regexp:  regexp.MustCompile(re),
		Cmd:     cmd,
		Matches: []string{},
		Results: map[string]string{},
		Enabled: true,
	}
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
