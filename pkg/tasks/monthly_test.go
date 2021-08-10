package tasks

import (
	"testing"
	"time"
)

func TestMonthlyDaysFrom(t *testing.T) {
	tests := map[string]struct {
		m        *monthly
		now      time.Time
		expected int
	}{
		"same day": {
			m:        &monthly{day: 6},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"next day": {
			m:        &monthly{day: 7},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"five days from now": {
			m:        &monthly{day: 11},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 5,
		},
		"previous day for month with 31 days": {
			m:        &monthly{day: 5},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
		"five days before for month with 31 days": {
			m:        &monthly{day: 1},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 26,
		},
		"previous day for month with 30 days": {
			m:        &monthly{day: 5},
			now:      time.Date(2021, time.June, 6, 0, 0, 0, 0, time.UTC),
			expected: 29,
		},
		"previous day for February": {
			m:        &monthly{day: 5},
			now:      time.Date(2021, time.February, 6, 0, 0, 0, 0, time.UTC),
			expected: 27,
		},
		"previous day for February leap year": {
			m:        &monthly{day: 5},
			now:      time.Date(2024, time.February, 6, 0, 0, 0, 0, time.UTC),
			expected: 28,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result := test.m.DaysFrom(test.now)
			if result != test.expected {
				t.Fatalf("result days %d not equal to expected days %d", result, test.expected)
			}
		})
	}
}

func TestNewMonthly(t *testing.T) {
	tests := map[string]struct {
		raw          *rawLine
		expectedDay  int
		expectedText string
		shouldErr    bool
	}{
		"empty": {
			raw:       &rawLine{},
			shouldErr: true,
		},
		"invalid date": {
			raw: &rawLine{
				date: "not a number",
			},
			shouldErr: true,
		},
		"non-empty": {
			raw: &rawLine{
				date: "12",
				text: "foo bar woo",
			},
			expectedDay:  12,
			expectedText: "foo bar woo",
			shouldErr:    false,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			res, err := newMonthly(test.raw)
			assertShouldError(t, test.shouldErr, err)
			if test.shouldErr {
				return
			}

			result, ok := res.(*monthly)
			if !ok {
				t.Fatal("type assertion failed on result")
			}
			if result.day != test.expectedDay {
				t.Fatalf("result days %d not equal to expected days %d", result.day, test.expectedDay)
			}
			if result.text != test.expectedText {
				t.Fatalf("result text '%s' not equal to expected text '%s'", result.text, test.expectedText)
			}
		})
	}

}
