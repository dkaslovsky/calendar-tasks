package tasks

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

type monthly struct {
	day  int
	text string
}

func newMonthly(raw *rawLine) (Task, error) {
	day, err := strconv.ParseInt(raw.date, 10, 0)
	if err != nil {
		return &monthly{}, fmt.Errorf("could not parse date: %v", err)
	}

	m := &monthly{
		day:  int(day),
		text: raw.text,
	}
	return m, nil
}

func (m *monthly) DaysFrom(t time.Time) int {
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

func (m *monthly) String() string {
	s, _ := json.MarshalIndent(map[string]string{
		"Day":  fmt.Sprint(m.day),
		"Text": m.text,
	}, "", "\t")
	return string(s)
}
