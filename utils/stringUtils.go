package utils

import (
	"bytes"
	"unicode"
)

func IsDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

// ParseCommand parses a command line and handle arguments in quotes.
// https://github.com/vrischmann/shlex/blob/master/shlex.go
func ParseCommand(s string) (res []string) {
	var buf bytes.Buffer
	insideQuotes := false
	for _, r := range s {
		switch {
		case unicode.IsSpace(r) && !insideQuotes:
			if buf.Len() > 0 {
				res = append(res, buf.String())
				buf.Reset()
			}
		case r == '"' || r == '\'':
			if insideQuotes {
				res = append(res, buf.String())
				buf.Reset()
				insideQuotes = false
				continue
			}
			insideQuotes = true
		default:
			buf.WriteRune(r)
		}
	}
	if buf.Len() > 0 {
		res = append(res, buf.String())
	}
	return
}
