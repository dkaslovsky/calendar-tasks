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
	date := fixDate(time.Now())

	args := cmdArgs{}
	args.attachArgs()

	taskChan := make(chan tasks.Task, 1000) // buffer size is large enough for a reasonable amount of tasks
	done := make(chan struct{})

	loader := tasks.NewLoader(taskChan, done)
	processor := tasks.NewProcessor(date, args.days, taskChan, done)

	loader.AddWeeklySource(args.weeklySources...)
	loader.AddMonthlySource(args.monthlySources...)
	loader.AddMultiDateSource(args.multiDateSources...)

	err := run(loader, processor)
	if err != nil {
		log.Fatal(err)
	}

	printTasks(processor, args.days, date)
}

func run(loader *tasks.Loader, processor *tasks.Processor) error {
	// start the processor and wait on it to finish before returning
	processor.Start()
	defer processor.Wait()

	// start the loader and return any errors
	return loader.Start()
}

func printTasks(processor *tasks.Processor, numDays int, startDate time.Time) {
	for day := 0; day <= numDays; day++ {
		tsks, ok := processor.GetTasks(day)
		if !ok {
			continue
		}

		// sort for consistent ordering
		sort.Slice(tsks, func(i, j int) bool {
			return strings.ToLower(tsks[i].String()) < strings.ToLower(tsks[j].String())
		})

		curDay := startDate.AddDate(0, 0, day)
		fmt.Printf("\n%s\n", curDay.Format("[Mon] Jan 2 2006"))
		for _, tsk := range tsks {
			fmt.Printf("\t-%s\n", tsk)
		}
	}
}

// fixDate returns a time.Time object matching the year, month, day (and location) of the argument
// and sets the hour to the middle of the day to avoid any boundary cases that can occur with
// e.g., daylight savings
func fixDate(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
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
