package main

import (
	"log"
	"os"

	"github.com/dkaslovsky/calendar-tasks/cmd"
)

const (
	name    = "calendar-tasks" // app name
	version = "0.3.0"          // hard-code version for now
)

func main() {
	log.SetFlags(0)

	err := cmd.Run(name, version, os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
