package nodeka

import (
	"regexp"

	"github.com/seandheath/go-mud-client/internal/client"
)

var inMap = false

// Handle seeing an overhead map line
func MapLine(re *regexp.Regexp, matches []string) {
	inMap = true
	client.Show(client.CurrentRaw)
	client.Gag = true
}

// Handle seeing an empty line
func EmptyLine(re *regexp.Regexp, matches []string) {
	if inMap {
		inMap = false
		client.Gag = true
	}
}
