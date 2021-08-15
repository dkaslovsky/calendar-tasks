package tasks

import (
	"sync"
	"time"
)

type Consumer struct {
	now     time.Time
	maxDays int

	in         <-chan Task
	loadCloser func() error

	lock  sync.RWMutex
	tasks map[int][]Task
}

func NewConsumer(now time.Time, maxDays int, in <-chan Task, loadCloser func() error) *Consumer {
	return &Consumer{
		now:     now,
		maxDays: maxDays,

		in:         in,
		loadCloser: loadCloser,

		tasks: make(map[int][]Task),
	}
}

func (c *Consumer) Start() error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range c.in {
			c.add(t)
		}
	}()

	err := c.loadCloser()
	if err != nil {
		return err
	}

	wg.Wait()
	return nil
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

func (c *Consumer) Get() map[int][]Task {
	return c.tasks
}
