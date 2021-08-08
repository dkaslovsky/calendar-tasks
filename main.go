package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/filter"
	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

func main() {
	now := time.Now()
	rs := days()
	//rs := months()

	f := filter.New(now)
	for _, r := range rs {
		f.Add(r)
	}

	for day, tasks := range f.GetTasksGrouped(10) {
		fmt.Printf("Day = %d\n", day)
		for _, task := range tasks {
			fmt.Println(task)
		}
	}
}

func days() []*tasks.Daily {
	fileName := os.Args[1]
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(b), "\n")

	ds := []*tasks.Daily{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		d, err := tasks.NewDaily(line)
		if err != nil {
			log.Fatal(err)
		}
		ds = append(ds, d)
	}

	return ds
}

func months() []*tasks.Monthly {
	fileName := os.Args[1]
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(b), "\n")

	ms := []*tasks.Monthly{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		m, err := tasks.NewMonthly(line)
		if err != nil {
			log.Fatal(err)
		}
		ms = append(ms, m)
	}

	return ms
}
