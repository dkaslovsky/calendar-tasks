# calendar-tasks
Simple CLI for tracking scheduled tasks

## Installation
Using Go >= 1.17, run
```
$ go install github.com/dkaslovsky/calendar-tasks@latest
```

## Usage
`calendar-tasks` is a simple tool for tracking scheduled tasks from the commandline rather than on a calendar.

To check tasks for the next `n` days, pass `n` as the argument to `calendar-tasks`.
For example:
```
$ calendar-tasks 3
[Sun] Sep 5 2021 (today)
    - walk dog
    - pay bills
[Tue] Sep 7 2021
    - workout
    - garbage night
    - return rental car
```
shows tasks for the next 3 days, starting with today.
Days without tasks are omitted from the output.

`calendar-tasks` properly handles leap years and months with fewer than 31 days.
For example, a task scheduled for the 30th of every month will not be skipped in February.
Instead, it will be shown on March 1 for leap years and March 2 for non-leap years.

A usage summary is displayed from the help menu:
```
$ calendar-tasks --help
calendar-tasks displays upcoming scheduled tasks

Tasks are read from files specified in comma-separated environment variables:
  CALENDAR_TASKS_WEEKLY_SOURCES 	source files for weekly tasks   	ex: CALENDAR_TASKS_WEEKLY_SOURCES="file1,file2,..."
  CALENDAR_TASKS_MONTHLY_SOURCES	source files for monthly tasks  	ex: CALENDAR_TASKS_MONTHLY_SOURCES="file1,file2,..."
  CALENDAR_TASKS_MULTIDATE_SOURCES	source files for multi-date tasks	ex: CALENDAR_TASKS_MULTIDATE_SOURCES="file1,file2,..."

Usage:
  calendar-tasks [args]
  calendar-tasks [flags]

Args:
  days int	 number of days from today to get tasks 	default: 0 (today)

Flags:
  -h, --help	 display usage information
  -v, --version	 display version information
```

## Task Source Files
Tasks are stored in text files, the paths to which are set using environment variables.
There are three types of supported task files: weekly, monthly, and multi-dated (see subsections below).

Paths to all weekly files are stored in the `CALENDAR_TASKS_WEEKLY_SOURCES` environment variable.

Paths to all monthly files are stored in the `CALENDAR_TASKS_MONTHLY_SOURCES` environment variable.

Paths to all multi-date files are stored in the `CALENDAR_TASKS_MULTIDATE_SOURCES` environment variable.

</br>

### Weekly Task Source Files
Weekly tasks are tasks that occur on the same day each week. Such tasks are stored in a file with each line having the form `<day-of-the-week>:<task>`. For example
```
Sun: Grocery shopping
Sun: Play with kids
Wed: Garbage night
Thu: Coffee with Amy
```
Note that each line contains only one task and that days can be repeated.

</br>

### Monthly Task Source Files
Monthly tasks are tasks that occur on the same day each month. Such tasks are stored in a file with each line having the form `<day-of-the-month>:<task>`. For example
```
3: Pay credit card bill
15: Meet with Alice and Bob
15: Poker night
30: Pay Mortgage
```
Note that each line contains only one task and that days of the month can be repeated.

</br>

### Multi-dated Task Source Files
Multi-dated tasks are tasks that occur on multiple dates. There are usually two types: multiple-month tasks and annual tasks.

*Multiple-month*

Multiple-month tasks are tasks that occur on the same day of multiple months. Such tasks are stored in a file with each line having the form `<month/month/.../month day-of-the-month>:<task>`. For example
```
Mar/Nov 1: Change smoke alarm batteries
Jan/Jun 15: Pay property taxes
Jan/Apr/Jul/Oct 30: Consulting appointment
```
Note that each line contains only one task and that dates can be repeated.

*Annual*

Annual tasks are tasks that occur on the same day each year. Such tasks are stored in a file with each line having the form `<month day-of-the-month>:<task>`. For example
```
Jan 12: Daughter's birthday
April 15: File taxes
May 1: Renew lease
```
Note that each line contains only one task and that dates can be repeated.

## Implementation Notes

### Why not use a structured file format?
While task files could have been structured as yaml, json, or some other standard format, but `calendar-tasks` uses the above format for ease of human-readable and manipulation by other commandline tools.
I find the current plain-text format easy to read and modify, but might also support structured files at some point in the future.

### Why not read the task files in serial?
Task files are generally small and therefore very fast to read.
The performance gain by reading them in parallel is negligible, if it even exists.
The reason for reading them in parallel is simply that it made this a more interesting project to develop.
