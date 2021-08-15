package tasks

import (
	"context"
	"io"
	"os"

	"golang.org/x/sync/errgroup"
)

type Loader struct {
	Ch chan Task

	ctx    context.Context
	cancel context.CancelFunc
	eg     *errgroup.Group
}

func NewLoader() *Loader {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	return &Loader{
		Ch: make(chan Task, 100), // TODO: make channel buffer size configurable

		ctx:    ctx,
		cancel: cancel,
		eg:     eg,
	}
}

func (l *Loader) Load(weekly []string, monthly []string, recurring []string) error {
	// pool of workers, one for each of weekly, monthly, recurring
	weeklyCh := make(chan io.ReadCloser)
	l.eg.Go(func() error {
		return l.scanWorker(weeklyCh, newWeekly)
	})
	monthlyCh := make(chan io.ReadCloser)
	l.eg.Go(func() error {
		return l.scanWorker(monthlyCh, newMonthly)
	})
	recurringCh := make(chan io.ReadCloser)
	l.eg.Go(func() error {
		return l.scanWorker(recurringCh, newRecurring)
	})

	for _, w := range weekly {
		f, err := os.Open(w)
		if err != nil {
			l.cancel()
			return err
		}
		weeklyCh <- f
	}

	for _, m := range monthly {
		f, err := os.Open(m)
		if err != nil {
			l.cancel()
			return err
		}
		monthlyCh <- f
	}

	for _, r := range recurring {
		f, err := os.Open(r)
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

func (l *Loader) Wait() error {
	defer close(l.Ch)
	return l.eg.Wait()
}

func (l *Loader) scanWorker(rcs <-chan io.ReadCloser, newTask func(*rawLine) (Task, error)) error {
	for r := range rcs {
		err := scan(l.ctx, r, newTask, l.Ch)
		if err != nil {
			return err
		}
	}
	return nil
}
