package cmd

import (
	"fmt"
	"runtime"
)

// colors for printing
const (
	reset  = "\033[0m"
	gray   = "\033[37m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

type color string

var (
	colorReset  color = reset
	colorToday  color = yellow
	colorPast   color = blue
	colorFuture color = gray
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
