package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

func main() {
	maxDays, _ := strconv.Atoi(os.Args[5])

	// now := time.Now()
	now := time.Date(2021, 8, 14, 18, 0, 0, 0, time.Local)

	loader := tasks.NewLoader()
	loader.AddWeekly(os.Args[1])
	loader.AddMonthly(os.Args[2])
	loader.AddRecurring(os.Args[3], os.Args[4])

	err := loader.Start()
	if err != nil {
		log.Fatal(err)
	}

	consumer := tasks.NewConsumer(now, maxDays, loader.Ch, loader.Close)
	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}

	for day, ts := range consumer.Get() {
		fmt.Printf("Day = %d\n", day)
		sort.Slice(ts, func(i, j int) bool {
			return strings.ToLower(ts[i].String()) > strings.ToLower(ts[j].String())
		})
		for _, task := range ts {
			fmt.Println(task)
		}
	}
}
