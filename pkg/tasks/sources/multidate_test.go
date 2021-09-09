package sources

import (
	"testing"
	"time"
)

func TestMultiDateDaysFrom(t *testing.T) {
	tests := map[string]struct {
		r        *MultiDate
		now      time.Time
		expected int
	}{
		"same day": {
			r: &MultiDate{
				month: time.August,
				day:   6,
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"next day": {
			r: &MultiDate{
				month: time.August,
				day:   7,
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"previous day": {
			r: &MultiDate{
				month: time.August,
				day:   5,
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 364,
		},
		"previous day including leap year starting past February": {
			r: &MultiDate{
				month: time.August,
				day:   5,
			},
			now:      time.Date(2023, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"previous day including leap year starting before February": {
			r: &MultiDate{
				month: time.January,
				day:   5,
			},
			now:      time.Date(2024, time.January, 6, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"previous day including leap year starting in February": {
			r: &MultiDate{
				month: time.February,
				day:   5,
			},
			now:      time.Date(2024, time.February, 6, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"previous day including leap year starting on February 28": {
			r: &MultiDate{
				month: time.February,
				day:   27,
			},
			now:      time.Date(2024, time.February, 28, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"previous day including leap year starting on February 29": {
			r: &MultiDate{
				month: time.February,
				day:   28,
			},
			now:      time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"next month": {
			r: &MultiDate{
				month: time.September,
				day:   6,
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"24 hours away": {
			r: &MultiDate{
				month: time.August,
				day:   15,
			},
			now:      time.Date(2021, time.August, 14, 18, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"within 24 hours": {
			r: &MultiDate{
				month: time.August,
				day:   15,
			},
			now:      time.Date(2021, time.August, 14, 18, 12, 0, 0, time.UTC),
			expected: 1,
		},
		"one second away from today": {
			r: &MultiDate{
				month: time.August,
				day:   15,
			},
			now:      time.Date(2021, time.August, 14, 23, 59, 59, 0, time.UTC),
			expected: 1,
		},
		"exactly same time": {
			r: &MultiDate{
				month: time.August,
				day:   15,
			},
			now:      time.Date(2021, time.August, 15, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result := test.r.DaysFrom(test.now)
			if result != test.expected {
				t.Fatalf("result days %d not equal to expected days %d", result, test.expected)
			}
		})
	}
}

func TestNewMultiDate(t *testing.T) {
	tests := map[string]struct {
		raw      *RawTask
		expected *MultiDate
	}{
		"single month": {
			raw: &RawTask{
				Date: "april 17",
				Text: "foo bar woo",
			},
			expected: &MultiDate{
				month: time.April,
				day:   17,
				text:  "foo bar woo",
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result, err := NewMultiDate(test.raw)
			if err != nil {
				t.Fatalf("unexpected non-nil error: %v", err)
			}
			if result.text != test.expected.text {
				t.Fatalf("result text '%s' not equal to expected text '%s'", result.text, test.expected.text)
			}
			if result.month != test.expected.month {
				t.Fatalf("result month '%s' not equal to expected month '%s'", result.month, test.expected.month)
			}
			if result.day != test.expected.day {
				t.Fatalf("result day '%d' not equal to expected day '%d'", result.day, test.expected.day)
			}
		})
	}

}

func TestNewMultiDateError(t *testing.T) {
	tests := map[string]struct {
		raw *RawTask
	}{
		"empty": {
			raw: &RawTask{},
		},
		"invalid month": {
			raw: &RawTask{
				Date: "xxx",
			},
		},
		"month without day": {
			raw: &RawTask{
				Date: "april",
			},
		},
		"month with invalid day": {
			raw: &RawTask{
				Date: "april 1xxx",
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			_, err := NewMultiDate(test.raw)
			if err == nil {
				t.Fatal("unexpected nil error")
			}
		})
	}

}
