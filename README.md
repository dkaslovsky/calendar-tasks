# Calendar-Tasks
Simple CLI for tracking scheduled tasks

## Usage

`calendar-tasks` is a simple tool for tracking scheduled tasks from the commandline rather than on a calendar.

To check tasks for the next 3 days, pass the task files to `calendar-tasks` (see below) together with `-d 3` for the number of days. For example:
```
$ calendar-tasks --weekly weekly-tasks.txt -d 3
[Sun] Sep 5 2021 (today)
    - walk dog
    - pay bills
[Tue] Sep 7 2021
    - workout
    - garbage night
    - return rental car
```

`calendar-tasks` handles leap years as well as differences in the lengths of months. For example, a task scheduled for the 30th of every month will not be skipped in February.
Instead, it will be shown on March 1 for leap years and March 2 for non-leap years.

A summary of the arguments that can be passed to `calendar-tasks` is displayed in its help menu:
```
$ calendar-tasks --help
calendar-tasks displays upcoming scheduled tasks

Usage:
  calendar-tasks [flags]

Flags:
  -d int
    	days ahead to get tasks
  -monthly value
    	monthly task source file path
  -multi value
    	multiDate task source file path
  -v	display version info
  -weekly value
    	weekly task source file path
```


## Task Files

Tasks are stored in text files which are pass to `calendar-tasks` as arguments.
There are three types of supported task files: annual, montly, and multi-dated

### Weekly
Weekly tasks are tasks that occur on the same day each week. Such tasks are stored in a file with each line having the form `<day of the week>:<task>`. For example
```
Sun: Grocery shopping
Sun: Play with kids
Wed: Garbage night
Thu: Coffee with Amy
```
Note that each line contains only one task and that days can be repeated. Pass these files to `calendar-tasks` using the `--weekly` flag.

### Monthly
Monthly tasks are tasks that occur on the same day each month. Such tasks are stored in a file with each line having the form `<day>:<task>`. For example
```
3: Pay credit card bill
15: Meet with Alice and Bob
15: Poker night
30: Pay Mortgage
```
Note that each line contains only one task and that days of the month can be repeated.
Pass these files to `calendar-tasks` using the `--monthly` flag.

### Multi-dated
Multi-dated tasks are tasks that occur on multiple dates. There are usually two types: multiple-month tasks and annual tasks.

*Multiple-month*

Multiple-month tasks are tasks that occur on the same day of multiple months. Such tasks are stored in a file with each line having the form `<month/month/.../month day>:<task>`. For example
```
Mar/Nov 1: Change smoke alarm batteries
Jan/Jun 15: Pay property taxes
Jan/Apr/Jul/Oct 30: Consulting appointment
```
Note that each line contains only one task and that dates can be repeated.
Pass these files to `calendar-tasks` using the `--multi` flag.

*Annual*

Annual tasks are tasks that occur on the same day each year. Such tasks are stored in a file with each line having the form `<month day>:<task>`. For example
```
Jan 12: Daughter's birthday
April 15: File taxes
May 1: Renew lease
```
Note that each line contains only one task and that dates can be repeated.
Pass these files to `calendar-tasks` using the `--multi` flag.

## Why the unstructured text files?
The task files certainly could have been structured as yaml, json, or some other structured format, but `calendar-tasks` uses the above format to be human-readable and compatible with manipulation by other commandline tools. I find the current plain-text format easy to read and modify, but might also support structured files at some point in the future.
