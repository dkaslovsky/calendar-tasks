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

	weekly, err := tasks.LoadWeekly(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range weekly {
		grouper.Add(t)
	}

	monthly, err := tasks.LoadMonthly(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range monthly {
		grouper.Add(t)
	}

	recurringMonthly, err := tasks.LoadRecurring(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range recurringMonthly {
		grouper.Add(t)
	}

	daily, err := tasks.LoadRecurring(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range daily {
		grouper.Add(t)
	}

	n, _ := strconv.Atoi(os.Args[5])

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
