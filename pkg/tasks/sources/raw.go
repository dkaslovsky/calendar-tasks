package sources

import (
	"fmt"
	"strings"
)

const (
	dateTextSeparator  = ":"
	multiDateSeparator = "/"
)

// RawTask represents an unprocessed task parsed from a line of an input source file
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
	text := cleanString(parts[1])

	for _, date := range dateParts {
		rt := &RawTask{
			Date: cleanString(date),
			Text: text,
		}
		rts = append(rts, rt)
	}
	return rts, nil
}

func cleanString(s string) string {
	return strings.TrimSpace(s)
}
