package tasks

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

const delim = ":"

type rawLine struct {
	date string
	text string
}

// func load(ctx context.Context, fileName string, newTask func(*rawLine) (Task, error), taskCh chan Task, done chan struct{}) error {
// 	f, err := os.Open(filepath.Clean(fileName))
// 	if err != nil {
// 		return fmt.Errorf("failed to load file: %v", err)
// 	}

// 	go scan(f, newTask, taskCh, done)
// 	return nil
// }

// func load(l *Loader, fileName string, newTask func(*rawLine) (Task, error)) error {
// 	f, err := os.Open(filepath.Clean(fileName))
// 	if err != nil {
// 		return fmt.Errorf("failed to load file: %v", err)
// 	}

// 	l.eg.Go(func() error {
// 		return scan(l.ctx, f, newTask, l.Ch)
// 	})
// 	return nil
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

// func scan(r io.ReadCloser, newTask func(*rawLine) (Task, error), taskCh chan Task, done chan struct{}) error {
// 	defer r.Close()
// 	nTasks := 0

// 	scanner := bufio.NewScanner(r)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if strings.ReplaceAll(line, " ", "") == "" {
// 			continue
// 		}
// 		rl, err := loadLine(line)
// 		if err != nil {
// 			return fmt.Errorf("failed to load line: %v", err)
// 		}
// 		t, err := newTask(rl)
// 		if err != nil {
// 			return fmt.Errorf("failed to parse line: %v", err)
// 		}

// 		taskCh <- t
// 		nTasks++
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return err
// 	}
// 	if nTasks == 0 {
// 		return errors.New("failed to load any tasks")
// 	}

// 	done <- struct{}{}
// 	return nil
// }

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
