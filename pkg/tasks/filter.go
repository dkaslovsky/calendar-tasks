package tasks

import "time"

// TODO: make thread safe?
type Filter struct {
	now   time.Time
	tasks map[int][]Task
}

func NewFilter(now time.Time) *Filter {
	return &Filter{
		now:   now,
		tasks: make(map[int][]Task),
	}
}

func (f *Filter) Add(t Task) {
	days := t.DaysFrom(f.now)
	if _, exists := f.tasks[days]; !exists {
		f.tasks[days] = []Task{}
	}
	f.tasks[days] = append(f.tasks[days], t)
}

func (f *Filter) GetTasksFlat(nDays int) []Task {
	tasks := []Task{}
	for day := 0; day < nDays; day++ {
		t, ok := f.tasks[day]
		if !ok {
			continue
		}
		tasks = append(tasks, t...)
	}
	return tasks
}

func (f *Filter) GetTasksGrouped(nDays int) map[int][]Task {
	ff := NewFilter(f.now)
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
