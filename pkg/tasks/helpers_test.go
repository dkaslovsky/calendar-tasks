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

func (tt *testTask) equal(other *testTask) bool { return tt.id == other.id }

func assertEqualTestTaskMap(t *testing.T, expected, actual map[int][]Task) {
	for day, etsks := range expected {
		atsks, ok := actual[day]
		if !ok {
			t.Fatalf("result missing task key %d", day)
		}

		assertEqualTestTaskSlice(t, etsks, atsks)
		// if len(atsks) != len(etsks) {
		// 	t.Fatalf("result number of tasks %d not equal to expected number of tasks %d for key %d", len(atsks), len(etsks), day)
		// }

		// etasks := []*testTask{}
		// atasks := []*testTask{}
		// for i := 0; i < len(etsks); i++ {
		// 	etsk, ok := etsks[i].(*testTask)
		// 	if !ok {
		// 		t.Fatalf("type assertion on expected task failed")
		// 	}
		// 	etasks = append(etasks, etsk)

		// 	atsk, ok := atsks[i].(*testTask)
		// 	if !ok {
		// 		t.Fatalf("type assertion on grouped task failed")
		// 	}
		// 	atasks = append(atasks, atsk)
		// }
		// sort.Slice(etasks, func(i, j int) bool {
		// 	return etasks[i].id > etasks[j].id
		// })
		// sort.Slice(atasks, func(i, j int) bool {
		// 	return atasks[i].id > atasks[j].id
		// })

		// for i := 0; i < len(etsks); i++ {
		// 	if !etasks[i].equal(atasks[i]) {
		// 		t.Fatalf("result task with id '%s' not equal to expected task with id '%s'", atasks[i].id, atasks[i].id)
		// 	}
		// }
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

func assertShouldError(t *testing.T, shouldErr bool, err error) {
	if shouldErr {
		if err == nil {
			t.Fatal("expected error but result err is nil")
		}
		return
	}
	if !shouldErr && err != nil {
		t.Fatalf("expected nil error but result err is %v", err)
	}
}
