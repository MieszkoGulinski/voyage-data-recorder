package main

import (
	"fmt"
	"log"
	"time"

	"github.com/racingmars/go3270"
	"gorm.io/gorm"
)

var screenTitles = []string{
	"Position",
	"Weather",
	"Battery status",
}

func formatUnixTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	return t.Format("2006-01-02 15:04:05")
}

// Individual formatters by table

func renderPositionScreen(lastTimestamp int64, db *gorm.DB) (
	screenContent go3270.Screen,
	newLastTimestamp int64,
	nextPageExists bool,
	err error,
){
	result, newLastTimestamp, nextPageExists, err := queryWithPagination(db, &Position{}, lastTimestamp, 20)

	screenContent = go3270.Screen{
		{Row: 0, Col: 0, Content: "Positions"},
		// Header
		{Row: 1, Col: 0, Content: "Time UTC"},
		{Row: 1, Col: 23, Content: "V kt"},
		{Row: 1, Col: 32, Content: "Dir"},
		{Row: 1, Col: 42, Content: "Mag.bear."},
	}

	for i, v := range result {
		screenContent = append(
			screenContent,
			go3270.Field{Row: i+2, Col: 0, Content: formatUnixTimestamp(v.Timestamp)},
			go3270.Field{Row: i+2, Col: 23, Content: fmt.Sprintf("%5.1f", v.Speed)},
			go3270.Field{Row: i+2, Col: 32, Content: fmt.Sprintf("%5.1f", v.Direction)},
			go3270.Field{Row: i+2, Col: 42, Content: fmt.Sprintf("%5.1f", v.MagBearing)},
		)
	}

	return
}

func renderWeatherScreen(lastTimestamp int64, db *gorm.DB) (
	screenContent go3270.Screen,
	newLastTimestamp int64,
	nextPageExists bool,
	err error,
){
	result, newLastTimestamp, nextPageExists, err := queryWithPagination(db, &Weather{}, lastTimestamp, 20)

	screenContent = go3270.Screen{
		{Row: 0, Col: 0, Content: "Weather"},
		// Header
		{Row: 1, Col: 0, Content: "Time UTC"},
		// ...
	}

	for i, v := range result {
		screenContent = append(
			screenContent,
			go3270.Field{Row: i+2, Col: 0, Content: formatUnixTimestamp(v.Timestamp)},
			go3270.Field{Row: i+2, Col: 23, Content: fmt.Sprintf("%5.1f", v.AirTemperature)},
		)
	}

	return
}

func renderBatteryScreen(lastTimestamp int64, db *gorm.DB) (
	screenContent go3270.Screen,
	newLastTimestamp int64,
	nextPageExists bool,
	err error,
){
	result, newLastTimestamp, nextPageExists, err := queryWithPagination(db, &Battery{}, lastTimestamp, 20)

	screenContent = go3270.Screen{
		{Row: 0, Col: 0, Content: "Battery"},
		// Header
		{Row: 1, Col: 0, Content: "Time UTC"},
		// ...
	}

	for i, v := range result {
		screenContent = append(
			screenContent,
			go3270.Field{Row: i+2, Col: 0, Content: formatUnixTimestamp(v.Timestamp)},
			go3270.Field{Row: i+2, Col: 23, Content: fmt.Sprintf("%3d", v.Percent)},
		)
	}

	return
}

func getLogger3270ScreenContent (
	currentScreenId int,
  historyStack []int64,
	db *gorm.DB,
) (
	screenContent go3270.Screen,
	lastTimestamp int64,
	nextPageExists bool,
	err error,
) {
	var currentPageLastTimestamp int64
	if (len(historyStack) > 0) {
		currentPageLastTimestamp = historyStack[len(historyStack)-1]
	}

	switch currentScreenId {
		case 0: screenContent, lastTimestamp, nextPageExists, err = renderPositionScreen(currentPageLastTimestamp, db)
		case 1: screenContent, lastTimestamp, nextPageExists, err = renderWeatherScreen(currentPageLastTimestamp, db)
		case 2: screenContent, lastTimestamp, nextPageExists, err = renderBatteryScreen(currentPageLastTimestamp, db)
	}
  
	if err != nil {
		log.Println(err)
		return
	}

	// Default size is 80x24 characters - TODO handle more screen sizes if available
	screenContent = append(screenContent,
		go3270.Field{Row: 23, Col: 0, Content: "F3 close F7 page up F8 page down F9 next table F10 prev table"},
	)

	return
}