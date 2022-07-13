package trigger

import (
	"regexp"
)

type Match struct {
	Matches []string
	*Trigger
}

type Func func(*Match)
type Trigger struct {
	*regexp.Regexp
	Cmd Func
}
