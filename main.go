package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/dkaslovsky/calendar-tasks/pkg/tasks"
)

func main() {
	loader := tasks.NewLoader()

	err := loader.Load(
		[]string{os.Args[1]},
		[]string{os.Args[2]},
		[]string{os.Args[3], os.Args[4]},
	)
	if err != nil {
		fmt.Println("XXXX")
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range loader.Ch {
			fmt.Println(t)
		}
	}()

	if err := loader.Wait(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	// now := time.Date(2021, 8, 14, 18, 0, 0, 0, time.Local)
	// grouper := tasks.NewGrouper(now)
	// grouper.Add(taskCh, doneCh, nReaders)

	// n, _ := strconv.Atoi(os.Args[5])
	// tasksGroups := grouper.Filter(n)
	// for day := 0; day <= n; day++ {
	// 	tasks, ok := tasksGroups[day]
	// 	if !ok {
	// 		continue
	// 	}
	// 	fmt.Printf("Day = %d\n", day)
	// 	for _, task := range tasks {
	// 		fmt.Println(task)
	// 	}
	// }
}
