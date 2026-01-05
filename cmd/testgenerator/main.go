package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/brutella/can"
)

func main() {
	interfaceName := flag.String("interface", "vcan0", "CAN interface name to which test data will be sent")
	bus, err := can.NewBusForInterfaceWithName(*interfaceName)
	if err != nil {
		log.Fatal(err)
	}

	// ConnectAndPublish blocks, we must call it in a new goroutine
	go func() {
    if err := bus.ConnectAndPublish(); err != nil {
        log.Fatal(err)
    }
	}()

	defer bus.Disconnect()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Submitting weather frame")

		var payload [8]uint8

		// Create values
		var temperature int16 = 331 // 33.1 C
		var pressure uint16 = 9981 // 998.1 hPa
		var appWindSpeed uint8 = 4 // 4 kt
		var appWindDir uint8 = 1 // 0=N, 1=NbE, 2=NNE, ..., 31=NbW
		var humidity uint8 = 86 // 86%

		// Form packet
		binary.BigEndian.PutUint16(payload[0:2], uint16(temperature)) // int16 must be converted to uint16 to be saved to 
		binary.BigEndian.PutUint16(payload[2:4], pressure)
		payload[4] = appWindSpeed
		payload[5] = appWindDir
		payload[6] = humidity
		// LATER: payload[7] indicates sensors fault, all zeros means valid reading

		frame := can.Frame{
			ID:   0x050,
			Data: payload,
		}
		err = bus.Publish(frame)
		if err != nil {
			log.Fatal(err)
		}
	}
}
