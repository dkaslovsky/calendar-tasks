package calendar

import (
	"testing"
	"time"
)

func TestDaysInMonth(t *testing.T) {
	tests := map[string]struct {
		tm       time.Time
		expected int
	}{
		"January": {
			tm:       time.Date(2021, time.January, 2, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"February non leap year": {
			tm:       time.Date(2021, time.February, 2, 0, 0, 0, 0, time.UTC),
			expected: 28,
		},
		"February leap year": {
			tm:       time.Date(2024, time.February, 2, 0, 0, 0, 0, time.UTC),
			expected: 29,
		},
		"March": {
			tm:       time.Date(2021, time.March, 2, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"April": {
			tm:       time.Date(2021, time.April, 2, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
		"May": {
			tm:       time.Date(2021, time.May, 2, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"June": {
			tm:       time.Date(2021, time.June, 2, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
		"July": {
			tm:       time.Date(2021, time.July, 2, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"August": {
			tm:       time.Date(2021, time.August, 2, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"September": {
			tm:       time.Date(2021, time.September, 2, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
		"October": {
			tm:       time.Date(2021, time.October, 2, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"November": {
			tm:       time.Date(2021, time.November, 2, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
		"December": {
			tm:       time.Date(2021, time.December, 2, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result := DaysInMonth(test.tm)
			if result != test.expected {
				t.Fatalf("result days %d not equal to expected days %d", result, test.expected)
			}
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	t.Run("leap year 2024", func(t *testing.T) {
		if !IsLeapYear(2024) {
			t.Fatal("expected 2024 to be a leap year")
		}
	})
	t.Run("not leap year 2023", func(t *testing.T) {
		if IsLeapYear(2023) {
			t.Fatal("expected 2023 to not be a leap year")
		}
	})
}
