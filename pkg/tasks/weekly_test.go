package tasks

import (
	"testing"
	"time"
)

func TestWeeklyDaysFrom(t *testing.T) {
	tests := map[string]struct {
		wk       *weekly
		tm       time.Time
		expected int
	}{
		"same day": {
			wk:       &weekly{day: time.Friday},
			tm:       time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"next day": {
			wk:       &weekly{day: time.Saturday},
			tm:       time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"previous day": {
			wk:       &weekly{day: time.Thursday},
			tm:       time.Date(2021, time.August, 6, 0, 0, 0, 0, time.UTC),
			expected: 6,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result := test.wk.DaysFrom(test.tm)
			if result != test.expected {
				t.Fatalf("result days %d not equal to expected days %d", result, test.expected)
			}
		})
	}
}

func Test_newWeekly(t *testing.T) {
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
			if test.shouldErr {
				if err == nil {
					t.Fatal("expected error but result err is nil")
				}
				return
			}
			if !test.shouldErr && err != nil {
				t.Fatalf("expected nil error but result err is %v", err)
			}
			result, ok := res.(*weekly)
			if !ok {
				t.Fatal("type assertion to *weekly failed on result")
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
