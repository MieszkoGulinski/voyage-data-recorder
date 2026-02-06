package main

import (
	"context"
	"datalogger/database"
	"datalogger/listeners"
	"datalogger/writer"
	"flag"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"
)

func main() {
	canInterfaceName := flag.String("interface", "can0", "CAN interface name")
	gpsdPort := flag.Int("gpsd-port", 2947, "gpsd port")
	httpPort := flag.Int("http-port", 8081, "HTTP listener port")
	diagnostics := flag.Bool("diagnostics", false, "Enable diagnostics")
	flag.Parse()

	db := database.CreateDatabaseWriterConnection()

	g, ctx := errgroup.WithContext(context.Background())

	channelsSet := writer.NewChannelsSet()

	g.Go(func() error {
		return listeners.StartCANListener(ctx, *canInterfaceName, *diagnostics, channelsSet)
	})
	g.Go(func() error {
		return listeners.StartGPSListener(ctx, *gpsdPort, *diagnostics, channelsSet)
	})
	g.Go(func() error {
		return listeners.StartHTTPListener(ctx, *httpPort, *diagnostics, channelsSet)
	})

	g.Go(func() error {
		// Summarizer and writer to DB
		return writer.StartWriter(ctx, db, channelsSet)
	})

	if *diagnostics {
		fmt.Println("Writer active")
	}

	err := g.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
