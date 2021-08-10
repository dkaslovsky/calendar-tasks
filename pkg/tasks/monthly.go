package tasks

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

// LoadMonthly loads monthly tasks from a file and returns a slice of objects satisfying the Task interface
func LoadMonthly(fileName string) ([]Task, error) {
	return load(fileName, newMonthly)
}

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
	diff := int(m.day - t.Day())
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
