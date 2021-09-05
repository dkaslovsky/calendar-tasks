package tasks

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks/sources"
)

func TestScan(t *testing.T) {
	tests := map[string]struct {
		r         io.ReadCloser
		newTask   func(*sources.RawLine) (Task, error)
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
			newTask: func(rl *sources.RawLine) (Task, error) {
				return &testTask{
					id:       rl.Text,
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
			newTask: func(rl *sources.RawLine) (Task, error) {
				return &testTask{
					id:       rl.Text,
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
			testDone := make(chan struct{})
			result := []Task{}
			go func() {
				for res := range resChan {
					result = append(result, res)
				}
				testDone <- struct{}{}
			}()

			err := scan(context.Background(), test.r, test.newTask, resChan)
			assertShouldError(t, test.shouldErr, err)
			if test.shouldErr {
				return
			}

			close(resChan)
			<-testDone
			assertEqualTestTaskSlice(t, test.expected, result)
		})
	}
}
