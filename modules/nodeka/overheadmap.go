package nodeka

// Handle seeing an overhead map line
func (m *Module) MapLine(text string) {
	m.Client.ShowOverhead(text)
	m.Client.CurrentRaw = "" // Gag
}
