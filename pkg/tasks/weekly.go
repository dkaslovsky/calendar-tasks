package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

// LoadWeekly loads weekly tasks from a file and sends them on a specified channel
func LoadWeekly(fileName string, taskCh chan Task, done chan struct{}) error {
	return load(fileName, newWeekly, taskCh, done)
}

type weekly struct {
	day  time.Weekday
	text string
}

func newWeekly(raw *rawLine) (Task, error) {
	day, err := calendar.ParseWeekday(raw.date)
	if err != nil {
		return nil, fmt.Errorf("could not parse date: %v", err)
	}

	w := &weekly{
		day:  day,
		text: raw.text,
	}
	return w, nil
}

func (w *weekly) DaysFrom(t time.Time) int {
	return calendar.DaysBetweenWeekdays(t.Weekday(), w.day)
}

func (w *weekly) String() string {
	s, _ := json.MarshalIndent(map[string]string{
		"Day":  w.day.String(),
		"Text": w.text,
	}, "", "\t")
	return string(s)
}
