package listeners

import (
	"context"
	"fmt"

	"github.com/brutella/can"
	"gorm.io/gorm"
)

type CANFrame struct {
	DB *gorm.DB
}

func StartCANListener(ctx context.Context, interfaceName string, diagnostics bool) error {
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

	bus.ConnectAndPublish()

	return nil
}
