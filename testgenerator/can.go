package testgenerator

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/brutella/can"
	"golang.org/x/sync/errgroup"
)

func StartCANTestDataGenerator(ctx context.Context, interfaceName string) error {
	bus, err := can.NewBusForInterfaceWithName(interfaceName)
	if err != nil {
		return fmt.Errorf("Error opening CAN interface %q: %w", interfaceName, err)
	}

	defer bus.Disconnect()
	g, ctx := errgroup.WithContext(ctx)

	// ConnectAndPublish blocks, we must call it in a new goroutine
	g.Go(func() error {
		return bus.ConnectAndPublish()
	})

	g.Go(func() error {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case <-ticker.C:
				if err := submitFrame(bus); err != nil {
					return fmt.Errorf("publish CAN frame: %w", err)
				}
			}
		}
	})

	return g.Wait()
}

func submitFrame(bus *can.Bus) error {
	fmt.Println("Submitting test weather frame")

	var payload [8]uint8

	// Create values
	var temperature int16 = 331 // 33.1 C
	var pressure uint16 = 9981  // 998.1 hPa
	var appWindSpeed uint8 = 4  // 4 kt
	var appWindDir uint8 = 1    // 0=N, 1=NbE, 2=NNE, ..., 31=NbW
	var humidity uint8 = 86     // 86%

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
	err := bus.Publish(frame)
	if err != nil {
		return fmt.Errorf("Error submitting frame to CAN: %w", err)
	}

	return nil
}
