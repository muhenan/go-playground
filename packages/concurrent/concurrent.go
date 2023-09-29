package concurrent

import (
	"fmt"
	"sync"
)

func Concurrent() {
	var wg sync.WaitGroup
	ch := make(chan int)

	// Start two goroutines
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			defer func() {
				fmt.Println("goroutine ", id, " completed")
				wg.Done()
			}()
			for j := 1; j <= 5; j++ {
				ch <- id*10 + j
				fmt.Println(id, " ", id*10+j)
			}
		}(i)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Read from the channel
	for num := range ch {
		fmt.Println("Received:", num)
	}
}
