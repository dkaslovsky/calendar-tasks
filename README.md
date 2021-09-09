# calendar-tasks
Simple CLI for tracking scheduled tasks

</br>

## Installation
Using Go >= 1.17, run
```
$ go install github.com/dkaslovsky/calendar-tasks@latest
```
Earlier versions of Go might use the `go get` command.

</br>

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
  CALENDAR_TASKS_WEEKLY_SOURCES		source files for weekly tasks		ex: CALENDAR_TASKS_WEEKLY_SOURCES="file1,file2,..."
  CALENDAR_TASKS_MONTHLY_SOURCES	source files for monthly tasks		ex: CALENDAR_TASKS_MONTHLY_SOURCES="file1,file2,..."
  CALENDAR_TASKS_ANNUAL_SOURCES		source files for annual tasks		ex: CALENDAR_TASKS_ANNUAL_SOURCES="file1,file2,..."

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
There are three types of supported task files: weekly, monthly, and annual (see descriptions below).

Paths to all weekly task files are stored in the `CALENDAR_TASKS_WEEKLY_SOURCES` environment variable.

Paths to all monthly task files are stored in the `CALENDAR_TASKS_MONTHLY_SOURCES` environment variable.

Paths to all annual task files are stored in the `CALENDAR_TASKS_ANNUAL_SOURCES` environment variable.

Each environment variable supports specifying multiple files so that the source files can be organized however a user wishes.
For example, it might be convenient to store each month's tasks in separate monthly task files.

</br>

### Weekly Task Source Files
Weekly tasks are tasks that occur on the same day each week. Such tasks are stored in a file with each line having the form `<day-of-the-week>:<task>`. For example
```
Sun: Grocery shopping
Sun: Play with kids
Wed: Garbage night
Thu: Coffee with Amy
Wed/Sat/Sun: Hiking
```
Note that each line contains only one task and that days can be repeated on multiple lines.
Tasks occurring on multiple days are indicated by using the forward-slash separator between days: `<day-of-the-week1/day-of-the-week2>/...:<task>`.
Days can be specified using their full name or three-letter abbreviation.

</br>

### Monthly Task Source Files
Monthly tasks are tasks that occur on the same day each month. Such tasks are stored in a file with each line having the form `<day-of-the-month>:<task>`. For example
```
3: Pay credit card bill
15: Meet with Alice and Bob
15: Poker night
25: Pay Mortgage
15/30: Pick up paycheck
```
Note that each line contains only one task and that days of the month can be repeated.
Tasks occurring on multiple days are indicated by using the forward-slash separator between days: `<day-of-the-month>/<day-of-the-month>/...:<task>`.

</br>

### Annual Task Source Files
Annual tasks are tasks that occur on a specific day of the year, specified by a month and a day.
Such tasks are stored in a file with each line having the form `<month day-of-the-month>:<task>`.
For example,
```
Jan 12: Daughter's birthday
April 15: File taxes
May 1: Renew lease
Mar 1/Nov 1: Change smoke alarm batteries
```
Note that each line contains only one task and that dates can be repeated.
Tasks occurring on multiple dates are indicated by using the forward-slash separator between dates: `<month day-of-the-month>/<month day-of-the-month>/...:<task>`.
Months can be specified using their full name or common abbreviation.

</br>

## Implementation Notes

### Why not use a structured file format?
While task files could have been structured as yaml, json, or some other standard format, `calendar-tasks` uses the above format for readability and ease of manipulation by other commandline tools.
The current plain-text format is easy to read and modify, but structured files might also be supported in a future version.

### Why not read the task files in serial?
Task files are generally small and fast to read.
The performance gain by reading them in parallel is negligible, if it even exists.
The reason for reading files in parallel is simply that it made this a more interesting project to develop.
