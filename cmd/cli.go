package cmd

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

const (
	name    = "calendar-tasks" // app name
	version = "0.1.1"          // hard-code version for now

	// environment variables
	envWeeklySources  = "CALENDAR_TASKS_WEEKLY_SOURCES"
	envMonthlySources = "CALENDAR_TASKS_MONTHLY_SOURCES"
	envAnnualSources  = "CALENDAR_TASKS_ANNUAL_SOURCES"

	printTimeFormat = "[Mon] Jan 2 2006" // format for displaying dates
)

// colors for printing
var (
	reset  = "\033[0m"
	white  = "\033[97m"
	yellow = "\033[33m"
	purple = "\033[35m"

	colorToday  = white
	colorPast   = purple
	colorFuture = yellow
)

// windows does not support color printing
func init() {
	if runtime.GOOS == "windows" {
		reset = ""
		colorToday = ""
		colorPast = ""
		colorFuture = ""
	}
}

type cmdArgs struct {
	days int
	back int

	version bool

	weeklySources  []string
	monthlySources []string
	annualSources  []string
}

// Run excutes the CLI
func Run(argsIn []string) error {
	flag.Usage = setUsage()

	args := &cmdArgs{}
	err := args.parseArgs(argsIn)
	if err != nil {
		return err
	}

	if args.version {
		printVersion()
		return nil
	}

	if args.numSources() == 0 {
		return fmt.Errorf("no source files provided, run `%s --help` for usage", name)
	}

	runDates := getRunDates(time.Now(), args)

	taskChan := make(chan tasks.Task, 1000) // buffer large enough for reasonable amount of tasks
	doneChan := make(chan struct{})

	loader := tasks.NewLoader(taskChan, doneChan)
	processor := tasks.NewProcessor(runDates.start, runDates.numDays, taskChan, doneChan)

	loader.AddWeeklySource(args.weeklySources...)
	loader.AddMonthlySource(args.monthlySources...)
	loader.AddAnnualSource(args.annualSources...)

	err = run(loader, processor)
	if err != nil {
		return err
	}

	printTasks(processor, runDates)
	return nil
}

func run(loader *tasks.Loader, processor *tasks.Processor) error {
	// start the processor and wait on it to finish before returning
	processor.Start()
	defer processor.Wait()
	// start the loader and return any errors
	return loader.Start()
}

type runDates struct {
	today   time.Time
	start   time.Time
	numDays int
}

func getRunDates(now time.Time, args *cmdArgs) *runDates {
	today := fixDate(now)
	start := today.AddDate(0, 0, -args.back)
	numDays := args.days + args.back
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
		var color string
		var curDayStr string
		switch curDay := dates.start.AddDate(0, 0, day); {
		case curDay == dates.today:
			curDayStr = curDay.Format(printTimeFormat) + " (today)"
			color = colorToday
		case curDay.After(dates.today):
			curDayStr = curDay.Format(printTimeFormat)
			color = colorFuture
		default:
			// past
			curDayStr = curDay.Format(printTimeFormat)
			color = colorPast
		}

		colorPrint(color, curDayStr, "\n")
		for _, tsk := range tsks {
			colorPrint(color, "\t-", tsk, "\n")
			numTasks++
		}
	}

	if numTasks == 0 {
		fmt.Println("no tasks")
	}
}

func colorPrint(color string, args ...interface{}) {
	fmt.Printf("%s%s%s", color, fmt.Sprint(args...), reset)
}

func setUsage() func() {
	return func() {
		fmt.Printf("%s displays upcoming scheduled tasks\n", name)
		fmt.Printf("\nTasks are read from files specified in comma-separated environment variables:\n")
		fmt.Printf("  %s\t\tsource files for weekly tasks\t\tex: %s=\"file1,file2,...\"\n", envWeeklySources, envWeeklySources)
		fmt.Printf("  %s\tsource files for monthly tasks\t\tex: %s=\"file1,file2,...\"\n", envMonthlySources, envMonthlySources)
		fmt.Printf("  %s\t\tsource files for annual tasks\t\tex: %s=\"file1,file2,...\"\n", envAnnualSources, envAnnualSources)
		fmt.Print("\nUsage:\n")
		fmt.Printf("  %s [flags] [args]\n", name)
		fmt.Printf("\nArgs:\n")
		fmt.Printf("  days int\t number of days from today to get tasks \tdefault: 0 (today)\n")
		fmt.Printf("\nFlags:\n")
		fmt.Printf("  -b, --back\t number of days back from today to get tasks \tdefault: 0\n")
		fmt.Printf("  -h, --help\t display usage information\n")
		fmt.Printf("  -v, --version\t display version information\n")
	}
}

func printVersion() {
	fmt.Printf("%s: v%s\n", name, version)
}

func (args *cmdArgs) parseArgs(argsIn []string) error {
	flag.IntVar(&args.back, "b", 0, "number of days back from today")
	flag.IntVar(&args.back, "back", 0, "number of days back from today")
	flag.BoolVar(&args.version, "v", false, "display version information")
	flag.BoolVar(&args.version, "version", false, "display version information")
	flag.Parse()

	if args.version {
		return nil
	}

	// parse environment variables
	args.weeklySources = parseStringSliceEnvVar(os.Getenv(envWeeklySources))
	args.monthlySources = parseStringSliceEnvVar(os.Getenv(envMonthlySources))
	args.annualSources = parseStringSliceEnvVar(os.Getenv(envAnnualSources))

	// run with defaults
	if flag.NArg() == 0 {
		args.days = 0
		return nil
	}

	// parse args
	dayStr := flag.Arg(0)
	days, err := strconv.Atoi(dayStr)
	if err != nil {
		return fmt.Errorf("unparsable integer argument: %s", dayStr)
	}
	args.days = days

	// ignore invalid values
	if args.days < 0 {
		args.days = 0
	}
	if args.back < 0 {
		args.back = 0
	}

	return nil
}

func (args *cmdArgs) numSources() int {
	return len(args.weeklySources) + len(args.monthlySources) + len(args.annualSources)
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
