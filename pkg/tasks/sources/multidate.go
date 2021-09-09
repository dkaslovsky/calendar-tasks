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
	nowYear := t.Year()
	nowMonth := t.Month()
	nowLoc := t.Location()

	curT := time.Date(nowYear, m.month, m.day, 0, 0, 0, 0, nowLoc)
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
	return int(math.Ceil(diff))
}

func (m *MultiDate) String() string {
	return m.text
}
