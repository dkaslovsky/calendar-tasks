package tasks

import "time"

// TODO: make Grouper thread safe?

// Grouper groups tasks by the number of days until they are to occur
type Grouper struct {
	now   time.Time
	tasks map[int][]Task
}

// NewGrouper constructs a Grouper
func NewGrouper(now time.Time) *Grouper {
	return &Grouper{
		now:   now,
		tasks: make(map[int][]Task),
	}
}

// Add adds a task to a Grouper
func (g *Grouper) Add(t Task) {
	days := t.DaysFrom(g.now)
	if _, exists := g.tasks[days]; !exists {
		g.tasks[days] = []Task{}
	}
	g.tasks[days] = append(g.tasks[days], t)
}

// Filter returns groups of tasks to occur within a specified number of days
func (g *Grouper) Filter(nDays int) map[int][]Task {
	gg := NewGrouper(g.now)
	// start the loop at today (0) and include nDays
	for day := 0; day <= nDays; day++ {
		tasks, ok := g.tasks[day]
		if !ok {
			continue
		}
		for _, t := range tasks {
			gg.Add(t)
		}
	}
	return gg.tasks
}
