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
`calendar-tasks` is a simple tool for tracking scheduled tasks from the command line.

To check tasks for the next `n` days, pass `n` as the argument to `calendar-tasks`.
For example:
```
$ calendar-tasks 3
[Sun] Sep 5 2021 (today)
    - pay bills
    - walk dog
[Tue] Sep 7 2021
    - garbage night
    - return rental car
    - workout
```
shows tasks for the next 3 days, starting with today.
Days without tasks are omitted from the output.

Tasks from previous days can be included in the output by specifying the number of days back from today to include with the `-b` or `--back` flag.

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
  CALENDAR_TASKS_SINGLE_SOURCES		source files for single tasks		ex: CALENDAR_TASKS_SINGLE_SOURCES="file1,file2,..."

Usage:
  calendar-tasks [flags] [args]

Args:
  days int	 number of days from today to get tasks 	default: 0 (today)

Flags:
  -b, --back	 number of days back from today to get tasks 	default: 0
  -h, --help	 display usage information
  -v, --version	 display version information
```

## Task Source Files
Tasks are stored in text files, the paths to which are set using environment variables.
There are four types of supported task files: weekly, monthly, annual, and single (see descriptions below).

- Paths to all weekly task files are stored in the `CALENDAR_TASKS_WEEKLY_SOURCES` environment variable.

- Paths to all monthly task files are stored in the `CALENDAR_TASKS_MONTHLY_SOURCES` environment variable.

- Paths to all annual task files are stored in the `CALENDAR_TASKS_ANNUAL_SOURCES` environment variable.

- Paths to all single task files are stored in the `CALENDAR_TASKS_SINGLE_SOURCES` environment variable.

Each environment variable supports specifying multiple files so that the source files can be organized however a user wishes.
For example, it might be convenient to store each month's tasks in separate monthly task files.
Specify multiple files with a comma-separated list.

</br>

### Weekly Task Source Files
Weekly tasks are tasks that occur on the same day each week. Such tasks are stored in a file with each line having the form `<day-of-the-week>: <task>`. For example
```
Sun: Grocery shopping
Sun: Play with kids
Wed: Garbage night
Thu: Coffee with Amy
Wed/Sat/Sun: Hiking
```
Note that each line contains only one task and that days can be repeated on multiple lines.
Tasks occurring on multiple days are indicated by using the forward-slash separator between days: `<day-of-the-week1/day-of-the-week2>/...: <task>`.
Days can be specified using their full name or three-letter abbreviation.

</br>

### Monthly Task Source Files
Monthly tasks are tasks that occur on the same day each month. Such tasks are stored in a file with each line having the form `<day-of-the-month>: <task>`. For example
```
3: Pay credit card bill
15: Meet with Alice and Bob
15: Poker night
25: Pay Mortgage
15/30: Pick up paycheck
```
Note that each line contains only one task and that days of the month can be repeated.
Tasks occurring on multiple days are indicated by using the forward-slash separator between days: `<day-of-the-month>/<day-of-the-month>/...: <task>`.

</br>

### Annual Task Source Files
Annual tasks are tasks that occur on a specific day of the year, specified by a month and a day.
Such tasks are stored in a file with each line having the form `<month day-of-the-month>: <task>`.
For example,
```
Jan 12: Daughter's birthday
April 15: File taxes
May 1: Renew lease
Mar 1/Nov 1: Change smoke alarm batteries
```
Note that each line contains only one task and that dates can be repeated.
Tasks occurring on multiple dates are indicated by using the forward-slash separator between dates: `<month day-of-the-month>/<month day-of-the-month>/...: <task>`.
Months can be specified using their full name or common abbreviation.

</br>

### Single Task Source Files
Single tasks are tasks that occur on a specific date, specified by a year, month, and day.
Single tasks, by definition, are not recurring.
Such tasks are stored in a file with each line having the form `<month day-of-the-month year>: <task>`.
For example,
```
Jan 12 2021: Pickup Alice from airport
Mar 20 2021: Drop off dog
Mar 20 2021: Trip to New York
```
Note that each line contains only one task and that dates can be repeated.
Tasks can occur on multiple dates, separated by the usual forward-slash (`/`) delimiter, however this concept makes less sense for single tasks than for those tasks that are recurring.
Months can be specified using their full name or common abbreviation.

Because single tasks are not recurring, it might be desirable to remove past single tasks from time to time.
`calendar-tasks` does not provide this functionality so as to keep its implementation minimal.
However, since tasks are stored in plaintext files, it is easy to prune tasks using standard command line tools.
A function for removing all single tasks prior to a specified date (and other useful functions) is available in a [gist](https://gist.github.com/dkaslovsky/d492bfb792133a46cb02c4a8c71372e3).

</br>

## Implementation Notes

### Why not use a structured file format?
While task files could have been structured as yaml, json, or some other standard format, `calendar-tasks` uses the above format for readability and ease of manipulation by other command line tools.
The current plain-text format is easy to read and modify, but structured files might also be supported in a future version.

### Why implement concurrent task file reads?
Task files are generally small and fast to read.
The performance gain by reading them with concurrent goroutines is negligible, if it even exists.
The reason for reading files concurrently is simply that it made `calendar-tasks` a more interesting project to develop.
