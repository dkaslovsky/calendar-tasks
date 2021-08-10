package tasks

import (
	"sort"
	"testing"
	"time"
)

type testTask struct {
	id       string
	daysFrom int
}

func (tt *testTask) DaysFrom(t time.Time) int { return tt.daysFrom }

func (tt *testTask) String() string { return "" }

func (tt *testTask) Equal(other *testTask) bool { return tt.id == other.id }

func assertEqualTestTaskMap(t *testing.T, expected, actual map[int][]Task) {
	for day, etsks := range expected {
		atsks, ok := actual[day]
		if !ok {
			t.Fatalf("result missing task key %d", day)
		}
		if len(atsks) != len(etsks) {
			t.Fatalf("result number of tasks %d not equal to expected number of tasks %d for key %d", len(atsks), len(etsks), day)
		}

		eTestTasks := []*testTask{}
		aTestTasks := []*testTask{}
		for i := 0; i < len(etsks); i++ {
			etsk, ok := etsks[i].(*testTask)
			if !ok {
				t.Fatalf("type assertion on expected task failed")
			}
			eTestTasks = append(eTestTasks, etsk)

			atsk, ok := atsks[i].(*testTask)
			if !ok {
				t.Fatalf("type assertion on grouped task failed")
			}
			aTestTasks = append(aTestTasks, atsk)
		}
		sort.Slice(eTestTasks, func(i, j int) bool {
			return eTestTasks[i].id > eTestTasks[j].id
		})
		sort.Slice(aTestTasks, func(i, j int) bool {
			return aTestTasks[i].id > aTestTasks[j].id
		})

		for i := 0; i < len(etsks); i++ {
			if !eTestTasks[i].Equal(aTestTasks[i]) {
				t.Fatalf("result task with id '%s' not equal to expected task with id '%s'", aTestTasks[i].id, aTestTasks[i].id)
			}
		}
	}
}

func TestAdd(t *testing.T) {
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
				g.Add(tsk)
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
