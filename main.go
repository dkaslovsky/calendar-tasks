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

	//now := time.Now()
	now := time.Date(2021, 8, 14, 0, 0, 0, 0, time.Local)
	fmt.Printf("\nDEBUG MODE - USING FIXED DATE %s\n", now)

	args := cmdArgs{}
	args.attachArgs()

	taskChan := make(chan tasks.Task, 1000) // buffer size is large enough for a reasonable amount of tasks
	done := make(chan struct{})

	loader := tasks.NewLoader(taskChan, done)
	processor := tasks.NewProcessor(now, args.days, taskChan, done)

	loader.AddWeeklySource(args.weeklySources...)
	loader.AddMonthlySource(args.monthlySources...)
	loader.AddMultiDateSource(args.multiDateSources...)

	err := run(loader, processor)
	if err != nil {
		log.Fatal(err)
	}

	printTasks(processor, args.days, now)
}

func run(loader *tasks.Loader, processor *tasks.Processor) error {
	// start the processor and wait on it to finish before returning
	processor.Start()
	defer processor.Wait()

	// start the loader and return any errors
	return loader.Start()
}

func printTasks(processor *tasks.Processor, numDays int, now time.Time) {
	for day := 0; day <= numDays; day++ {
		tsks, ok := processor.GetTasks(day)
		if !ok {
			continue
		}

		// sort for consistent ordering
		sort.Slice(tsks, func(i, j int) bool {
			return strings.ToLower(tsks[i].String()) < strings.ToLower(tsks[j].String())
		})

		dayStr := now.AddDate(0, 0, day).Format("[Mon] Jan 2 2006")
		fmt.Printf("\n%s\n", dayStr)
		for _, tsk := range tsks {
			fmt.Printf("\t-%s\n", tsk)
		}
	}
}

type cmdArgs struct {
	days             int
	weeklySources    stringSliceArg
	monthlySources   stringSliceArg
	multiDateSources stringSliceArg
}

func (args *cmdArgs) attachArgs() {
	flag.IntVar(&args.days, "days", 0, "days ahead to get tasks")
	flag.Var(&args.weeklySources, "weekly", "weekly task source file")
	flag.Var(&args.monthlySources, "monthly", "monthly task source file")
	flag.Var(&args.multiDateSources, "multi", "multiDate task source file")
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
