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
	Ch chan Task

	weekly    []string
	monthly   []string
	recurring []string

	ctx    context.Context
	cancel context.CancelFunc
	eg     *errgroup.Group
}

func NewLoader() *Loader {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	return &Loader{
		Ch: make(chan Task, 100), // TODO: make channel buffer size configurable

		weekly:    []string{},
		monthly:   []string{},
		recurring: []string{},

		ctx:    ctx,
		cancel: cancel,
		eg:     eg,
	}
}

func (l *Loader) AddWeekly(s ...string) {
	l.weekly = append(l.weekly, s...)
}

func (l *Loader) AddMonthly(s ...string) {
	l.monthly = append(l.monthly, s...)
}

func (l *Loader) AddRecurring(s ...string) {
	l.recurring = append(l.recurring, s...)
}

func (l *Loader) Start() error {
	// one worker for each type of task (weekly, monthly, recurring)
	weeklyCh := make(chan io.ReadCloser)
	l.eg.Go(func() error {
		return l.scan(weeklyCh, newWeekly)
	})
	monthlyCh := make(chan io.ReadCloser)
	l.eg.Go(func() error {
		return l.scan(monthlyCh, newMonthly)
	})
	recurringCh := make(chan io.ReadCloser)
	l.eg.Go(func() error {
		return l.scan(recurringCh, newRecurring)
	})

	// process each task file
	for _, fp := range l.weekly {
		f, err := os.Open(filepath.Clean(fp))
		if err != nil {
			l.cancel()
			return err
		}
		weeklyCh <- f
	}
	for _, fp := range l.monthly {
		f, err := os.Open(filepath.Clean(fp))
		if err != nil {
			l.cancel()
			return err
		}
		monthlyCh <- f
	}
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
	return nil
}

func (l *Loader) Close() error {
	defer close(l.Ch)
	return l.eg.Wait()
}

func (l *Loader) scan(rcs <-chan io.ReadCloser, newTask func(*rawLine) (Task, error)) error {
	for r := range rcs {
		err := scan(l.ctx, r, newTask, l.Ch)
		if err != nil {
			return err
		}
	}
	return nil
}

// type scanPool struct {
// 	weeklyCh    chan io.ReadCloser
// 	monthlyCh   chan io.ReadCloser
// 	recurringCh chan io.ReadCloser

// 	loader *Loader
// }

// func newScanPool(loader *Loader) *scanPool {
// 	return &scanPool{
// 		weeklyCh:    make(chan io.ReadCloser),
// 		monthlyCh:   make(chan io.ReadCloser),
// 		recurringCh: make(chan io.ReadCloser),

// 		loader: loader,
// 	}
// }

// func (s *scanPool) Start() {
// 	s.loader.eg.Go(func() error {
// 		return s.loader.scan(s.weeklyCh, newWeekly)
// 	})
// 	s.loader.eg.Go(func() error {
// 		return s.loader.scan(s.monthlyCh, newMonthly)
// 	})
// 	s.loader.eg.Go(func() error {
// 		return s.loader.scan(s.recurringCh, newRecurring)
// 	})
// }

// func (s *scanPool) Close() {
// 	close(s.weeklyCh)
// 	close(s.monthlyCh)
// 	close(s.recurringCh)
// }

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
