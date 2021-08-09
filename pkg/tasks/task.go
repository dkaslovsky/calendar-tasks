package tasks

import "time"

type Task interface {
	DaysFrom(time.Time) int
	String() string
}
