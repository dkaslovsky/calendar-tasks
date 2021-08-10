package tasks

import "time"

// Task represents a task to occur on a specified date(s)
type Task interface {
	DaysFrom(time.Time) int
	String() string
}
