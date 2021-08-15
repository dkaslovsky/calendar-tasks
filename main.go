package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

func main() {
	now := time.Now()
	grouper := tasks.NewGrouper(now)

	n, _ := strconv.Atoi(os.Args[5])

	taskCh := make(chan tasks.Task, 100)
	doneCh := make(chan struct{}, 4)

	nReaders := 4
	err := tasks.LoadWeekly(os.Args[1], taskCh, doneCh)
	if err != nil {
		log.Fatal(err)
	}
	err = tasks.LoadMonthly(os.Args[2], taskCh, doneCh)
	if err != nil {
		log.Fatal(err)
	}
	err = tasks.LoadRecurring(os.Args[3], taskCh, doneCh)
	if err != nil {
		log.Fatal(err)
	}
	err = tasks.LoadRecurring(os.Args[4], taskCh, doneCh)
	if err != nil {
		log.Fatal(err)
	}

	grouper.Add(taskCh, doneCh, nReaders)

	tasksGroups := grouper.Filter(n)
	for day := 0; day <= n; day++ {
		tasks, ok := tasksGroups[day]
		if !ok {
			continue
		}
		fmt.Printf("Day = %d\n", day)
		for _, task := range tasks {
			fmt.Println(task)
		}
	}
}
