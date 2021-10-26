package cmd

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

const (
	// environment variables
	envWeeklySources  = "CALENDAR_TASKS_WEEKLY_SOURCES"
	envMonthlySources = "CALENDAR_TASKS_MONTHLY_SOURCES"
	envAnnualSources  = "CALENDAR_TASKS_ANNUAL_SOURCES"
	envSingleSources  = "CALENDAR_TASKS_SINGLE_SOURCES"

	// format for displaying dates
	printTimeFormat = "[Mon] Jan 2 2006"
)

type cliOpts struct {
	name    string
	version string

	days int
	back int

	printVersion bool

	weeklySources  []string
	monthlySources []string
	annualSources  []string
	singleSources  []string
}

// Run excutes the CLI
func Run(name string, version string, argsIn []string) error {
	opts := &cliOpts{
		name:    name,
		version: version,
	}

	flag.Usage = setUsage(opts)

	err := parseArgs(argsIn, opts)
	if err != nil {
		return err
	}

	if opts.printVersion {
		printVersion(opts)
		return nil
	}

	return run(opts)
}

func run(opts *cliOpts) error {
	runDates := getRunDates(time.Now(), opts)

	taskChan := make(chan tasks.Task, 1000) // buffer large enough for reasonable amount of tasks
	doneChan := make(chan struct{})

	loader := tasks.NewLoader(taskChan, doneChan)
	processor := tasks.NewProcessor(runDates.start, runDates.numDays, taskChan, doneChan)

	loader.AddWeeklySource(opts.weeklySources...)
	loader.AddMonthlySource(opts.monthlySources...)
	loader.AddAnnualSource(opts.annualSources...)

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

func parseArgs(argsIn []string, opts *cliOpts) error {
	flag.IntVar(&opts.back, "b", 0, "number of days back from today")
	flag.IntVar(&opts.back, "back", 0, "number of days back from today")
	flag.BoolVar(&opts.printVersion, "v", false, "display version information")
	flag.BoolVar(&opts.printVersion, "version", false, "display version information")
	flag.Parse()

	if opts.printVersion {
		return nil
	}

	if opts.back < 0 {
		return fmt.Errorf("invalid negative value: --back %d", opts.back)
	}

	// parse environment variables
	opts.weeklySources = parseStringSliceEnvVar(os.Getenv(envWeeklySources))
	opts.monthlySources = parseStringSliceEnvVar(os.Getenv(envMonthlySources))
	opts.annualSources = parseStringSliceEnvVar(os.Getenv(envAnnualSources))
	opts.singleSources = parseStringSliceEnvVar(os.Getenv(envSingleSources))
	if (len(opts.weeklySources) + len(opts.monthlySources) + len(opts.annualSources) + len(opts.singleSources)) == 0 {
		return fmt.Errorf("no source files provided, run `%s --help` for usage", opts.name)
	}

	// run with defaults
	if flag.NArg() == 0 {
		opts.days = 0
		return nil
	}

	// parse day arg
	dayStr := flag.Arg(0)
	days, err := strconv.Atoi(dayStr)
	if err != nil {
		return fmt.Errorf("unparsable integer argument: %s", dayStr)
	}
	if days < 0 {
		return fmt.Errorf("invalid negative argument: %d", days)
	}
	opts.days = days

	return nil
}

// parseStringSliceEnvVar parses a comma-separated environment variable into a slice of string
func parseStringSliceEnvVar(envStr string) []string {
	parsed := []string{}
	if envStr == "" {
		return parsed
	}
	split := strings.Split(envStr, ",")
	for _, s := range split {
		parsed = append(parsed, strings.TrimSpace(s))
	}
	return parsed
}

func setUsage(opts *cliOpts) func() {
	return func() {
		fmt.Printf("%s displays upcoming scheduled tasks\n", opts.name)
		fmt.Printf("\nTasks are read from files specified in comma-separated environment variables:\n")
		fmt.Printf("  %s\t\tsource files for weekly tasks\t\tex: %s=\"file1,file2,...\"\n", envWeeklySources, envWeeklySources)
		fmt.Printf("  %s\tsource files for monthly tasks\t\tex: %s=\"file1,file2,...\"\n", envMonthlySources, envMonthlySources)
		fmt.Printf("  %s\t\tsource files for annual tasks\t\tex: %s=\"file1,file2,...\"\n", envAnnualSources, envAnnualSources)
		fmt.Printf("  %s\t\tsource files for single tasks\t\tex: %s=\"file1,file2,...\"\n", envSingleSources, envSingleSources)
		fmt.Print("\nUsage:\n")
		fmt.Printf("  %s [flags] [args]\n", opts.name)
		fmt.Printf("\nArgs:\n")
		fmt.Printf("  days int\t number of days from today to get tasks \tdefault: 0 (today)\n")
		fmt.Printf("\nFlags:\n")
		fmt.Printf("  -b, --back\t number of days back from today to get tasks \tdefault: 0\n")
		fmt.Printf("  -h, --help\t display usage information\n")
		fmt.Printf("  -v, --version\t display version information\n")
	}
}

func printVersion(opts *cliOpts) {
	fmt.Printf("%s: v%s\n", opts.name, opts.version)
}

type runDates struct {
	today   time.Time
	start   time.Time
	numDays int
}

func getRunDates(now time.Time, opts *cliOpts) *runDates {
	today := fixDate(now)
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
