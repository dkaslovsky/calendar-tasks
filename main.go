package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

func main() {
	now := time.Now()
	rs := days()
	//rs := months()

	upcoming := make(map[int][]taskItem)
	for _, r := range rs {
		appendUpcoming(upcoming, r, now)
	}

	for until, rr := range upcoming {
		fmt.Println(until)
		for _, r := range rr {
			fmt.Println(r)
		}
		fmt.Print("--------\n")
	}
}

type taskItem interface {
	DaysFrom(time.Time) int
	String() string
}

func appendUpcoming(upcoming map[int][]taskItem, r taskItem, now time.Time) {
	until := r.DaysFrom(now)
	if _, ok := upcoming[until]; !ok {
		upcoming[until] = []taskItem{}
	}
	upcoming[until] = append(upcoming[until], r)
}

func days() []*tasks.Daily {
	fileName := os.Args[1]
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(b), "\n")

	ds := []*tasks.Daily{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		d, err := tasks.NewDaily(line)
		if err != nil {
			log.Fatal(err)
		}
		ds = append(ds, d)
	}

	return ds
}

func months() []*tasks.Monthly {
	fileName := os.Args[1]
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(b), "\n")

	ms := []*tasks.Monthly{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		m, err := tasks.NewMonthly(line)
		if err != nil {
			log.Fatal(err)
		}
		ms = append(ms, m)
	}

	return ms
}
