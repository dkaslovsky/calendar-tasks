package tasks

import (
	"testing"
	"time"
)

func TestWeeklyDaysFrom(t *testing.T) {
	tests := map[string]struct {
		w        *weekly
		now      time.Time
		expected int
	}{
		"same day": {
			w:        &weekly{day: time.Friday},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"next day": {
			w:        &weekly{day: time.Saturday},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"five days from now": {
			w:        &weekly{day: time.Wednesday},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 5,
		},
		"previous day": {
			w:        &weekly{day: time.Thursday},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 6,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result := test.w.DaysFrom(test.now)
			if result != test.expected {
				t.Fatalf("result days %d not equal to expected days %d", result, test.expected)
			}
		})
	}
}

func TestNewWeekly(t *testing.T) {
	tests := map[string]struct {
		raw          *rawLine
		expectedDay  time.Weekday
		expectedText string
		shouldErr    bool
	}{
		"empty": {
			raw:       &rawLine{},
			shouldErr: true,
		},
		"invalid date": {
			raw: &rawLine{
				date: "funday",
			},
			shouldErr: true,
		},
		"non-empty": {
			raw: &rawLine{
				date: "Monday",
				text: "foo bar woo",
			},
			expectedDay:  time.Monday,
			expectedText: "foo bar woo",
			shouldErr:    false,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			res, err := newWeekly(test.raw)
			assertShouldError(t, test.shouldErr, err)
			if test.shouldErr {
				return
			}
			result, ok := res.(*weekly)
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
