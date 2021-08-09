package tasks

import (
	"encoding/json"
	"fmt"
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

type Recurring struct {
	Dates []*date
	Text  string
}

func NewRecurring(raw *rawLine) (Task, error) {
	dateParts := strings.SplitN(raw.date, " ", 2)
	if len(dateParts) != 2 {
		return &Recurring{}, fmt.Errorf("invalid recurring date [%s]", raw.date)
	}
	day, err := strconv.ParseInt(dateParts[1], 10, 0)
	if err != nil {
		return &Recurring{}, fmt.Errorf("could not parse date: %v", err)
	}
	months := strings.Split(dateParts[0], monthDelim)
	if len(months) == 0 || len(months) > 12 {
		return &Recurring{}, fmt.Errorf("invalid recurring date [%s]", raw.date)
	}

	dates := []*date{}
	for _, m := range months {
		month, err := calendar.ParseMonth(m)
		if err != nil {
			return &Recurring{}, fmt.Errorf("invalid recurring date [%s]", raw.date)
		}
		dates = append(dates, &date{
			Month: month,
			Day:   int(day),
		})
	}

	m := &Recurring{
		Dates: dates,
		Text:  raw.text,
	}
	return m, nil
}

func (m *Recurring) DaysFrom(t time.Time) int {
	nowYear := t.Year()
	nowLoc := t.Location()

	curDiff := 10e8
	for _, date := range m.Dates {
		curT := time.Date(nowYear, date.Month, date.Day, 0, 0, 0, 0, nowLoc)
		diff := curT.Sub(t).Hours() / 24
		if diff < 0 {
			diff += 365
		}
		if diff < curDiff {
			curDiff = diff
		}
	}

	// cast to integer does not lose precision because day is finest time granularity available
	return int(curDiff)
}

func (m *Recurring) String() string {
	s, _ := json.MarshalIndent(m, "", "\t")
	return string(s)
}

func LoadRecurring(fileName string) ([]Task, error) {
	return Load(fileName, NewRecurring)
}
