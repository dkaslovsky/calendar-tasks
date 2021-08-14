package tasks

import (
	"sync"
	"time"
)

// Grouper groups tasks by the number of days until they are to occur
type Grouper struct {
	now time.Time

	lock  sync.RWMutex
	tasks map[int][]Task
}

// NewGrouper constructs a Grouper
func NewGrouper(now time.Time) *Grouper {
	return &Grouper{
		now:   now,
		tasks: make(map[int][]Task),
	}
}

// Add adds tasks from one or more channels to a Grouper
func (g *Grouper) Add(chs ...<-chan Task) {
	wg := &sync.WaitGroup{}
	for i, ch := range chs {
		wg.Add(1)
		go func(ch <-chan Task, ii int) {
			defer wg.Done()
			for t := range ch {
				g.add(t)
			}
		}(ch, i)
	}
	wg.Wait()
}

func (g *Grouper) add(t Task) {
	g.lock.Lock()
	defer g.lock.Unlock()
	days := t.DaysFrom(g.now)
	if _, exists := g.tasks[days]; !exists {
		g.tasks[days] = []Task{}
	}
	g.tasks[days] = append(g.tasks[days], t)
}

// Filter returns groups of tasks to occur within a specified number of days
func (g *Grouper) Filter(nDays int) map[int][]Task {
	g.lock.RLock()
	defer g.lock.RUnlock()
	gg := NewGrouper(g.now)
	// start the loop at today (0) and include nDays
	for day := 0; day <= nDays; day++ {
		tasks, ok := g.tasks[day]
		if !ok {
			continue
		}
		for _, t := range tasks {
			gg.add(t)
		}
	}
	return gg.tasks
}
