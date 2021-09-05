package cmd

import (
	"errors"
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

const (
	name    = "calendar-tasks"
	version = "0.0.1" //hard-code version for now

	printTimeFormat = "[Mon] Jan 2 2006"
)

type cmdArgs struct {
	days             int
	weeklySources    stringSliceArg
	monthlySources   stringSliceArg
	multiDateSources stringSliceArg
	version          bool
}

func (args *cmdArgs) parseArgs(argsIn []string) error {
	fs := flag.NewFlagSet("calendar-tasks", flag.ExitOnError)
	fs.IntVar(&args.days, "d", 0, "days ahead to get tasks")
	fs.Var(&args.weeklySources, "weekly", "weekly task source file path")
	fs.Var(&args.monthlySources, "monthly", "monthly task source file path")
	fs.Var(&args.multiDateSources, "multi", "multiDate task source file path")
	fs.BoolVar(&args.version, "v", false, "display version info")

	setUsage(fs)
	return fs.Parse(argsIn)
}

// Run excutes the CLI
func Run(argsIn []string) error {
	args := &cmdArgs{}
	err := setupArgs(args, argsIn)
	if err != nil {
		return err
	}

	if args.version {
		printVersion()
		return nil
	}

	date := fixDate(time.Now())

	taskChan := make(chan tasks.Task, 1000) // buffer large enough for reasonable amount of tasks
	done := make(chan struct{})

	loader := tasks.NewLoader(taskChan, done)
	processor := tasks.NewProcessor(date, args.days, taskChan, done)

	loader.AddWeeklySource(args.weeklySources...)
	loader.AddMonthlySource(args.monthlySources...)
	loader.AddMultiDateSource(args.multiDateSources...)

	err = start(loader, processor)
	if err != nil {
		return err
	}

	printTasks(processor, args.days, date)
	return nil
}

func start(loader *tasks.Loader, processor *tasks.Processor) error {
	// start the processor and wait on it to finish before returning
	processor.Start()
	defer processor.Wait()

	// start the loader and return any errors
	return loader.Start()
}

// fixDate returns a time.Time object matching the year, month, day (and location) of the argument
// and sets the hour to the middle of the day to avoid any boundary cases that can occur with
// e.g., daylight savings
func fixDate(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
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
		curDayStr := curDay.Format(printTimeFormat)
		if day == 0 {
			curDayStr += " (today)"
		}
		fmt.Println(curDayStr)

		for _, tsk := range tsks {
			fmt.Printf("\t-%s\n", tsk)
		}
	}
}

func setupArgs(args *cmdArgs, argsIn []string) error {
	if len(argsIn) < 2 {
		return args.parseArgs([]string{"--help"})
	}
	err := args.parseArgs(argsIn[1:])
	if err != nil {
		return err
	}
	if args.version {
		return nil
	}
	if len(args.weeklySources)+len(args.monthlySources)+len(args.multiDateSources) == 0 {
		return errors.New("no source files provided")
	}
	return nil
}

// stringSliceArg is a flag that accumulates strings
type stringSliceArg []string

func (s *stringSliceArg) String() string {
	return strings.Join(*s, " ")
}

func (s *stringSliceArg) Set(val string) error {
	*s = append(*s, val)
	return nil
}

func setUsage(fs *flag.FlagSet) {
	fs.Usage = func() {
		fmt.Printf("%s displays upcoming scheduled tasks\n\n", name)
		fmt.Print("Usage:\n")
		fmt.Printf("  %s [flags]\n\n", name)
		fmt.Printf("Flags:\n")
		fs.PrintDefaults()
	}
}

func printVersion() {
	fmt.Printf("%s: v%s\n", name, version)
}
