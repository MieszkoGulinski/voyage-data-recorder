package main

import (
	"context"
	"datalogger/testgenerator"
	"flag"
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"
)

func main() {
	interfaceName := flag.String("interface", "vcan0", "CAN interface name to which test data will be sent")
	activateGPS := flag.Bool("gps", false, "Activate gpsd mock?")
	gpsPort := flag.Int("port", 2497, "When using --gps option, port on which gpsd mock will be sending data")
	flag.Parse()

	g, _ := errgroup.WithContext(context.Background()) // TODO: the second argument is context, pass it to goroutines below

	if *activateGPS {
		fmt.Println("Starting GPS test data generator")
		g.Go(func() error {
			return testgenerator.StartGPSTestDataGenerator(*gpsPort)
		})
	}

	g.Go(func() error {
		return testgenerator.StartCANTestDataGenerator(*interfaceName)
	})

	err := g.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
