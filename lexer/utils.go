package lexer

import "unicode"

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isEOL reports whether r is an end of line (EOL)
// TODO: check for "\r\n"
func isEOL(r rune) bool {
	return r == '\n'
}
