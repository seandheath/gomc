package trigger

import (
	"regexp"
)

type Match struct {
	Matches []string
	Trigger *Trigger
}

type Func func(*Match)
type Trigger struct {
	Re  *regexp.Regexp
	Cmd Func
}
