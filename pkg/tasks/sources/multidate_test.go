package sources

import (
	"sort"
	"testing"
	"time"
)

func assertEqualDateSlice(t *testing.T, expected, actual []*date) {
	if len(expected) != len(actual) {
		t.Fatalf("number of result dates %d not equal to number of expected dates %d", len(actual), len(expected))
	}

	actualDates := actual
	expectedDates := expected
	sort.Slice(actualDates, func(i, j int) bool {
		return actualDates[i].month > actualDates[j].month
	})
	sort.Slice(expectedDates, func(i, j int) bool {
		return expectedDates[i].month > expectedDates[j].month
	})

	for i := 0; i < len(expectedDates); i++ {
		a := actualDates[i]
		e := expectedDates[i]
		if a.month != e.month || a.day != e.day {
			t.Fatalf("unexpected result date %v", a)
		}
	}
}

func TestMultiDateDaysFrom(t *testing.T) {
	tests := map[string]struct {
		r        *MultiDate
		now      time.Time
		expected int
	}{
		"same day with single date": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.August,
						day:   6,
					},
				},
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"same day with multiple dates": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.August,
						day:   6,
					},
					{
						month: time.November,
						day:   6,
					},
				},
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"next day with single date": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.August,
						day:   7,
					},
				},
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"next day with multiple dates": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.August,
						day:   7,
					},
					{
						month: time.November,
						day:   7,
					},
				},
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"previous day": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.August,
						day:   5,
					},
				},
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 364,
		},
		"previous day including leap year starting past February": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.August,
						day:   5,
					},
				},
			},
			now:      time.Date(2023, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"previous day including leap year starting before February": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.January,
						day:   5,
					},
				},
			},
			now:      time.Date(2024, time.January, 6, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"previous day including leap year starting in February": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.February,
						day:   5,
					},
				},
			},
			now:      time.Date(2024, time.February, 6, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"previous day including leap year starting on February 28": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.February,
						day:   27,
					},
				},
			},
			now:      time.Date(2024, time.February, 28, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"previous day including leap year starting on February 29": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.February,
						day:   28,
					},
				},
			},
			now:      time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			expected: 365,
		},
		"next month single date": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.September,
						day:   6,
					},
				},
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"next month multiple dates": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.September,
						day:   6,
					},
					{
						month: time.November,
						day:   6,
					},
				},
			},
			now:      time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"within 24 hours": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.August,
						day:   15,
					},
				},
			},
			now:      time.Date(2021, time.August, 14, 18, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"one second away from today": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.August,
						day:   15,
					},
				},
			},
			now:      time.Date(2021, time.August, 14, 23, 59, 59, 0, time.UTC),
			expected: 1,
		},
		"exactly same time": {
			r: &MultiDate{
				dates: []*date{
					{
						month: time.August,
						day:   15,
					},
				},
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
		raw           *RawTask
		expectedDates []*date
		expectedText  string
	}{
		"single month": {
			raw: &RawTask{
				Date: "april 1",
				Text: "foo bar woo",
			},
			expectedDates: []*date{
				{
					month: time.April,
					day:   1,
				},
			},
			expectedText: "foo bar woo",
		},
		"multiple months": {
			raw: &RawTask{
				Date: "april/may 12",
				Text: "foo bar woo",
			},
			expectedDates: []*date{
				{
					month: time.April,
					day:   12,
				},
				{
					month: time.May,
					day:   12,
				},
			},
			expectedText: "foo bar woo",
		},
		"multiple months unordered": {
			raw: &RawTask{
				Date: "april/may/january 15",
				Text: "foo bar woo",
			},
			expectedDates: []*date{
				{
					month: time.April,
					day:   15,
				},
				{
					month: time.May,
					day:   15,
				},
				{
					month: time.January,
					day:   15,
				},
			},
			expectedText: "foo bar woo",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result, err := NewMultiDate(test.raw)
			if err != nil {
				t.Fatalf("unexpected non-nil error: %v", err)
			}
			if result.text != test.expectedText {
				t.Fatalf("result text '%s' not equal to expected text '%s'", result.text, test.expectedText)
			}
			assertEqualDateSlice(t, test.expectedDates, result.dates)
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
		"invalid second month": {
			raw: &RawTask{
				Date: "april/xxx",
			},
		},
		"empty second month": {
			raw: &RawTask{
				Date: "april/",
			},
		},
		"month without day": {
			raw: &RawTask{
				Date: "april/may",
			},
		},
		"month with invalid day": {
			raw: &RawTask{
				Date: "april/may xxx",
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
