package calendar

import (
	"testing"
	"time"
)

func TestDaysBetweenWeekdays(t *testing.T) {
	tests := map[string]struct {
		current  time.Weekday
		target   time.Weekday
		expected int
	}{
		"same day": {
			current:  time.Tuesday,
			target:   time.Tuesday,
			expected: 0,
		},
		"next day": {
			current:  time.Tuesday,
			target:   time.Wednesday,
			expected: 1,
		},
		"day before": {
			current:  time.Tuesday,
			target:   time.Monday,
			expected: 6,
		},
		"wrap around": {
			current:  time.Saturday,
			target:   time.Sunday,
			expected: 1,
		},
		"full range": {
			current:  time.Sunday,
			target:   time.Saturday,
			expected: 6,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result := DaysBetweenWeekdays(test.current, test.target)
			if result != test.expected {
				t.Fatalf("result days %d not equal to expected days %d", result, test.expected)
			}
		})
	}
}
