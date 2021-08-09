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

const monthDelim = "/"

type date struct {
	Month time.Month
	Day   int
}

type RepeatedMonthly struct {
	Dates []*date
	Text  string
}

func NewRepeatedMonthly(line string) (*RepeatedMonthly, error) {
	raw, err := loadLine(line)
	if err != nil {
		return &RepeatedMonthly{}, nil
	}

	dateParts := strings.SplitN(raw.date, " ", 2)
	if len(dateParts) != 2 {
		return &RepeatedMonthly{}, fmt.Errorf("invalid repeated monthly date [%s]", raw.date)
	}
	day, err := strconv.ParseInt(dateParts[1], 10, 0)
	if err != nil {
		return &RepeatedMonthly{}, fmt.Errorf("could not parse date: %v", err)
	}
	months := strings.Split(dateParts[0], monthDelim)
	if len(months) == 0 || len(months) > 12 {
		return &RepeatedMonthly{}, fmt.Errorf("invalid repeated monthly date [%s]", raw.date)
	}

	dates := []*date{}
	for _, m := range months {
		month, err := calendar.ParseMonth(m)
		if err != nil {
			return &RepeatedMonthly{}, fmt.Errorf("invalid repeated monthly date [%s]", raw.date)
		}
		dates = append(dates, &date{
			Month: month,
			Day:   int(day),
		})
	}

	m := &RepeatedMonthly{
		Dates: dates,
		Text:  raw.text,
	}
	return m, nil
}

func (m *RepeatedMonthly) DaysFrom(t time.Time) int {
	return 0
}

func (m *RepeatedMonthly) String() string {
	s, _ := json.MarshalIndent(m, "", "\t")
	return string(s)
}

func LoadRepeatedMonthly(fileName string) ([]*RepeatedMonthly, error) {
	ms := []*RepeatedMonthly{}

	b, err := os.ReadFile(fileName)
	if err != nil {
		return ms, fmt.Errorf("failed to read file: %v", err)
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		m, err := NewRepeatedMonthly(line)
		if err != nil {
			return ms, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}
