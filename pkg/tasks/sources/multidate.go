package sources

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

const dayHours int = 24 * 60 * 60

// MultiDate represents a task with multiple dates
type MultiDate struct {
	month time.Month
	day   int
	text  string
}

// NewMultiDate constructs a MultiDate
func NewMultiDate(raw *RawTask) (*MultiDate, error) {
	dateParts := strings.SplitN(raw.Date, " ", 2)
	if len(dateParts) != 2 {
		return &MultiDate{}, fmt.Errorf("invalid multiple date [%s]", raw.Date)
	}

	month, err := calendar.ParseMonth(dateParts[0])
	if err != nil {
		return &MultiDate{}, fmt.Errorf("invalid multiple date [%s]", raw.Date)
	}
	day, err := strconv.ParseInt(dateParts[1], 10, 0)
	if err != nil {
		return &MultiDate{}, fmt.Errorf("could not parse date: %v", err)
	}

	m := &MultiDate{
		month: month,
		day:   int(day),
		text:  raw.Text,
	}
	return m, nil
}

// DaysFrom calculates the number of days until a task's date
func (m *MultiDate) DaysFrom(t time.Time) int {
	// set the task's year as the input year
	mTime := time.Date(t.Year(), m.month, m.day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	tUnix := t.Unix()

	days := int(mTime.Unix()-tUnix) / dayHours
	if days >= 0 {
		return days
	}
	// wrap task date to the next year and recalculate day difference
	mTime = mTime.AddDate(1, 0, 0)
	return int(mTime.Unix()-tUnix) / dayHours
}

func (m *MultiDate) String() string {
	return m.text
}
