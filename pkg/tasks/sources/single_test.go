package sources

import (
	"testing"
	"time"
)

// func TestSingleDaysFrom(t *testing.T) {
// 	t.Fatal("write tests")
// }

func TestNewSingle(t *testing.T) {
	tests := map[string]struct {
		raw      *RawTask
		expected *Single
	}{
		"single date": {
			raw: &RawTask{
				Date: "april 17 2021",
				Text: "foo bar woo",
			},
			expected: &Single{
				month: time.April,
				day:   17,
				year:  2021,
				text:  "foo bar woo",
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result, err := NewSingle(test.raw)
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
			if result.year != test.expected.year {
				t.Fatalf("result year '%d' not equal to expected year '%d'", result.year, test.expected.year)
			}
		})
	}
}

func TestNewSingleError(t *testing.T) {
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
		"month only": {
			raw: &RawTask{
				Date: "april",
			},
		},
		"month and day only": {
			raw: &RawTask{
				Date: "april 17",
			},
		},
		"month with invalid day": {
			raw: &RawTask{
				Date: "april 1xxx 2021",
			},
		},
		"invalid month with valid day and year": {
			raw: &RawTask{
				Date: "xxx 17 2021",
			},
		},
		"invalid year with valid day and month": {
			raw: &RawTask{
				Date: "April 17 -2021",
			},
		},
		"day out of range": {
			raw: &RawTask{
				Date: "April 32 2021",
			},
		},
		"negative day": {
			raw: &RawTask{
				Date: "April -3 2021",
			},
		},
		"day is zero": {
			raw: &RawTask{
				Date: "April 0 2021",
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			_, err := NewSingle(test.raw)
			if err == nil {
				t.Fatal("unexpected nil error")
			}
		})
	}
}
