package sources

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

// Monthly represents a monthly task
type Monthly struct {
	day  int
	text string
}

// NewMonthly constructs a Monthly
func NewMonthly(raw *RawTask) (*Monthly, error) {
	day, err := strconv.ParseInt(raw.Date, 10, 0)
	if err != nil {
		return &Monthly{}, fmt.Errorf("could not parse date: %v", err)
	}
	if day <= 0 || day > 31 {
		return &Monthly{}, fmt.Errorf("could not parse date: %v", err)
	}

	m := &Monthly{
		day:  int(day),
		text: raw.Text,
	}
	return m, nil
}

// DaysFrom calculates the number of days until a task's date
func (m *Monthly) DaysFrom(t time.Time) int {
	// handle the case where the day is bigger than the number of days in the month
	if d := calendar.DaysInMonth(t.AddDate(0, -1, 0)); d < m.day {
		diff := m.day - (t.Day() + d)
		if diff >= 0 {
			return diff
		}
	}

	diff := m.day - t.Day()
	if diff >= 0 {
		return diff
	}
	return diff + calendar.DaysInMonth(t)
}

func (m *Monthly) String() string {
	return m.text
}
