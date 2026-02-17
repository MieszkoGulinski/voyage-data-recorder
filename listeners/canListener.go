package listeners

import (
	"context"
	"datalogger/writer"
	"encoding/binary"
	"fmt"

	"github.com/brutella/can"
)

// StartCANListener starts a CAN listener on the given interface and passes received frames to the appropriate channels.
func StartCANListener(ctx context.Context, interfaceName string, diagnostics bool, channelsSet *writer.ChannelsSet) error {
	if diagnostics {
		fmt.Printf("CAN listener starting on interface %s\n", interfaceName)
	}
	bus, err := can.NewBusForInterfaceWithName(interfaceName)
	if err != nil {
		return fmt.Errorf("Error setting up CAN listener on interface %s: %w", interfaceName, err)
	}

	bus.SubscribeFunc(func(frame can.Frame) {
		// TODO replace this with actual processing
		if err := decodeCANFrame(frame, channelsSet); err != nil {
			fmt.Printf("Error decoding CAN frame: %v\n", err)
		}
	})

	errCh := make(chan error, 1)

	go func() {
		errCh <- bus.ConnectAndPublish()
	}()

	select {
	case <-ctx.Done():
		if diagnostics {
			fmt.Println("CAN listener shutting down")
		}
		bus.Disconnect()
		return nil

	case err := <-errCh:
		return err
	}
}

// decodeCANFrame decodes a received frame to appropriate struct, and passes it to appropriate channel.
//
// Unknown frame types are ignored.
//
// Some values inside the CAN frame represent signed integers. Since bit operations are performed assuming
// unsigned integers, we need to bit cast them to signed.
func decodeCANFrame(frame can.Frame, channelsSet *writer.ChannelsSet) error {
	switch frame.ID {
	case 0x050:
		// Weather frame
		temperature := int16(binary.BigEndian.Uint16(frame.Data[0:2]))
		pressure := binary.BigEndian.Uint16(frame.Data[2:4])
		appWindSpeed := frame.Data[4]
		appWindDir := frame.Data[5]
		humidity := frame.Data[6]
		// TODO error bits from Data[7] - if any, set appropriate fields to nil
		channelsSet.WeatherCh <- writer.WeatherMessage{
			Timestamp:    nil,
			Temperature:  &temperature,
			Pressure:     &pressure,
			AppWindSpeed: &appWindSpeed,
			AppWindDir:   &appWindDir,
			Humidity:     &humidity,
		}
	default:
		// ignore unknown frame types
	}
	return nil
}
