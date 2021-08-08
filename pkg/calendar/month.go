package calendar

import "time"

// DaysInMonth calculates the number of days in the month of the time.Time object
func DaysInMonth(t time.Time) int {
	year, month, _ := t.Date()
	// start from the first of the current month, go forward one month and back one day
	monthEndDate := time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).AddDate(0, 1, -1)
	return monthEndDate.Day()
}
