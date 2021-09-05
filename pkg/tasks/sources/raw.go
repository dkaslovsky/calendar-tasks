package sources

import (
	"fmt"
	"strings"
)

const delim = ":"

type RawLine struct {
	Date string
	Text string
}

func LoadLine(line string) (*RawLine, error) {
	parts := strings.SplitN(line, delim, 2)
	if len(parts) != 2 {
		return &RawLine{}, fmt.Errorf("invalid line [%s]", line)
	}

	r := &RawLine{
		Date: parts[0],
		Text: cleanText(parts[1]),
	}
	return r, nil
}

func cleanText(s string) string {
	return strings.TrimLeft(s, " ")
}
