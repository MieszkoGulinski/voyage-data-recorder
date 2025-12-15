package viewer3270

import (
	"datalogger/database"
	"datalogger/formatters"

	"fmt"
	"log"

	"github.com/racingmars/go3270"
	"gorm.io/gorm"
)

var screenTitles = []string{
	"Position",
	"Weather",
	"Battery status",
}

// Individual formatters by table

func renderPositionScreen(lastTimestamp int64, db *gorm.DB) (
	screenContent go3270.Screen,
	newLastTimestamp int64,
	nextPageExists bool,
	err error,
){
	result, newLastTimestamp, nextPageExists, err := database.QueryWithPagination(db, &database.Position{}, lastTimestamp, 20)

	screenContent = go3270.Screen{
		{Row: 0, Col: 0, Content: "Positions", Intense: true},
		// Header
		{Row: 1, Col: 0, Content: "Time UTC"},
		{Row: 1, Col: 22, Content: "Lat"},
		{Row: 1, Col: 35, Content: "Lon"},
		{Row: 1, Col: 48, Content: "V kt"},
		{Row: 1, Col: 55, Content: "  Dir"},
		{Row: 1, Col: 63, Content: "Mag.b."},
	}

	for i, v := range result {
		screenContent = append(
			screenContent,
			go3270.Field{Row: i+2, Col: 0, Content: formatters.FormatUnixTimestamp(v.Timestamp)},
			go3270.Field{Row: i+2, Col: 22, Content: formatters.FormatLatitude(v.Latitude)},
			go3270.Field{Row: i+2, Col: 35, Content: formatters.FormatLongitude(v.Longitude)},
			go3270.Field{Row: i+2, Col: 48, Content: formatters.FormatNumber("%4.1f", v.SpeedOverGround)},
			go3270.Field{Row: i+2, Col: 55, Content: formatters.FormatNumber("%05.1f", v.CourseOverGround)},
			go3270.Field{Row: i+2, Col: 63, Content: formatters.FormatNumber("%05.1f", v.SpeedOverWater)},
			go3270.Field{Row: i+2, Col: 73, Content: formatters.FormatNumber("%05.1f", v.MagneticBearing)},
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
	result, newLastTimestamp, nextPageExists, err := database.QueryWithPagination(db, &database.Weather{}, lastTimestamp, 20)

	screenContent = go3270.Screen{
		{Row: 0, Col: 0, Content: "Weather", Intense: true},
		// Header
		{Row: 1, Col: 0, Content: "Time UTC"},
		{Row: 1, Col: 23, Content: "Air C"},
		{Row: 1, Col: 28, Content: "Wat C"},
		{Row: 1, Col: 33, Content: "P hPa"},
		{Row: 1, Col: 38, Content: "Sun"},
		{Row: 1, Col: 43, Content: "W m/s"},
		{Row: 1, Col: 48, Content: "W dir"},
	}

	for i, v := range result {
		screenContent = append(
			screenContent,
			go3270.Field{Row: i+2, Col: 0, Content: formatters.FormatUnixTimestamp(v.Timestamp)},
			go3270.Field{Row: i+2, Col: 23, Content: formatters.FormatNumber("%5.1f", v.AirTemperature)},
			go3270.Field{Row: i+2, Col: 28, Content: formatters.FormatNumber("%5.1f", v.WaterTemperature)},
			go3270.Field{Row: i+2, Col: 33, Content: formatters.FormatNumber("%4.0f", v.Pressure)},
			go3270.Field{Row: i+2, Col: 38, Content: formatters.FormatNumber("%5.1f", v.WindSpeed)},
			go3270.Field{Row: i+2, Col: 43, Content: formatters.FormatNumber("%5.1f", v.AirTemperature)},
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
	result, newLastTimestamp, nextPageExists, err := database.QueryWithPagination(db, &database.Battery{}, lastTimestamp, 20)

	screenContent = go3270.Screen{
		{Row: 0, Col: 0, Content: "Battery", Intense: true},
		// Header
		{Row: 1, Col: 0, Content: "Time UTC"},
		{Row: 1, Col: 23, Content: "%"},
		{Row: 1, Col: 28, Content: "%/hr"},
	}

	for i, v := range result {
		screenContent = append(
			screenContent,
			go3270.Field{Row: i+2, Col: 0, Content: formatters.FormatUnixTimestamp(v.Timestamp)},
			go3270.Field{Row: i+2, Col: 23, Content: fmt.Sprintf("%3d", v.Percent)},
			go3270.Field{Row: i+2, Col: 28, Content: fmt.Sprintf("%5.1f", v.ChangeRate)},
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
		go3270.Field{Row: 23, Col: 0, Content: "F3 close F7 page up F8 page down F9 next table F10 prev table", Color: go3270.Green},
	)

	return
}