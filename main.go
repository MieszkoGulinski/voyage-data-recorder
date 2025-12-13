package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

  go func(){
		defer wg.Done()
		start3270Server()
	}()

	go func(){
		defer wg.Done()
		startHTTPServer()
	}()

	fmt.Println("Servers listening - press Ctrl+C to stop")

	wg.Wait()
}