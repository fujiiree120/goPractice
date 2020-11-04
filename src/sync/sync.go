package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println(time.Now())
	count := 0
	var mu sync.Mutex
	//	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		//	wg.Add(1)
		go func() {
			mu.Lock()

			count++
			fmt.Println(i)
			//	wg.Done()
			defer mu.Unlock()
		}()
	}

	//	wg.Wait()
	fmt.Println()
	fmt.Printf("\n%d\n", count)
	fmt.Println(time.Now())
}
