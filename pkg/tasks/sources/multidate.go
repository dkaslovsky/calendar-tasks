package sources

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

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
	day, err := strconv.ParseInt(dateParts[1], 10, 0)
	if err != nil {
		return &MultiDate{}, fmt.Errorf("could not parse date: %v", err)
	}
	month, err := calendar.ParseMonth(dateParts[0])
	if err != nil {
		return &MultiDate{}, fmt.Errorf("invalid multiple date [%s]", raw.Date)
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
	nowYear := t.Year()
	nowMonth := t.Month()
	nowLoc := t.Location()

	type date struct {
		month time.Month
		day   int
	}

	dates := []*date{
		{
			month: m.month,
			day:   m.day,
		},
	}

	curDiff := 10e8 // any large value > 366 will work
	for _, date := range dates {
		curT := time.Date(nowYear, date.month, date.day, 0, 0, 0, 0, nowLoc)
		diff := curT.Sub(t).Hours() / 24
		if diff < 0 {
			if nowMonth > time.February && calendar.IsLeapYear(nowYear+1) {
				diff += 366
			} else if nowMonth <= time.February && calendar.IsLeapYear(nowYear) {
				diff += 366
			} else {
				diff += 365
			}
		}
		if diff < curDiff {
			curDiff = diff
		}
	}

	return int(math.Ceil(curDiff))
}

func (m *MultiDate) String() string {
	return m.text
}
