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

func load(fileName string, newTask func(*rawLine) (Task, error)) ([]Task, error) {
	f, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		return []Task{}, fmt.Errorf("failed to load file: %v", err)
	}
	defer f.Close()
	return scan(f, newTask)
}

func scan(r io.Reader, newTask func(*rawLine) (Task, error)) ([]Task, error) {
	tasks := []Task{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.ReplaceAll(line, " ", "") == "" {
			continue
		}
		rl, err := loadLine(line)
		if err != nil {
			return tasks, fmt.Errorf("failed to load line: %v", err)
		}
		t, err := newTask(rl)
		if err != nil {
			return tasks, fmt.Errorf("failed to parse line: %v", err)
		}
		tasks = append(tasks, t)
	}
	if err := scanner.Err(); err != nil {
		return tasks, err
	}
	if len(tasks) == 0 {
		return tasks, errors.New("failed to load any tasks")
	}
	return tasks, scanner.Err()
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
