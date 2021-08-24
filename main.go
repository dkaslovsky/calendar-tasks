package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

func main() {

	// now := time.Now()
	now := time.Date(2021, 8, 14, 18, 0, 0, 0, time.Local)
	maxDays, _ := strconv.Atoi(os.Args[5])

	taskChan := make(chan tasks.Task, 100) // TODO: make buffer size configurable
	done := make(chan struct{})

	loader := tasks.NewLoader(taskChan, done)
	loader.AddWeekly(os.Args[1])
	loader.AddMonthly(os.Args[2])
	loader.AddRecurring(os.Args[3], os.Args[4])

	consumer := tasks.NewConsumer(now, maxDays, taskChan, done)

	// processing:
	// start the consumer, start loading from files, wait for the consumer to finish consuming
	consumer.Start()
	err := loader.Start()
	if err != nil {
		log.Fatal(err)
	}
	consumer.Wait()

	printTasks(consumer.Tasks(), maxDays)
}

func printTasks(tsMap map[int][]tasks.Task, n int) {
	for day := 0; day <= n; day++ {
		fmt.Printf("Day = %d\n", day)
		ts, ok := tsMap[day]
		if !ok {
			fmt.Println("No tasks")
		}
		sort.Slice(ts, func(i, j int) bool {
			return strings.ToLower(ts[i].String()) > strings.ToLower(ts[j].String())
		})
		for _, task := range ts {
			fmt.Println(task)
		}
	}
}
