package tasks

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

type Daily struct {
	Day  time.Weekday
	Text string
}

func NewDaily(line string) (*Daily, error) {
	raw, err := loadLine(line)
	if err != nil {
		return &Daily{}, nil
	}

	day, err := calendar.ParseWeekday(raw.date)
	if err != nil {
		return &Daily{}, fmt.Errorf("could not parse date: %v", err)
	}

	d := &Daily{
		Day:  day,
		Text: raw.text,
	}
	return d, nil
}

func (d *Daily) DaysFrom(t time.Time) int {
	return calendar.DaysBetweenWeekdays(t.Weekday(), d.Day)
}

func (d *Daily) String() string {
	s, _ := json.MarshalIndent(map[string]string{
		"Day":  d.Day.String(),
		"Text": d.Text,
	}, "", "\t")
	return string(s)
}
