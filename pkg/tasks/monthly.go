package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/calendar"
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
	return diff + calendar.DaysInMonth(t)
}

func (m *Monthly) String() string {
	s, _ := json.MarshalIndent(m, "", "\t")
	return string(s)
}

func LoadMonthly(fileName string) ([]*Monthly, error) {
	ms := []*Monthly{}

	b, err := os.ReadFile(fileName)
	if err != nil {
		return ms, fmt.Errorf("failed to read file: %v", err)
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		m, err := NewMonthly(line)
		if err != nil {
			return ms, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}
