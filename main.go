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

	weekly, err := tasks.LoadWeekly(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	monthly, err := tasks.LoadMonthly(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	recurringMonthly, err := tasks.LoadRecurring(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	daily, err := tasks.LoadRecurring(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}

	grouper.Add(weekly, monthly, recurringMonthly, daily)

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
