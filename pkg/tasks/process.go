package tasks

import (
	"sync"
	"time"
)

// Processor groups and filters tasks
type Processor struct {
	now     time.Time
	maxDays int

	in   <-chan Task
	done <-chan struct{}

	wg    *sync.WaitGroup
	lock  sync.RWMutex
	tasks map[int][]Task
}

// NewProcessor constructs a Processor
func NewProcessor(now time.Time, maxDays int, in <-chan Task, done <-chan struct{}) *Processor {
	return &Processor{
		now:     now,
		maxDays: maxDays,

		in:   in,
		done: done,

		wg:    &sync.WaitGroup{},
		lock:  sync.RWMutex{},
		tasks: make(map[int][]Task),
	}
}

// Start launches the Processor goroutine
func (p *Processor) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case t := <-p.in:
				p.add(t)
			case <-p.done:
				p.drain()
				return
			}
		}
	}()
}

// Wait waits for the Processor goroutine to finish
func (p *Processor) Wait() {
	p.wg.Wait()
}

// Tasks returnns processed tasks
func (p *Processor) Tasks() map[int][]Task {
	return p.tasks
}

func (p *Processor) drain() {
	for {
		select {
		case t := <-p.in:
			p.add(t)
		default:
			return
		}
	}
}

func (p *Processor) add(t Task) {
	days := t.DaysFrom(p.now)
	if days > p.maxDays {
		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	if _, exists := p.tasks[days]; !exists {
		p.tasks[days] = []Task{}
	}
	p.tasks[days] = append(p.tasks[days], t)
}
