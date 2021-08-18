package tasks

import (
	"sync"
	"time"
)

type Consumer struct {
	now     time.Time
	maxDays int

	in   <-chan Task
	done <-chan struct{}

	wg    *sync.WaitGroup
	lock  sync.RWMutex
	tasks map[int][]Task
}

func NewConsumer(now time.Time, maxDays int, in <-chan Task, done <-chan struct{}) *Consumer {
	return &Consumer{
		now:     now,
		maxDays: maxDays,

		in:   in,
		done: done,

		wg:    &sync.WaitGroup{},
		lock:  sync.RWMutex{},
		tasks: make(map[int][]Task),
	}
}

func (c *Consumer) Start() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case t := <-c.in:
				c.add(t)
			case <-c.done:
				c.drain()
				return
			}
		}
	}()
}

func (c *Consumer) Wait() {
	c.wg.Wait()
}

func (c *Consumer) Tasks() map[int][]Task {
	return c.tasks
}

func (c *Consumer) drain() {
	for {
		select {
		case t := <-c.in:
			c.add(t)
		default:
			return
		}
	}
}

func (c *Consumer) add(t Task) {
	days := t.DaysFrom(c.now)
	if days > c.maxDays {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exists := c.tasks[days]; !exists {
		c.tasks[days] = []Task{}
	}
	c.tasks[days] = append(c.tasks[days], t)
}
