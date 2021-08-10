package tasks

import (
	"sort"
	"testing"
	"time"
)

func TestRecurringDaysFrom(t *testing.T) {
	tests := map[string]struct {
		r        *recurring
		now      time.Time
		expected int
	}{
		"same day with single date": {
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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
			r: &recurring{
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

func Test_newRecurring(t *testing.T) {
	tests := map[string]struct {
		raw           *rawLine
		expectedDates []*date
		expectedText  string
		shouldErr     bool
	}{
		"empty": {
			raw:       &rawLine{},
			shouldErr: true,
		},
		"invalid month": {
			raw: &rawLine{
				date: "xxx",
			},
			shouldErr: true,
		},
		"invalid second month": {
			raw: &rawLine{
				date: "april/xxx",
			},
			shouldErr: true,
		},
		"empty second month": {
			raw: &rawLine{
				date: "april/",
			},
			shouldErr: true,
		},
		"month without day": {
			raw: &rawLine{
				date: "april/may",
			},
			shouldErr: true,
		},
		"month with invalid day": {
			raw: &rawLine{
				date: "april/may xxx",
			},
			shouldErr: true,
		},
		"single month": {
			raw: &rawLine{
				date: "april 1",
				text: "foo bar woo",
			},
			expectedDates: []*date{
				{
					month: time.April,
					day:   1,
				},
			},
			expectedText: "foo bar woo",
			shouldErr:    false,
		},
		"multiple months": {
			raw: &rawLine{
				date: "april/may 12",
				text: "foo bar woo",
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
			shouldErr:    false,
		},
		"multiple months unordered": {
			raw: &rawLine{
				date: "april/may/january 15",
				text: "foo bar woo",
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
			shouldErr:    false,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			res, err := newRecurring(test.raw)
			if test.shouldErr {
				if err == nil {
					t.Fatal("expected error but result err is nil")
				}
				return
			}
			if !test.shouldErr && err != nil {
				t.Fatalf("expected nil error but result err is %v", err)
			}
			result, ok := res.(*recurring)
			if !ok {
				t.Fatal("type assertion to *weekly failed on result")
			}
			if result.text != test.expectedText {
				t.Fatalf("result text '%s' not equal to expected text '%s'", result.text, test.expectedText)
			}

			if len(result.dates) != len(test.expectedDates) {
				t.Fatalf("number of result dates %d not equal to number of expected dates %d", len(result.dates), len(test.expectedDates))
			}

			resultDates := result.dates
			expectedDates := test.expectedDates
			sort.Slice(resultDates, func(i, j int) bool {
				return resultDates[i].month > resultDates[j].month
			})
			sort.Slice(expectedDates, func(i, j int) bool {
				return expectedDates[i].month > expectedDates[j].month
			})
			for i, d := range resultDates {
				ed := expectedDates[i]
				if d.month != ed.month || d.day != ed.day {
					t.Fatalf("unexpected result date %v", d)
				}
			}
		})
	}

}