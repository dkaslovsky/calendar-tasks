package reminders

import (
	"fmt"
	"strings"
)

const delim = ":"

type rawReminder struct {
	date string
	text string
}

func loadLine(line string) (*rawReminder, error) {
	parts := strings.Split(line, delim)
	if len(parts) != 2 {
		return &rawReminder{}, fmt.Errorf("invalid line [%s]", line)
	}

	r := &rawReminder{
		date: parts[0],
		text: cleanText(parts[1]),
	}
	return r, nil
}

func cleanText(s string) string {
	return strings.TrimLeft(s, " ")
}
