package main

import (
	"datalogger/database"
	"datalogger/viewer3270"
	"datalogger/viewerhttp"
	"flag"
	"fmt"
	"sync"
)

func main() {
	httpPort := flag.Int("port", 8080, "Port on which JSON API and HTML is served")
	tn3270Port := flag.Int("tn3270-port", 3270, "Port on which JSON API and HTML is served")
	flag.Parse()

	db := database.CreateDatabaseReaderConnection()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		viewer3270.Start3270Server(db, *tn3270Port)
	}()

	go func() {
		defer wg.Done()
		viewerhttp.StartHTTPServer(db, *httpPort)
	}()

	fmt.Println("Servers listening - press Ctrl+C to stop")

	wg.Wait()
}
