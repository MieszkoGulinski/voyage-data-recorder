package main

import (
	"fmt"
	"log"
	"net"

	"github.com/racingmars/go3270"
	"gorm.io/gorm"
)

func onConnect3270(conn net.Conn, db *gorm.DB) {
	defer conn.Close()
	devinfo, err := go3270.NegotiateTelnet(conn)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Successfully created a 3270 connection")
	fmt.Println(devinfo.TerminalType())

	// Main loop
	currentScreenId := 0
	screensCount := len(screenTitles)
	var historyStack []int64

	for {
		screenContent, lastTimestamp, nextPageExists, err := getLogger3270ScreenContent(currentScreenId, historyStack, db)
		if err != nil {
			log.Println(err)
			return
		}
		// As we don't read user input, only the function keys, we set fieldValues to nil
		response, err := go3270.ShowScreenOpts(screenContent, nil, conn,
			go3270.ScreenOpts{
				Codepage: devinfo.Codepage(),
			})
		if err != nil {
			log.Println(err)
			return
		}
		// Respond to keys
		if response.AID == go3270.AIDPF3 {
			fmt.Println("Closing connection")
			return
		}

		if response.AID == go3270.AIDPF7 && lastTimestamp != 0 && nextPageExists {
			historyStack = append(historyStack, lastTimestamp)
		}

		if response.AID == go3270.AIDPF8 && len(historyStack) > 0 {
			historyStack = historyStack[:len(historyStack)-1]
		}

		if response.AID == go3270.AIDPF9 {
			historyStack = historyStack[:0] // clears the slice
			currentScreenId = (currentScreenId + 1) % screensCount
		}

		if response.AID == go3270.AIDPF10 {
			historyStack = historyStack[:0]
			currentScreenId = (currentScreenId - 1 + screensCount) % screensCount
		}
	}
}

func start3270Server(db *gorm.DB) {
	ln, err := net.Listen("tcp", ":3270")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go onConnect3270(conn, db)
	}
}