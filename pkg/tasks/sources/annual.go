package sources

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

const dayHours int = 24 * 60 * 60

// Annual represents an annual task
type Annual struct {
	month time.Month
	day   int
	text  string
}

// NewAnnual constructs an Annual
func NewAnnual(raw *RawTask) (*Annual, error) {
	dateParts := strings.SplitN(raw.Date, " ", 2)
	if len(dateParts) != 2 {
		return &Annual{}, fmt.Errorf("invalid annual date [%s]", raw.Date)
	}

	month, err := calendar.ParseMonth(dateParts[0])
	if err != nil {
		return &Annual{}, fmt.Errorf("invalid annual date [%s]", raw.Date)
	}
	day, err := strconv.ParseInt(dateParts[1], 10, 0)
	if err != nil {
		return &Annual{}, fmt.Errorf("could not parse date: %v", err)
	}

	a := &Annual{
		month: month,
		day:   int(day),
		text:  raw.Text,
	}
	return a, nil
}

// DaysFrom calculates the number of days until a task's date
func (a *Annual) DaysFrom(t time.Time) int {
	// set the task's year as the input year
	aTime := time.Date(t.Year(), a.month, a.day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	tUnix := t.Unix()

	days := int(aTime.Unix()-tUnix) / dayHours
	if days >= 0 {
		return days
	}
	// wrap task date to the next year and recalculate day difference
	aTime = aTime.AddDate(1, 0, 0)
	return int(aTime.Unix()-tUnix) / dayHours
}

func (a *Annual) String() string {
	return a.text
}
