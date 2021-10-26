package sources

import (
	"testing"
	"time"
)

func TestMonthlyDaysFrom(t *testing.T) {
	tests := map[string]struct {
		m        *Monthly
		now      time.Time
		expected int
	}{
		"same day": {
			m:        &Monthly{day: 6},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"next day": {
			m:        &Monthly{day: 7},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"five days from now": {
			m:        &Monthly{day: 11},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 5,
		},
		"previous day for month with 31 days": {
			m:        &Monthly{day: 5},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
		"five days before for month with 31 days": {
			m:        &Monthly{day: 1},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 26,
		},
		"previous day for month with 30 days": {
			m:        &Monthly{day: 5},
			now:      time.Date(2021, time.June, 6, 0, 0, 0, 0, time.UTC),
			expected: 29,
		},
		"previous day for February": {
			m:        &Monthly{day: 5},
			now:      time.Date(2021, time.February, 6, 0, 0, 0, 0, time.UTC),
			expected: 27,
		},
		"previous day for February leap year": {
			m:        &Monthly{day: 5},
			now:      time.Date(2024, time.February, 6, 0, 0, 0, 0, time.UTC),
			expected: 28,
		},
		"scheduled for 30th in February": {
			m:        &Monthly{day: 30},
			now:      time.Date(2021, time.February, 20, 0, 0, 0, 0, time.UTC),
			expected: 10,
		},
		"scheduled for 30th with month rolled over to March": {
			m:        &Monthly{day: 30},
			now:      time.Date(2021, time.March, 1, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"scheduled for 30th with month rolled over to March in leap year": {
			m:        &Monthly{day: 30},
			now:      time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"scheduled for 29th with month rolled over to March in leap year": {
			m:        &Monthly{day: 29},
			now:      time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC),
			expected: 28,
		},
		"scheduled for 31st with month rolled over to October": {
			m:        &Monthly{day: 31},
			now:      time.Date(2021, time.October, 1, 0, 0, 0, 0, time.UTC),
			expected: 0,
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
		raw          *RawTask
		expectedDay  int
		expectedText string
	}{
		"valid": {
			raw: &RawTask{
				Date: "12",
				Text: "foo bar woo",
			},
			expectedDay:  12,
			expectedText: "foo bar woo",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result, err := NewMonthly(test.raw)
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

func TestNewMonthlyError(t *testing.T) {
	tests := map[string]struct {
		raw *RawTask
	}{
		"empty": {
			raw: &RawTask{},
		},
		"invalid date": {
			raw: &RawTask{
				Date: "not a number",
			},
		},
		"day is zero": {
			raw: &RawTask{
				Date: "0",
			},
		},
		"day is negative": {
			raw: &RawTask{
				Date: "-1",
			},
		},
		"day is out of range": {
			raw: &RawTask{
				Date: "32",
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			_, err := NewMonthly(test.raw)
			if err == nil {
				t.Fatal("unexpected nil error")
			}
		})
	}
}
