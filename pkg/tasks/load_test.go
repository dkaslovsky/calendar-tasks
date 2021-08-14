package tasks

import (
	"io"
	"strings"
	"testing"
)

func TestLoadLine(t *testing.T) {
	tests := map[string]struct {
		line      string
		expected  *rawLine
		shouldErr bool
	}{
		"empty": {
			line:      "",
			shouldErr: true,
		},
		"empty with spaces": {
			line:      "    ",
			shouldErr: true,
		},
		"no delimiter": {
			line:      "foobar",
			shouldErr: true,
		},
		"multiple delimiters": {
			line: "foo:bar:baz",
			expected: &rawLine{
				date: "foo",
				text: "bar:baz",
			},
			shouldErr: false,
		},
		"valid": {
			line: "foo:bar",
			expected: &rawLine{
				date: "foo",
				text: "bar",
			},
			shouldErr: false,
		},
		"valid with spaces": {
			line: "foo:  bar",
			expected: &rawLine{
				date: "foo",
				text: "bar",
			},
			shouldErr: false,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			result, err := loadLine(test.line)
			assertShouldError(t, test.shouldErr, err)
			if test.shouldErr {
				return
			}
			if result.date != test.expected.date || result.text != test.expected.text {
				t.Fatalf("result %v does not equal expected %v", result, test.expected)
			}
		})
	}
}

func TestScan(t *testing.T) {
	tests := map[string]struct {
		r         io.ReadCloser
		newTask   func(*rawLine) (Task, error)
		expected  []Task
		shouldErr bool
	}{
		"empty": {
			r:         io.NopCloser(strings.NewReader("")),
			shouldErr: true,
		},
		"empty with newlines and spaces": {
			r:         io.NopCloser(strings.NewReader("\n  \n\n  ")),
			shouldErr: true,
		},
		"valid": {
			r: io.NopCloser(strings.NewReader("Saturday: cook\nMonday: clean")),
			newTask: func(rl *rawLine) (Task, error) {
				return &testTask{
					id:       rl.text,
					daysFrom: 1,
				}, nil
			},
			expected: []Task{
				&testTask{
					id:       "cook",
					daysFrom: 1,
				},
				&testTask{
					id:       "clean",
					daysFrom: 1,
				},
			},
			shouldErr: false,
		},
		"valid with newlines and spaces": {
			r: io.NopCloser(strings.NewReader("     \nSaturday: cook\n  \nMonday: clean\n\n")),
			newTask: func(rl *rawLine) (Task, error) {
				return &testTask{
					id:       rl.text,
					daysFrom: 1,
				}, nil
			},
			expected: []Task{
				&testTask{
					id:       "cook",
					daysFrom: 1,
				},
				&testTask{
					id:       "clean",
					daysFrom: 1,
				},
			},
			shouldErr: false,
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {

			resChan := make(chan Task, 100)
			done := make(chan struct{})
			result := []Task{}
			go func() {
				for res := range resChan {
					result = append(result, res)
				}
				done <- struct{}{}
			}()

			err := scan(test.r, test.newTask, resChan)
			assertShouldError(t, test.shouldErr, err)
			if test.shouldErr {
				return
			}

			<-done
			assertEqualTestTaskSlice(t, test.expected, result)
		})
	}
}
