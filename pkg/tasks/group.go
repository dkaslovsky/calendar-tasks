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

// Add adds tasks from a channel
func (g *Grouper) Add(ch <-chan Task, done <-chan struct{}, n int) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		nDone := 0
		for nDone < n {
			select {
			case t := <-ch:
				g.add(t)
			case <-done:
				nDone++
			}
		}
		// drain remaining tasks
		for {
			select {
			case t := <-ch:
				g.add(t)
			default:
				return
			}
		}
	}()
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
