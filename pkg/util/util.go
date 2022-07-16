package util

// protect triggers by removing delimiters (right now hard-coded as semicolon)
func SwapSemi(b []byte) []byte {
	for i, bt := range b {
		if bt == ';' {
			b[i] = ':'
		}
	}
	return b
}

// Remove carriage returns and newlines
func TrimEnd(b []byte) []byte {
	// If the last character is a newline, remove it and recurse
	if len(b) > 0 {
		if b[len(b)-1] == '\r' || b[len(b)-1] == '\n' {
			if len(b) > 1 {
				// Remove the newline
				b = b[:len(b)-1]
			} else {
				// The string only had newlines
				return []byte{}
			}
			// Check again
			return TrimEnd(b)
		}
	}
	// No newlines, return the string
	return b
}
