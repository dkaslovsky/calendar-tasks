package tasks

import (
	"testing"
	"time"
)

func TestAddOne(t *testing.T) {
	tests := map[string]struct {
		tasks         []*testTask
		expectedTasks map[int][]Task
	}{
		"empty": {
			tasks:         []*testTask{},
			expectedTasks: make(map[int][]Task),
		},
		"single task": {
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
		"multiple tasks same key": {
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
		"multiple tasks different key": {
			tasks: []*testTask{
				{
					id:       "a",
					daysFrom: 4,
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
				4: {
					&testTask{
						id:       "a",
						daysFrom: 4,
					},
				},
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			g := NewGrouper(time.Now())
			for _, tsk := range test.tasks {
				g.add(tsk)
			}
			assertEqualTestTaskMap(t, test.expectedTasks, g.tasks)
		})
	}
}

func TestFilter(t *testing.T) {
	tasks := map[int][]Task{
		1: {
			&testTask{
				id:       "c",
				daysFrom: 1,
			},
			&testTask{
				id:       "b",
				daysFrom: 1,
			},
		},
		2: {
			&testTask{
				id:       "d",
				daysFrom: 2,
			},
			&testTask{
				id:       "e",
				daysFrom: 2,
			},
			&testTask{
				id:       "f",
				daysFrom: 2,
			},
		},
		4: {
			&testTask{
				id:       "a",
				daysFrom: 4,
			},
		},
	}

	tests := map[string]struct {
		nDays         int
		tasks         map[int][]Task
		expectedTasks map[int][]Task
	}{
		"filter today (with tasks)": {
			nDays: 0,
			tasks: map[int][]Task{
				0: {
					&testTask{
						id:       "x",
						daysFrom: 0,
					},
					&testTask{
						id:       "y",
						daysFrom: 0,
					},
				},
			},
			expectedTasks: map[int][]Task{
				0: {
					&testTask{
						id:       "x",
						daysFrom: 0,
					},
					&testTask{
						id:       "y",
						daysFrom: 0,
					},
				},
			},
		},
		"filter today (no tasks)": {
			nDays:         0,
			tasks:         tasks,
			expectedTasks: make(map[int][]Task),
		},
		"filter for next 1 day": {
			nDays: 1,
			tasks: tasks,
			expectedTasks: map[int][]Task{
				1: {
					&testTask{
						id:       "c",
						daysFrom: 1,
					},
					&testTask{
						id:       "b",
						daysFrom: 1,
					},
				},
			},
		},
		"filter for next 2 days": {
			nDays: 2,
			tasks: tasks,
			expectedTasks: map[int][]Task{
				1: {
					&testTask{
						id:       "c",
						daysFrom: 1,
					},
					&testTask{
						id:       "b",
						daysFrom: 1,
					},
				},
				2: {
					&testTask{
						id:       "d",
						daysFrom: 2,
					},
					&testTask{
						id:       "e",
						daysFrom: 2,
					},
					&testTask{
						id:       "f",
						daysFrom: 2,
					},
				},
			},
		},
		"filter for next 3 days": {
			nDays: 3,
			tasks: tasks,
			expectedTasks: map[int][]Task{
				1: {
					&testTask{
						id:       "c",
						daysFrom: 1,
					},
					&testTask{
						id:       "b",
						daysFrom: 1,
					},
				},
				2: {
					&testTask{
						id:       "d",
						daysFrom: 2,
					},
					&testTask{
						id:       "e",
						daysFrom: 2,
					},
					&testTask{
						id:       "f",
						daysFrom: 2,
					},
				},
			},
		},
		"filter for next 10 days": {
			nDays: 10,
			tasks: tasks,
			expectedTasks: map[int][]Task{
				1: {
					&testTask{
						id:       "c",
						daysFrom: 1,
					},
					&testTask{
						id:       "b",
						daysFrom: 1,
					},
				},
				2: {
					&testTask{
						id:       "d",
						daysFrom: 2,
					},
					&testTask{
						id:       "e",
						daysFrom: 2,
					},
					&testTask{
						id:       "f",
						daysFrom: 2,
					},
				},
				4: {
					&testTask{
						id:       "a",
						daysFrom: 4,
					},
				},
			},
		},
		"filter for next 1000 days": {
			nDays: 1000,
			tasks: tasks,
			expectedTasks: map[int][]Task{
				1: {
					&testTask{
						id:       "c",
						daysFrom: 1,
					},
					&testTask{
						id:       "b",
						daysFrom: 1,
					},
				},
				2: {
					&testTask{
						id:       "d",
						daysFrom: 2,
					},
					&testTask{
						id:       "e",
						daysFrom: 2,
					},
					&testTask{
						id:       "f",
						daysFrom: 2,
					},
				},
				4: {
					&testTask{
						id:       "a",
						daysFrom: 4,
					},
				},
			},
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			g := NewGrouper(time.Now()) // time is not relevant to this test
			g.tasks = test.tasks

			result := g.Filter(test.nDays)
			assertEqualTestTaskMap(t, test.expectedTasks, result)
		})
	}
}
