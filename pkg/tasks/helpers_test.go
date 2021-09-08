package tasks

import (
	"sort"
	"testing"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks/sources"
)

type testTask struct {
	id       string
	daysFrom int
}

func (tt *testTask) DaysFrom(t time.Time) int { return tt.daysFrom }

func (tt *testTask) String() string { return "" }

func (tt *testTask) equal(other *testTask) bool { return tt.id == other.id }

func newTestTask(rl *sources.RawTask) (Task, error) {
	return &testTask{
		id:       rl.Text,
		daysFrom: 1,
	}, nil
}

func assertEqualTestTaskMap(t *testing.T, expected, actual map[int][]Task) {
	for day, etsks := range expected {
		atsks, ok := actual[day]
		if !ok {
			t.Fatalf("result missing task key %d", day)
		}
		assertEqualTestTaskSlice(t, etsks, atsks)
	}
}

func assertEqualTestTaskSlice(t *testing.T, expected, actual []Task) {
	if len(actual) != len(expected) {
		t.Fatalf("result number of tasks %d not equal to expected number of tasks %d", len(actual), len(expected))
	}

	etasks := []*testTask{}
	atasks := []*testTask{}

	for i := 0; i < len(expected); i++ {
		etsk, ok := expected[i].(*testTask)
		if !ok {
			t.Fatalf("type assertion on expected task failed")
		}
		etasks = append(etasks, etsk)

		atsk, ok := actual[i].(*testTask)
		if !ok {
			t.Fatalf("type assertion on result task failed")
		}
		atasks = append(atasks, atsk)
	}

	sort.Slice(etasks, func(i, j int) bool {
		return etasks[i].id > etasks[j].id
	})
	sort.Slice(atasks, func(i, j int) bool {
		return atasks[i].id > atasks[j].id
	})

	for i := 0; i < len(etasks); i++ {
		if !etasks[i].equal(atasks[i]) {
			t.Fatalf("result task with id '%s' not equal to expected task with id '%s'", atasks[i].id, etasks[i].id)
		}
	}
}
