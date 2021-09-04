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
	loader.AddWeeklySource(os.Args[1])
	loader.AddMonthlySource(os.Args[2])
	loader.AddRecurringSource(os.Args[3], os.Args[4])

	processor := tasks.NewProcessor(now, maxDays, taskChan, done)

	err := run(loader, processor)
	if err != nil {
		log.Fatal(err)
	}

	printTasks(processor.Tasks(), maxDays)
}

func run(loader *tasks.Loader, processor *tasks.Processor) error {
	// start the processor and wait on it to finish before returning
	processor.Start()
	defer processor.Wait()

	// start loading
	return loader.Start()
}

func printTasks(tsMap map[int][]tasks.Task, n int) {
	for day := 0; day <= n; day++ {
		ts, ok := tsMap[day]
		if !ok {
			continue
		}

		// sort for consistent ordering (current implementation is broken, will need proper String())
		sort.Slice(ts, func(i, j int) bool {
			return strings.ToLower(ts[i].String()) > strings.ToLower(ts[j].String())
		})

		fmt.Printf("Day = %d\n", day)
		for _, task := range ts {
			fmt.Println(task)
		}
	}
}
