package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Teamlead explanation: The final answer is not 1000 because the counter++ operation is not atomic (it consists of separate read, increment, and write steps), leading to a data race where multiple goroutines concurrently overwrite each other's results.")
	fmt.Println("--------------------------------------------------")
	var counter int
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}