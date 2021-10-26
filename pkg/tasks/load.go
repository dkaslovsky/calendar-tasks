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

	weekly  []string
	monthly []string
	annual  []string
	single  []string

	ctx context.Context
	eg  *errgroup.Group
}

// NewLoader constructs a Loader
func NewLoader(ch chan Task, done chan struct{}) *Loader {
	eg, ctx := errgroup.WithContext(context.Background())
	return &Loader{
		ch:   ch,
		done: done,

		weekly:  []string{},
		monthly: []string{},
		annual:  []string{},

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

// AddAnnualSource adds the name of a source file from which annual tasks are loaded
func (l *Loader) AddAnnualSource(s ...string) {
	l.annual = append(l.annual, s...)
}

// AddSingleSource adds the name of a source file from which single tasks are loaded
func (l *Loader) AddSingleSource(s ...string) {
	l.annual = append(l.annual, s...)
}

// Start launches the goroutines that load each task type
func (l *Loader) Start() error {
	defer func() {
		l.done <- struct{}{}
	}()

	// start one worker for each type of task (weekly, monthly, annual, single)
	numWorkers := 4
	weeklyCh := make(chan string, len(l.weekly))
	l.eg.Go(func() error {
		return l.scan(weeklyCh, newWeeklyTask)
	})
	monthlyCh := make(chan string, len(l.monthly))
	l.eg.Go(func() error {
		return l.scan(monthlyCh, newMonthlyTask)
	})
	annualCh := make(chan string, len(l.annual))
	l.eg.Go(func() error {
		return l.scan(annualCh, newAnnualTask)
	})
	singleCh := make(chan string, len(l.single))
	l.eg.Go(func() error {
		return l.scan(singleCh, newSingleTask)
	})

	// send each file on the appropriate channel to be processed
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
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
		for _, fp := range l.annual {
			annualCh <- fp
		}
	}()
	go func() {
		defer wg.Done()
		for _, fp := range l.single {
			singleCh <- fp
		}
	}()
	wg.Wait()

	close(weeklyCh)
	close(monthlyCh)
	close(annualCh)
	close(singleCh)
	return l.eg.Wait()
}

type newTaskF func(*sources.RawTask) (Task, error)

// scan is a worker that loads the tasks from file names it receives on a channel
func (l *Loader) scan(fileCh <-chan string, newTask newTaskF) error {
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

func scan(ctx context.Context, r io.ReadCloser, newTask newTaskF, taskCh chan Task) error {
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
		rawTasks, err := sources.ParseLine(line)
		if err != nil {
			return fmt.Errorf("failed to load line: %v", err)
		}
		for _, rawTask := range rawTasks {
			t, err := newTask(rawTask)
			if err != nil {
				return fmt.Errorf("failed to parse line: %v", err)
			}

			taskCh <- t
			nTasks++
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	if nTasks == 0 {
		return errors.New("failed to load any tasks")
	}
	return nil
}

func newWeeklyTask(r *sources.RawTask) (Task, error) {
	return sources.NewWeekly(r)
}

func newMonthlyTask(r *sources.RawTask) (Task, error) {
	return sources.NewMonthly(r)
}

func newAnnualTask(r *sources.RawTask) (Task, error) {
	return sources.NewAnnual(r)
}

func newSingleTask(r *sources.RawTask) (Task, error) {
	return sources.NewSingle(r)
}
