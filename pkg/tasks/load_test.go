package tasks

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	tests := map[string]struct {
		r        io.ReadCloser
		expected []Task
	}{
		"valid": {
			r: io.NopCloser(strings.NewReader("Saturday: cook\nMonday: clean")),
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
		},
		"valid with newlines and spaces": {
			r: io.NopCloser(strings.NewReader("     \nSaturday: cook\n  \nMonday: clean\n\n")),
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
		},
		"empty": {
			r:        io.NopCloser(strings.NewReader("")),
			expected: []Task{},
		},
		"empty with newlines and spaces": {
			r:        io.NopCloser(strings.NewReader("\n  \n\n  ")),
			expected: []Task{},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			// setup
			resChan := make(chan Task, 100)
			testDone := make(chan struct{})
			result := []Task{}
			go func() {
				for res := range resChan {
					result = append(result, res)
				}
				testDone <- struct{}{}
			}()

			err := scan(context.Background(), test.r, newTestTask, resChan)

			// shutdown
			close(resChan)
			<-testDone

			if err != nil {
				t.Fatalf("unexpected non-nil error: %v", err)
			}
			assertEqualTestTaskSlice(t, test.expected, result)
		})
	}
}
