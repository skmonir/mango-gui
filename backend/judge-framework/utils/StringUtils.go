package utils

import (
	"bufio"
	"bytes"
	"math/rand"
	"strings"
	"unicode"
)

const CHA = "abcdefghijklmnopqrstuvwxyz0123456789"

// RandString n is the length. a-z 0-9
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CHA[rand.Intn(len(CHA))]
	}
	return string(b)
}

func IsDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func TrimIO(io string) string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(io))
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	io = strings.Join(lines, "\n")
	return io
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
