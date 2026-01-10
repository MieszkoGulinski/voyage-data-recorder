package main

import (
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

func sendGPSTestData(port int) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("GPSD test data generator listening on " + addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		log.Println("Client connected:", conn.RemoteAddr())
		go sendGPSExampleToClient(conn)
	}
}

func sendGPSExampleToClient(conn net.Conn) {
	defer conn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
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
			log.Println("write error:", err)
			return
		}
	}
}
