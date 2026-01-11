package testgenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// TPVMessage represents the GPSD TPV JSON structure
type TPVMessage struct {
	Class string  `json:"class"`
	Mode  int     `json:"mode"`
	Time  string  `json:"time"`
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Alt   float64 `json:"alt"`
	Speed float64 `json:"speed"`
	Climb float64 `json:"climb"`
}

// generateTPV returns a sample TPV message
func generateTPV() TPVMessage {
	return TPVMessage{
		Class: "TPV",
		Mode:  3, // 3D fix
		Time:  time.Now().UTC().Format(time.RFC3339Nano),
		Lat:   54.306, // latitude
		Lon:   15.753, // longitude
		Alt:   0.5,    // meters
		Speed: 1.2,    // m/s
		Climb: 0.0,    // m/s
	}
}

func StartGPSTestDataGenerator(ctx context.Context, port int) error {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Error opening GPSD TCP interface: %w", err)
	}
	defer listener.Close()

	fmt.Println("GPSD test data generator listening on " + addr)

	// Stop listener when context is canceled
	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err() // graceful shutdown
			default:
				log.Println("accept error:", err)
				continue
			}
		}
		fmt.Println("Client connected:", conn.RemoteAddr())
		go sendGPSExampleToClient(ctx, conn)
	}
}

func sendGPSExampleToClient(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("GPS client disconnecting due to shutdown:", conn.RemoteAddr())
			return

		case <-ticker.C:
			sendSingleGPSFrameToClient(conn)
		}
	}
}

func sendSingleGPSFrameToClient(conn net.Conn) {
	fmt.Println("Submitting test GPS TPV frame")

	tpv := generateTPV()
	data, err := json.Marshal(tpv)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	// gpsd clients expect newline-delimited JSON
	data = append(data, '\n')

	_, err = conn.Write(data)
	if err != nil {
		log.Println("TCP write error:", err)
		return
	}
}
