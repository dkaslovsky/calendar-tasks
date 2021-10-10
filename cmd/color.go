package cmd

import (
	"fmt"
	"runtime"
)

// colors for printing
const (
	reset  = "\033[0m"
	white  = "\033[97m"
	yellow = "\033[33m"
	purple = "\033[35m"
)

type color string

var (
	colorReset  color = reset
	colorToday  color = white
	colorPast   color = purple
	colorFuture color = yellow
)

// windows does not support color printing
func init() {
	if runtime.GOOS == "windows" {
		colorReset = ""
		colorToday = ""
		colorPast = ""
		colorFuture = ""
	}
}

func colorPrint(clr color, args ...interface{}) {
	fmt.Printf("%s%s%s", clr, fmt.Sprint(args...), colorReset)
}
