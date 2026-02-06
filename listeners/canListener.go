package listeners

import (
	"context"
	"datalogger/writer"
	"fmt"

	"github.com/brutella/can"
)

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
		fmt.Printf("Frame: %v\n", frame)
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
