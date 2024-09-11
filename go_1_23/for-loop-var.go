package main

import (
	"log"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// go1.23 : no need to capture i
			log.Printf("%d ", i)
		}()
	}
	wg.Wait()

	log.Println("done")
}
