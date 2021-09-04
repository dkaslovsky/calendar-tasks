package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

func main() {

	args := cmdArgs{}
	args.attachArgs()

	taskChan := make(chan tasks.Task, 1000) // buffer size is large enough for a reasonable amount of tasks
	done := make(chan struct{})

	loader := tasks.NewLoader(taskChan, done)
	processor := tasks.NewProcessor(time.Now(), args.days, taskChan, done)

	loader.AddWeeklySource(args.weeklySources...)
	loader.AddMonthlySource(args.monthlySources...)
	loader.AddRecurringSource(args.recurringSources...)

	err := run(loader, processor)
	if err != nil {
		log.Fatal(err)
	}

	printTasks(processor.Tasks(), args.days)
}

func run(loader *tasks.Loader, processor *tasks.Processor) error {
	// start the processor and wait on it to finish before returning
	processor.Start()
	defer processor.Wait()

	// start the loader and return any errors
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

type cmdArgs struct {
	days             int
	weeklySources    stringSliceArg
	monthlySources   stringSliceArg
	recurringSources stringSliceArg
}

func (args *cmdArgs) attachArgs() {
	flag.IntVar(&args.days, "days", 0, "days ahead to get tasks")
	flag.Var(&args.weeklySources, "weekly", "weekly task source file")
	flag.Var(&args.monthlySources, "monthly", "monthly task source file")
	flag.Var(&args.recurringSources, "recurring", "recurring task source file")
	flag.Parse()
}

type stringSliceArg []string

func (s *stringSliceArg) String() string {
	return strings.Join(*s, " ")
}

func (s *stringSliceArg) Set(val string) error {
	*s = append(*s, val)
	return nil
}
