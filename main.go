package main

import (
	"log"
	"os"

	"github.com/dkaslovsky/calendar-tasks/cmd"
)

func main() {
	log.SetFlags(0)

	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
