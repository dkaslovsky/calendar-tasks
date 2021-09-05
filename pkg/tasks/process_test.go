package tasks

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	tests := map[string]struct {
		maxDays       int
		tasks         []*testTask
		expectedTasks map[int][]Task
	}{
		"empty": {
			maxDays:       100,
			tasks:         []*testTask{},
			expectedTasks: make(map[int][]Task),
		},
		"single task": {
			maxDays: 4,
			tasks: []*testTask{
				{
					id:       "a",
					daysFrom: 4,
				},
			},
			expectedTasks: map[int][]Task{
				4: {
					&testTask{
						id:       "a",
						daysFrom: 4,
					},
				},
			},
		},
		"single task beyond maxDays": {
			maxDays: 3,
			tasks: []*testTask{
				{
					id:       "a",
					daysFrom: 4,
				},
			},
			expectedTasks: make(map[int][]Task),
		},
		"multiple tasks same key": {
			maxDays: 100,
			tasks: []*testTask{
				{
					id:       "b",
					daysFrom: 4,
				},
				{
					id:       "a",
					daysFrom: 4,
				},
			},
			expectedTasks: map[int][]Task{
				4: {
					&testTask{
						id:       "b",
						daysFrom: 4,
					},
					&testTask{
						id:       "a",
						daysFrom: 4,
					},
				},
			},
		},
		"multiple tasks same key beyond maxDays": {
			maxDays: 3,
			tasks: []*testTask{
				{
					id:       "b",
					daysFrom: 4,
				},
				{
					id:       "a",
					daysFrom: 4,
				},
			},
			expectedTasks: make(map[int][]Task),
		},
		"multiple tasks different key": {
			maxDays: 100,
			tasks: []*testTask{
				{
					id:       "a",
					daysFrom: 8,
				},
				{
					id:       "c",
					daysFrom: 5,
				},
				{
					id:       "b",
					daysFrom: 5,
				},
			},
			expectedTasks: map[int][]Task{
				5: {
					&testTask{
						id:       "b",
						daysFrom: 5,
					},
					&testTask{
						id:       "c",
						daysFrom: 5,
					},
				},
				8: {
					&testTask{
						id:       "a",
						daysFrom: 8,
					},
				},
			},
		},
		"multiple tasks different key with one beyond maxDays": {
			maxDays: 6,
			tasks: []*testTask{
				{
					id:       "a",
					daysFrom: 8,
				},
				{
					id:       "c",
					daysFrom: 5,
				},
				{
					id:       "b",
					daysFrom: 5,
				},
			},
			expectedTasks: map[int][]Task{
				5: {
					&testTask{
						id:       "b",
						daysFrom: 5,
					},
					&testTask{
						id:       "c",
						daysFrom: 5,
					},
				},
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			p := NewProcessor(time.Now(), test.maxDays, make(chan Task), make(chan struct{}))
			for _, tsk := range test.tasks {
				p.add(tsk)
			}
			assertEqualTestTaskMap(t, test.expectedTasks, p.tasks)
		})
	}
}
