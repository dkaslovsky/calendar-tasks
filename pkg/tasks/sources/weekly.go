package sources

import (
	"fmt"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

// Weekly represents a weekly task
type Weekly struct {
	day  time.Weekday
	text string
}

// NewWeekly constructs a Weekly
func NewWeekly(raw *RawTask) (*Weekly, error) {
	day, err := calendar.ParseWeekday(raw.Date)
	if err != nil {
		return nil, fmt.Errorf("could not parse date: %v", err)
	}

	w := &Weekly{
		day:  day,
		text: raw.Text,
	}
	return w, nil
}

// DaysFrom calculates the number of days until a task's date
func (w *Weekly) DaysFrom(t time.Time) int {
	return calendar.DaysBetweenWeekdays(t.Weekday(), w.day)
}

func (w *Weekly) String() string {
	return w.text
}
