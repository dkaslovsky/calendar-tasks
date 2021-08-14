package tasks

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const delim = ":"

type rawLine struct {
	date string
	text string
}

func load(fileName string, newTask func(*rawLine) (Task, error)) (<-chan Task, error) {
	out := make(chan Task, 100)

	f, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		return out, fmt.Errorf("failed to load file: %v", err)
	}

	go scan(f, newTask, out)
	return out, nil
}

func scan(r io.ReadCloser, newTask func(*rawLine) (Task, error), out chan Task) error {
	defer close(out)
	defer r.Close()
	nTasks := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
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

		out <- t
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
