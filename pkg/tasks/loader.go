package tasks

import (
	"context"
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
	for _, w := range weekly {
		f, err := os.Open(w)
		if err != nil {
			l.cancel()
			return err
		}

		l.eg.Go(func() error {
			return scan(l.ctx, f, newWeekly, l.Ch)
		})
	}

	for _, m := range monthly {
		f, err := os.Open(m)
		if err != nil {
			l.cancel()
			return err
		}

		l.eg.Go(func() error {
			return scan(l.ctx, f, newMonthly, l.Ch)
		})
	}

	for _, r := range recurring {
		f, err := os.Open(r)
		if err != nil {
			l.cancel()
			return err
		}

		l.eg.Go(func() error {
			return scan(l.ctx, f, newRecurring, l.Ch)
		})
	}

	return nil
}

func (l *Loader) Wait() error {
	defer close(l.Ch)
	return l.eg.Wait()
}
