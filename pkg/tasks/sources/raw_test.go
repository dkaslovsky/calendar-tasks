package sources

import (
	"fmt"
	"sort"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := map[string]struct {
		line     string
		expected []*RawTask
	}{
		"valid": {
			line: "foo:bar",
			expected: []*RawTask{
				{
					Date: "foo",
					Text: "bar",
				},
			},
		},
		"valid with spaces": {
			line: "foo:  bar",
			expected: []*RawTask{
				{
					Date: "foo",
					Text: "bar",
				},
			},
		},
		"multiple delimiters": {
			line: "foo:bar:baz",
			expected: []*RawTask{
				{
					Date: "foo",
					Text: "bar:baz",
				},
			},
		},
		"valid with multiple annual dates": {
			line: "mar 1/april 2/jan 3: foo",
			expected: []*RawTask{
				{
					Date: "jan 3",
					Text: "foo",
				},
				{
					Date: "mar 1",
					Text: "foo",
				},
				{
					Date: "april 2",
					Text: "foo",
				},
			},
		},
		"valid with multiple weekly dates": {
			line: "sat/sun: foo/bar",
			expected: []*RawTask{
				{
					Date: "sat",
					Text: "foo/bar",
				},
				{
					Date: "sun",
					Text: "foo/bar",
				},
			},
		},
		"valid with multiple monthly dates": {
			line: "4/25/31:woo",
			expected: []*RawTask{
				{
					Date: "4",
					Text: "woo",
				},
				{
					Date: "25",
					Text: "woo",
				},
				{
					Date: "31",
					Text: "woo",
				},
			},
		},
		"valid with multiple dates and spaces": {
			line: "date1 / date2 / date3:woo",
			expected: []*RawTask{
				{
					Date: "date1",
					Text: "woo",
				},
				{
					Date: "date2",
					Text: "woo",
				},
				{
					Date: "date3",
					Text: "woo",
				},
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result, err := ParseLine(test.line)
			if err != nil {
				t.Fatalf("unexpected non-nil error: %v", err)
			}
			if len(result) != len(test.expected) {
				t.Fatalf("number of results %d not equal to expected number of results %d", len(result), len(test.expected))
			}
			sort.Slice(result, func(i, j int) bool {
				return testRawTaskSortKey(result[i]) > testRawTaskSortKey(result[j])
			})
			sort.Slice(test.expected, func(i, j int) bool {
				return testRawTaskSortKey(test.expected[i]) > testRawTaskSortKey(test.expected[j])
			})
			for i := 0; i < len(result); i++ {
				r := result[i]
				e := test.expected[i]
				if r.Date != e.Date || r.Text != e.Text {
					t.Fatalf("result %v does not equal expected %v", r, e)
				}
			}
		})
	}
}

func TestLoadLineError(t *testing.T) {
	tests := map[string]struct {
		line string
	}{
		"empty": {
			line: "",
		},
		"empty with spaces": {
			line: "    ",
		},
		"no date-text separator one string": {
			line: "foobar",
		},
		"no date-text separator multiple strings": {
			line: "foobar xyz aaabbb",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			_, err := ParseLine(test.line)
			if err == nil {
				t.Fatal("unexpected nil error")
			}
		})
	}
}

// testRawTaskSortKey creates a string for sorting RawTasks
func testRawTaskSortKey(r *RawTask) string {
	return fmt.Sprintf("Date-%s-Text-%s", r.Date, r.Text)
}
