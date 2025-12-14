package main

import (
	"fmt"
	"sync"
)

func main() {
	db := createDatabaseConnection()

	var wg sync.WaitGroup
	wg.Add(2)

  go func(){
		defer wg.Done()
		start3270Server(db)
	}()

	go func(){
		defer wg.Done()
		startHTTPServer(db)
	}()

	fmt.Println("Servers listening - press Ctrl+C to stop")

	wg.Wait()
}