package reminders

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Monthly struct {
	Day  int
	Text string
}

func NewMonthly(line string) (*Monthly, error) {
	raw, err := loadLine(line)
	if err != nil {
		return &Monthly{}, nil
	}

	day, err := strconv.ParseInt(raw.date, 10, 0)
	if err != nil {
		return &Monthly{}, fmt.Errorf("could not parse date: %v", err)
	}

	m := &Monthly{
		Day:  int(day),
		Text: raw.text,
	}
	return m, nil
}

func (m *Monthly) DaysFrom(t time.Time) int {
	diff := int(m.Day - t.Day())
	if diff >= 0 {
		return diff
	}
	return diff + 31 // get days in current month instead of hardcoded 31
}

func (m *Monthly) String() string {
	s, _ := json.MarshalIndent(m, "", "\t")
	return string(s)
}

func (m *Monthly) GetText() string {
	return m.Text
}
