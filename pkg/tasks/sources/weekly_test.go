package sources

import (
	"testing"
	"time"
)

func TestWeeklyDaysFrom(t *testing.T) {
	tests := map[string]struct {
		w        *Weekly
		now      time.Time
		expected int
	}{
		"same day": {
			w:        &Weekly{day: time.Friday},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"next day": {
			w:        &Weekly{day: time.Saturday},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"five days from now": {
			w:        &Weekly{day: time.Wednesday},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 5,
		},
		"previous day": {
			w:        &Weekly{day: time.Thursday},
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
		raw          *RawLine
		expectedDay  time.Weekday
		expectedText string
	}{
		"valid": {
			raw: &RawLine{
				Date: "Monday",
				Text: "foo bar woo",
			},
			expectedDay:  time.Monday,
			expectedText: "foo bar woo",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result, err := NewWeekly(test.raw)
			if err != nil {
				t.Fatalf("unexpected non-nil error: %v", err)
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

func TestNewWeeklyError(t *testing.T) {
	tests := map[string]struct {
		raw *RawLine
	}{
		"empty": {
			raw: &RawLine{},
		},
		"invalid date": {
			raw: &RawLine{
				Date: "funday",
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			_, err := NewWeekly(test.raw)
			if err == nil {
				t.Fatal("unexpected nil error")
			}
		})
	}
}
