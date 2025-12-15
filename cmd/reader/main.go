package main

import (
	"datalogger/database"
	"datalogger/viewer3270"
	"datalogger/viewerhttp"
	"fmt"
	"sync"
)

func main() {
	db := database.CreateDatabaseReaderConnection()

	var wg sync.WaitGroup
	wg.Add(2)

  go func(){
		defer wg.Done()
		viewer3270.Start3270Server(db)
	}()

	go func(){
		defer wg.Done()
		viewerhttp.StartHTTPServer(db)
	}()

	fmt.Println("Servers listening - press Ctrl+C to stop")

	wg.Wait()
}