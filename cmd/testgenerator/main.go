package main

import (
	"datalogger/testgenerator"
	"flag"
	"fmt"
)

func main() {
	interfaceName := flag.String("interface", "vcan0", "CAN interface name to which test data will be sent")
	activateGPS := flag.Bool("gps", false, "Activate gpsd mock?")
	gpsPort := flag.Int("port", 2497, "When using --gps option, port on which gpsd mock will be sending data")
	flag.Parse()

	if *activateGPS {
		fmt.Println("Starting GPS test data generator")
		go testgenerator.StartGPSTestDataGenerator(*gpsPort)
	}

	testgenerator.StartCANTestDataGenerator(*interfaceName)
}
