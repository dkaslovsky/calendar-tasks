package tasks

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks/sources"
	"golang.org/x/sync/errgroup"
)

// Loader loads raw tasks to be sent for processing
type Loader struct {
	ch   chan Task
	done chan struct{}

	weekly    []string
	monthly   []string
	multiDate []string

	ctx context.Context
	eg  *errgroup.Group
}

// NewLoader constructs a Loader
func NewLoader(ch chan Task, done chan struct{}) *Loader {
	// ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(context.Background())
	return &Loader{
		ch:   ch,
		done: done,

		weekly:    []string{},
		monthly:   []string{},
		multiDate: []string{},

		ctx: ctx,
		eg:  eg,
	}
}

// AddWeeklySource adds the name of a source file from which weekly tasks are loaded
func (l *Loader) AddWeeklySource(s ...string) {
	l.weekly = append(l.weekly, s...)
}

// AddMonthlySource adds the name of a source file from which monthly tasks are loaded
func (l *Loader) AddMonthlySource(s ...string) {
	l.monthly = append(l.monthly, s...)
}

// AddMultiDateSource adds the name of a source file from which multiDate tasks are loaded
func (l *Loader) AddMultiDateSource(s ...string) {
	l.multiDate = append(l.multiDate, s...)
}

// Start launches the goroutines that load each task type
func (l *Loader) Start() error {
	defer func() {
		l.done <- struct{}{}
	}()

	// start one worker for each type of task (weekly, monthly, multiDate)
	weeklyCh := make(chan string, len(l.weekly))
	l.eg.Go(func() error {
		return l.scan(weeklyCh, newWeeklyTask)
	})
	monthlyCh := make(chan string, len(l.monthly))
	l.eg.Go(func() error {
		return l.scan(monthlyCh, newMonthlyTask)
	})
	multiDateCh := make(chan string, len(l.multiDate))
	l.eg.Go(func() error {
		return l.scan(multiDateCh, newMultiDateTask)
	})

	// send each file on the appropriate channel to be processed
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		for _, fp := range l.weekly {
			weeklyCh <- fp
		}
	}()
	go func() {
		defer wg.Done()
		for _, fp := range l.monthly {
			monthlyCh <- fp
		}
	}()
	go func() {
		defer wg.Done()
		for _, fp := range l.multiDate {
			multiDateCh <- fp
		}
	}()
	wg.Wait()

	close(weeklyCh)
	close(monthlyCh)
	close(multiDateCh)
	return l.eg.Wait()
}

// scan is a worker that loads the tasks from file names it receives on a channel
func (l *Loader) scan(fileCh <-chan string, newTask func(*sources.RawLine) (Task, error)) error {
	for fp := range fileCh {
		f, err := os.Open(filepath.Clean(fp))
		if err != nil {
			return err
		}
		err = scan(l.ctx, f, newTask, l.ch)
		if err != nil {
			return err
		}
	}
	return nil
}

func scan(ctx context.Context, r io.ReadCloser, newTask func(*sources.RawLine) (Task, error), taskCh chan Task) error {
	defer r.Close() //nolint
	nTasks := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()
		if strings.ReplaceAll(line, " ", "") == "" {
			continue
		}
		rl, err := sources.LoadLine(line)
		if err != nil {
			return fmt.Errorf("failed to load line: %v", err)
		}
		t, err := newTask(rl)
		if err != nil {
			return fmt.Errorf("failed to parse line: %v", err)
		}

		taskCh <- t
		nTasks++
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	if nTasks == 0 {
		return errors.New("failed to load any tasks")
	}
	return nil
}

func newWeeklyTask(r *sources.RawLine) (Task, error) {
	return sources.NewWeekly(r)
}

func newMonthlyTask(r *sources.RawLine) (Task, error) {
	return sources.NewMonthly(r)
}

func newMultiDateTask(r *sources.RawLine) (Task, error) {
	return sources.NewMultiDate(r)
}
