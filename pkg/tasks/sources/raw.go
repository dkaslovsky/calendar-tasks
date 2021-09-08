package sources

import (
	"fmt"
	"strings"
)

const (
	dateTextSeparator  = ":"
	multiDateSeparator = "/"
)

// RawLine represents a line from an input source file - RENAME TO rawTask
type RawTask struct {
	Date string
	Text string
}

// ParseLine parses a line from an input source file into a slice of one or more RawLines
func ParseLine(line string) ([]*RawTask, error) {
	rts := []*RawTask{}

	parts := strings.SplitN(line, dateTextSeparator, 2)
	if len(parts) != 2 {
		return rts, fmt.Errorf("invalid line [%s]", line)
	}

	dateParts := strings.Split(parts[0], multiDateSeparator)
	text := cleanText(parts[1])

	for _, date := range dateParts {
		rt := &RawTask{
			Date: date,
			Text: text,
		}
		rts = append(rts, rt)
	}
	return rts, nil
}

func cleanText(s string) string {
	return strings.TrimLeft(s, " ")
}
