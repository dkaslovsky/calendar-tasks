package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// environment variables
	envWeeklySources  = "CALENDAR_TASKS_WEEKLY_SOURCES"
	envMonthlySources = "CALENDAR_TASKS_MONTHLY_SOURCES"
	envAnnualSources  = "CALENDAR_TASKS_ANNUAL_SOURCES"
	envSingleSources  = "CALENDAR_TASKS_SINGLE_SOURCES"

	// format for date flag input
	inputDateFormat = "2006-01-02"
)

type cliOpts struct {
	days         int
	back         int
	date         time.Time
	printVersion bool

	weeklySources  []string
	monthlySources []string
	annualSources  []string
	singleSources  []string
}

func parseArgs(argsIn []string, opts *cliOpts) error {
	var date string

	flag.StringVar(&date, "d", "", "starting date (YYY-MM-DD)")
	flag.StringVar(&date, "date", "", "starting date (YYY-MM-DD)")
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

	if date == "" {
		opts.date = time.Now()
	} else {
		var err error
		opts.date, err = time.Parse(inputDateFormat, date)
		if err != nil {
			return fmt.Errorf("invalid date: --date %s does not match YYYY-MM-DD format", date)
		}
	}

	// parse environment variables
	opts.weeklySources = parseStringSliceEnvVar(os.Getenv(envWeeklySources))
	opts.monthlySources = parseStringSliceEnvVar(os.Getenv(envMonthlySources))
	opts.annualSources = parseStringSliceEnvVar(os.Getenv(envAnnualSources))
	opts.singleSources = parseStringSliceEnvVar(os.Getenv(envSingleSources))
	if (len(opts.weeklySources) + len(opts.monthlySources) + len(opts.annualSources) + len(opts.singleSources)) == 0 {
		return errors.New("no source files provided, use --help for usage")
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

type appInfo struct {
	name    string
	version string
}

func (info *appInfo) setUsage() func() {
	return func() {
		fmt.Printf("%s displays upcoming scheduled tasks\n", info.name)
		fmt.Printf("\nTasks are read from files specified in comma-separated environment variables:\n")
		fmt.Printf("  %s\t\tsource files for weekly tasks\t\tex: %s=\"file1,file2,...\"\n", envWeeklySources, envWeeklySources)
		fmt.Printf("  %s\tsource files for monthly tasks\t\tex: %s=\"file1,file2,...\"\n", envMonthlySources, envMonthlySources)
		fmt.Printf("  %s\t\tsource files for annual tasks\t\tex: %s=\"file1,file2,...\"\n", envAnnualSources, envAnnualSources)
		fmt.Printf("  %s\t\tsource files for single tasks\t\tex: %s=\"file1,file2,...\"\n", envSingleSources, envSingleSources)
		fmt.Print("\nUsage:\n")
		fmt.Printf("  %s [flags] [args]\n", info.name)
		fmt.Printf("\nArgs:\n")
		fmt.Printf("  days int\t number of days from date to get tasks \t\tdefault: 0 (today)\n")
		fmt.Printf("\nFlags:\n")
		fmt.Printf("  -b, --back\t number of days back from date to get tasks \tdefault: 0 (none)\n")
		fmt.Printf("  -d, --date\t date in YYYY-MM-DD format \t\t\tdefault: today\n")
		fmt.Printf("  -h, --help\t display usage information\n")
		fmt.Printf("  -v, --version\t display version information\n")
	}
}

func (info *appInfo) printVersion() {
	fmt.Printf("%s: v%s\n", info.name, info.version)
}
