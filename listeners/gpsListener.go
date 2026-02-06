package listeners

import (
	"context"
	"fmt"

	"github.com/stratoberry/go-gpsd"
)

func StartGPSListener(ctx context.Context, port int, diagnostics bool) error {
	gps, err := gpsd.Dial(fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return fmt.Errorf("Error setting up GPS listener on port %d: %w", port, err)
	}
	if diagnostics {
		fmt.Printf("GPS listener starting on port %d\n", port)
	}

	gps.AddFilter("TPV", func(r any) {
		msg, ok := r.(*gpsd.TPVReport)
		if !ok {
			return
		}
		if diagnostics {
			// TODO replace with proper logging
			fmt.Printf("TPV: %v\n", msg)
		}
	})

	doneCh := gps.Watch()
	select {
	case <-ctx.Done():
		gps.Close()
		return nil
	case <-doneCh:
		gps.Close()
		return nil
	}
}
