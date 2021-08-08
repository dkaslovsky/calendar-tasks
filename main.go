package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dkaslovsky/reminders/pkg/reminders"
)

func main() {
	now := time.Now()
	//rs := days()
	rs := months()

	for _, r := range rs {
		fmt.Println(r)
	}

	upcoming := make(map[int][]string)
	for _, r := range rs {
		appendUpcoming(upcoming, r, now)
	}

	for until, text := range upcoming {
		fmt.Println(until)
		for _, t := range text {
			fmt.Printf("\t%s\n", t)
		}
	}
}

type reminder interface {
	DaysFrom(time.Time) int
	GetText() string
}

func appendUpcoming(upcoming map[int][]string, r reminder, now time.Time) {
	until := r.DaysFrom(now)
	if _, ok := upcoming[until]; !ok {
		upcoming[until] = []string{}
	}
	upcoming[until] = append(upcoming[until], r.GetText())
}

func days() []*reminders.Daily {
	fileName := os.Args[1]
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(b), "\n")

	ds := []*reminders.Daily{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		d, err := reminders.NewDaily(line)
		if err != nil {
			log.Fatal(err)
		}
		ds = append(ds, d)
	}

	return ds
}

func months() []*reminders.Monthly {
	fileName := os.Args[1]
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(b), "\n")

	ms := []*reminders.Monthly{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		m, err := reminders.NewMonthly(line)
		if err != nil {
			log.Fatal(err)
		}
		ms = append(ms, m)
	}

	return ms
}
