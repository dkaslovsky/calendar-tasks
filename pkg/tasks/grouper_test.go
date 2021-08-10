package tasks

import (
	"testing"
	"time"
)

type testTask struct {
	id       string
	daysFrom int
}

func (tt *testTask) DaysFrom(t time.Time) int { return tt.daysFrom }

func (tt *testTask) String() string { return "" }

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

			for day, etsks := range test.expectedTasks {
				gtsks, ok := g.tasks[day]
				if !ok {
					t.Fatalf("result missing task key %d", day)
				}
				if len(gtsks) != len(etsks) {
					t.Fatalf("result number of tasks %d not equal to expected number of tasks %d for key %d", len(gtsks), len(etsks), day)
				}
				for i, etsk := range etsks {
					etsk, ok := etsk.(*testTask)
					if !ok {
						t.Fatalf("bad test set up, type assertion on expected task failed")
					}
					gtsk, ok := gtsks[i].(*testTask)
					if !ok {
						t.Fatalf("type assertion on grouped task failed")
					}
					if gtsk.id != etsk.id {
						t.Fatalf("result task id '%s' not equal to expected task id '%s'", gtsk.id, etsk.id)
					}
				}
			}
		})
	}
}
