package nodeka

import "github.com/seandheath/go-mud-client/internal/client"

var inMap = false

// Handle seeing an overhead map line
func MapLine() {
	inMap = true
	client.ShowOverhead(client.CurrentRaw)
	client.Gag = true
}

// Handle seeing an empty line
func EmptyLine() {
	if inMap {
		inMap = false
		client.Gag = true
	}
}
