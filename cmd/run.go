package cmd

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

// format for displaying dates
const printTimeFormat = "[Mon] Jan 2 2006"

// Run excutes the CLI
func Run(name string, version string, argsIn []string) error {
	info := &appInfo{
		name:    name,
		version: version,
	}

	flag.Usage = info.setUsage()

	opts := &cliOpts{}
	err := parseArgs(argsIn, opts)
	if err != nil {
		return err
	}

	if opts.printVersion {
		info.printVersion()
		return nil
	}

	return run(opts)
}

func run(opts *cliOpts) error {
	runDates := getRunDates(opts)

	taskChan := make(chan tasks.Task, 1000) // buffer large enough for reasonable amount of tasks
	doneChan := make(chan struct{})

	loader := tasks.NewLoader(taskChan, doneChan)
	processor := tasks.NewProcessor(runDates.start, runDates.numDays, taskChan, doneChan)

	loader.AddWeeklySource(opts.weeklySources...)
	loader.AddMonthlySource(opts.monthlySources...)
	loader.AddAnnualSource(opts.annualSources...)
	loader.AddSingleSource(opts.singleSources...)

	err := processTasks(loader, processor)
	if err != nil {
		return err
	}

	printTasks(processor, runDates)
	return nil
}

// processTasks starts the processor and loader and waits on the processor before returning
func processTasks(loader *tasks.Loader, processor *tasks.Processor) error {
	// start the processor and wait on it to finish before returning
	processor.Start()
	defer processor.Wait()
	// start the loader and return any errors
	return loader.Start()
}

func printTasks(processor *tasks.Processor, dates *runDates) {
	numTasks := 0

	for day := 0; day <= dates.numDays; day++ {
		tsks, ok := processor.GetTasks(day)
		if !ok {
			continue
		}

		// sort for consistent ordering
		sort.Slice(tsks, func(i, j int) bool {
			return strings.ToLower(tsks[i].String()) < strings.ToLower(tsks[j].String())
		})

		// format printing
		var clr color
		var curDayStr string
		switch curDay := dates.start.AddDate(0, 0, day); {
		case curDay == dates.today:
			curDayStr = curDay.Format(printTimeFormat) + " (today)"
			clr = colorToday
		case curDay.After(dates.today):
			curDayStr = curDay.Format(printTimeFormat)
			clr = colorFuture
		default:
			// past
			curDayStr = curDay.Format(printTimeFormat)
			clr = colorPast
		}

		colorPrint(clr, curDayStr, "\n")
		for _, tsk := range tsks {
			colorPrint(clr, "\t-", tsk, "\n")
			numTasks++
		}
	}

	if numTasks == 0 {
		fmt.Println("no tasks")
	}
}

type runDates struct {
	today   time.Time
	start   time.Time
	numDays int
}

func getRunDates(opts *cliOpts) *runDates {
	today := fixDate(opts.date)
	start := today.AddDate(0, 0, -opts.back)
	numDays := opts.days + opts.back
	return &runDates{
		today:   today,
		start:   start,
		numDays: numDays,
	}
}

// fixDate returns a time.Time object matching the year, month, day (and location) of the argument
// and sets the hour to the middle of the day to avoid any boundary cases that can occur with
// e.g., daylight savings
func fixDate(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
}
