package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

func main() {
	maxDays, _ := strconv.Atoi(os.Args[5])

	// now := time.Now()
	now := time.Date(2021, 8, 14, 18, 0, 0, 0, time.Local)

	loader := tasks.NewLoader()
	consumer := tasks.NewConsumer(now, maxDays, loader.Ch, loader.Wait)

	err := loader.Load(
		[]string{os.Args[1]},
		[]string{os.Args[2]},
		[]string{os.Args[3], os.Args[4]},
	)
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}

	for day, ts := range consumer.Get() {
		fmt.Printf("Day = %d\n", day)
		sort.Slice(ts, func(i, j int) bool {
			return ts[i].String() > ts[j].String()
		})
		for _, task := range ts {
			fmt.Println(task)
		}
	}
}
