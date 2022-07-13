package trigger

import (
	"regexp"
)

type Func func(*Trigger)
type Trigger struct {
	*regexp.Regexp
	Cmd     Func
	Matches []string
	Results map[string]string
}
