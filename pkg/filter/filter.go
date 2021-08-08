package filter

import "time"

type task interface {
	DaysFrom(time.Time) int
}

// TODO: make thread safe?
type Filter struct {
	now   time.Time
	tasks map[int][]task
}

func New(now time.Time) *Filter {
	return &Filter{
		now:   now,
		tasks: make(map[int][]task),
	}
}

func (f *Filter) Add(t task) {
	days := t.DaysFrom(f.now)
	if _, exists := f.tasks[days]; !exists {
		f.tasks[days] = []task{}
	}
	f.tasks[days] = append(f.tasks[days], t)
}

func (f *Filter) GetTasksFlat(nDays int) []task {
	tasks := []task{}
	for day := 0; day < nDays; day++ {
		t, ok := f.tasks[day]
		if !ok {
			continue
		}
		tasks = append(tasks, t...)
	}
	return tasks
}

func (f *Filter) GetTasksGrouped(nDays int) map[int][]task {
	ff := New(f.now)
	for day := 0; day < nDays; day++ {
		tasks, ok := f.tasks[day]
		if !ok {
			continue
		}
		for _, t := range tasks {
			ff.Add(t)
		}
	}
	return ff.tasks
}
