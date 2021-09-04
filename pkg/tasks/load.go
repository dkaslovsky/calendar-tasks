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

	"golang.org/x/sync/errgroup"
)

type Loader struct {
	ch   chan Task
	done chan struct{}

	weekly    []string
	monthly   []string
	recurring []string

	ctx    context.Context
	cancel context.CancelFunc
	eg     *errgroup.Group
}

func NewLoader(ch chan Task, done chan struct{}) *Loader {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	return &Loader{
		ch:   ch,
		done: done,

		weekly:    []string{},
		monthly:   []string{},
		recurring: []string{},

		ctx:    ctx,
		cancel: cancel,
		eg:     eg,
	}
}

func (l *Loader) AddWeeklySource(s ...string) {
	l.weekly = append(l.weekly, s...)
}

func (l *Loader) AddMonthlySource(s ...string) {
	l.monthly = append(l.monthly, s...)
}

func (l *Loader) AddRecurringSource(s ...string) {
	l.recurring = append(l.recurring, s...)
}

func (l *Loader) Start() error {
	defer func() {
		l.done <- struct{}{}
	}()

	// start one worker for each type of task (weekly, monthly, recurring)
	weeklyCh := make(chan io.ReadCloser, len(l.weekly))
	l.eg.Go(func() error {
		return l.scan(weeklyCh, newWeekly)
	})
	monthlyCh := make(chan io.ReadCloser, len(l.monthly))
	l.eg.Go(func() error {
		return l.scan(monthlyCh, newMonthly)
	})
	recurringCh := make(chan io.ReadCloser, len(l.recurring))
	l.eg.Go(func() error {
		return l.scan(recurringCh, newRecurring)
	})

	// process each weekly task file
	for _, fp := range l.weekly {
		f, err := os.Open(filepath.Clean(fp))
		if err != nil {
			l.cancel()
			return err
		}
		weeklyCh <- f
	}
	// process each monthly task file
	for _, fp := range l.monthly {
		f, err := os.Open(filepath.Clean(fp))
		if err != nil {
			l.cancel()
			return err
		}
		monthlyCh <- f
	}
	// process each recurring task file
	for _, fp := range l.recurring {
		f, err := os.Open(filepath.Clean(fp))
		if err != nil {
			l.cancel()
			return err
		}
		recurringCh <- f
	}

	close(weeklyCh)
	close(monthlyCh)
	close(recurringCh)
	return l.eg.Wait()
}

func (l *Loader) scan(rcs <-chan io.ReadCloser, newTask func(*rawLine) (Task, error)) error {
	for r := range rcs {
		err := scan(l.ctx, r, newTask, l.ch)
		if err != nil {
			return err
		}
	}
	return nil
}

func scan(ctx context.Context, r io.ReadCloser, newTask func(*rawLine) (Task, error), taskCh chan Task) error {
	defer r.Close()
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
		rl, err := loadLine(line)
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

const delim = ":"

type rawLine struct {
	date string
	text string
}

func loadLine(line string) (*rawLine, error) {
	parts := strings.SplitN(line, delim, 2)
	if len(parts) != 2 {
		return &rawLine{}, fmt.Errorf("invalid line [%s]", line)
	}

	r := &rawLine{
		date: parts[0],
		text: cleanText(parts[1]),
	}
	return r, nil
}

func cleanText(s string) string {
	return strings.TrimLeft(s, " ")
}
