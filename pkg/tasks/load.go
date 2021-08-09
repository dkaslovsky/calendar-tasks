package tasks

import (
	"fmt"
	"os"
	"strings"
)

const delim = ":"

type rawLine struct {
	date string
	text string
}

func Load(fileName string, newTask func(*rawLine) (Task, error)) ([]Task, error) {
	tasks := []Task{}

	b, err := os.ReadFile(fileName)
	if err != nil {
		return tasks, fmt.Errorf("failed to load file: %v", err)
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		rl, err := loadLine(line)
		if err != nil {
			return tasks, fmt.Errorf("failed to load file: %v", err)
		}
		t, err := newTask(rl)
		if err != nil {
			return tasks, fmt.Errorf("failed to parse line: %v", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func loadLine(line string) (*rawLine, error) {
	parts := strings.Split(line, delim)
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
