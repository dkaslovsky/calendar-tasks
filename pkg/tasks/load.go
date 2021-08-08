package tasks

import (
	"fmt"
	"strings"
)

const delim = ":"

type rawLine struct {
	date string
	text string
}

func loadLine(line string) (*rawLine, error) {
	parts := strings.Split(line, delim)
	if len(parts) != 2 {
		return &rawLine{}, fmt.Errorf("invalid line [%s]", line)
	}

	r := &rawLine{
		date: parts[0],
		text: cleanText(parts[1]),
	}
	return r, nil
}

func cleanText(s string) string {
	return strings.TrimLeft(s, " ")
}
