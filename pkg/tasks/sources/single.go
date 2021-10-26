package sources

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

// Single represents a single-occurrence task
type Single struct {
	day   int
	month time.Month
	year  int
	text  string
}

// NewSingle constructs a Single
func NewSingle(raw *RawTask) (*Single, error) {
	dateParts := strings.SplitN(raw.Date, " ", 3)
	if len(dateParts) != 3 {
		return &Single{}, fmt.Errorf("invalid single date [%s]", raw.Date)
	}

	month, err := calendar.ParseMonth(dateParts[0])
	if err != nil {
		return &Single{}, fmt.Errorf("invalid single date [%s]", raw.Date)
	}
	day, err := strconv.ParseInt(dateParts[1], 10, 0)
	if err != nil {
		return &Single{}, fmt.Errorf("could not parse date: %v", err)
	}
	if day <= 0 || day > 31 {
		return &Single{}, fmt.Errorf("could not parse date: %v", err)
	}
	year, err := strconv.ParseInt(dateParts[2], 10, 0)
	if err != nil {
		return &Single{}, fmt.Errorf("could not parse date: %v", err)
	}
	if year < 0 {
		return &Single{}, fmt.Errorf("could not parse date: %v", err)
	}

	s := &Single{
		day:   int(day),
		month: month,
		year:  int(year),
		text:  raw.Text,
	}
	return s, nil
}

// DaysFrom calculates the number of days until a task's date
func (s *Single) DaysFrom(t time.Time) int {
	sTime := time.Date(s.year, s.month, s.day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	days := calendar.UnixToDaysFloored(sTime.Unix() - t.Unix())
	return int(days)
}

func (s *Single) String() string {
	return s.text
}
