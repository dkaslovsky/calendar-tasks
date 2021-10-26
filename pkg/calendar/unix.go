package calendar

import "time"

// UnixToDaysFloored converts a unix timestamp to days (rounding down)
func UnixToDaysFloored(t int64) int64 {
	return t / int64(24*time.Hour.Seconds())
}
