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

	daily, err := tasks.LoadDaily(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	monthly, err := tasks.LoadMonthly(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range daily {
		f.Add(d)
	}
	for _, m := range monthly {
		f.Add(m)
	}

	n, _ := strconv.Atoi(os.Args[3])

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
