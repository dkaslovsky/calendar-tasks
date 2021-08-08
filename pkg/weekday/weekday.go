package weekday

import (
	"fmt"
	"strings"
	"time"
)

var weekdays = map[string]time.Weekday{}

func init() {
	for d := time.Sunday; d <= time.Saturday; d++ {
		name := strings.ToLower(d.String())
		weekdays[name] = d
		weekdays[name[:3]] = d
	}
}

// Parse converts a string to its corresponding time.Weekday
func Parse(s string) (time.Weekday, error) {
	var day time.Weekday
	day, ok := weekdays[strings.ToLower(s)]
	if !ok {
		return day, fmt.Errorf("invalid weekday [%s]", day)
	}
	return day, nil
}

// DaysBetween returns the number of days between two time.Weekdays
func DaysBetween(now, target time.Weekday) int {
	diff := int(target - now)
	if diff >= 0 {
		return diff
	}
	return diff + 7
}
