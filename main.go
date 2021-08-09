package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/filter"
	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

func main() {
	now := time.Now()
	f := filter.New(now)

	weekly, err := tasks.LoadWeekly(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range weekly {
		f.Add(d)
	}

	monthly, err := tasks.LoadMonthly(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range monthly {
		f.Add(m)
	}

	repeatedMonthly, err := tasks.LoadRepeatedMonthly(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range repeatedMonthly {
		f.Add(m)
	}

	daily, err := tasks.LoadRepeatedMonthly(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range daily {
		f.Add(m)
	}

	n, _ := strconv.Atoi(os.Args[5])

	tasksByDay := f.GetTasksGrouped(n)
	for day := 0; day <= n; day++ {
		tasks, ok := tasksByDay[day]
		if !ok {
			continue
		}
		fmt.Printf("Day = %d\n", day)
		for _, task := range tasks {
			fmt.Println(task)
		}
	}
}
