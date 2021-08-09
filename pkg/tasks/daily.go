package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
)

type Weekly struct {
	Day  time.Weekday
	Text string
}

func NewWeekly(line string) (*Weekly, error) {
	raw, err := loadLine(line)
	if err != nil {
		return &Weekly{}, nil
	}

	day, err := calendar.ParseWeekday(raw.date)
	if err != nil {
		return &Weekly{}, fmt.Errorf("could not parse date: %v", err)
	}

	d := &Weekly{
		Day:  day,
		Text: raw.text,
	}
	return d, nil
}

func (d *Weekly) DaysFrom(t time.Time) int {
	return calendar.DaysBetweenWeekdays(t.Weekday(), d.Day)
}

func (d *Weekly) String() string {
	s, _ := json.MarshalIndent(map[string]string{
		"Day":  d.Day.String(),
		"Text": d.Text,
	}, "", "\t")
	return string(s)
}

func LoadWeekly(fileName string) ([]*Weekly, error) {
	ds := []*Weekly{}

	b, err := os.ReadFile(fileName)
	if err != nil {
		return ds, fmt.Errorf("failed to read file: %v", err)
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		d, err := NewWeekly(line)
		if err != nil {
			return ds, err
		}
		ds = append(ds, d)
	}
	return ds, nil
}
