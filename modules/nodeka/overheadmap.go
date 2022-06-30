package nodeka

var inMap = false

// Handle seeing an overhead map line
func MapLine(text string) {
	inMap = true
	Client.ShowOverhead(Client.CurrentRaw)
	Client.Gag = true
}

func EmptyLine(text string) {
	if inMap {
		inMap = false
		Client.Gag = true
	}
}
