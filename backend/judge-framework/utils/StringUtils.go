package utils

import (
	"bufio"
	"strings"
)

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
