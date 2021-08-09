package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

type Weekly struct {
	Day  time.Weekday
	Text string
}

func NewWeekly(raw *rawLine) (Task, error) {
	day, err := calendar.ParseWeekday(raw.date)
	if err != nil {
		return nil, fmt.Errorf("could not parse date: %v", err)
	}

	w := &Weekly{
		Day:  day,
		Text: raw.text,
	}
	return w, nil
}

func (w *Weekly) DaysFrom(t time.Time) int {
	return calendar.DaysBetweenWeekdays(t.Weekday(), w.Day)
}

func (w *Weekly) String() string {
	s, _ := json.MarshalIndent(map[string]string{
		"Day":  w.Day.String(),
		"Text": w.Text,
	}, "", "\t")
	return string(s)
}

func LoadWeekly(fileName string) ([]Task, error) {
	return Load(fileName, NewWeekly)
}
