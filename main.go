package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("main started")
	for i := 0; i < 4; i++ {
		go func(id int) {
			defer wg.Done()
			runtime.LockOSThread()
			for {
				fmt.Printf("Goroutine %d on thread\n", id)
				time.Sleep(10 * time.Second)
			}
		}(i)
	}
	wg.Wait()
}
