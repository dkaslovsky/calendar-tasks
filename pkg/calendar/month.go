package calendar

import (
	"fmt"
	"strings"
	"time"
)

var months = make(map[string]time.Month)

// similar to github.com/icza/gox/timex
func init() {
	for d := time.January; d <= time.December; d++ {
		name := strings.ToLower(d.String())
		months[name] = d
		months[name[:3]] = d
		if len(name) > 3 {
			months[name[:4]] = d
		}
	}
}

// ParseMonth converts a string to its corresponding time.Month
func ParseMonth(s string) (time.Month, error) {
	var month time.Month
	month, ok := months[strings.ToLower(s)]
	if !ok {
		return month, fmt.Errorf("invalid month [%s]", s)
	}
	return month, nil
}

// DaysInMonth calculates the number of days in the month of the time.Time object
func DaysInMonth(t time.Time) int {
	year, month, _ := t.Date()
	// start from the first of the current month, go forward one month and back one day
	monthEndDate := time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, -1)
	return monthEndDate.Day()
}
