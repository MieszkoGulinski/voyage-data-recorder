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
		{Row: 0, Col: 0, Content: "Positions", Color: go3270.Green},
		// Header
		{Row: 1, Col: 0, Content: "Time UTC"},
		{Row: 1, Col: 22, Content: "Lat"},
		{Row: 1, Col: 35, Content: "Lon"},
		{Row: 1, Col: 48, Content: "SOG"},
		{Row: 1, Col: 55, Content: "COG"},
		{Row: 1, Col: 63, Content: "SOW"},
		{Row: 1, Col: 71, Content: "Mag.b."},
	}

	for i, v := range result {
		screenContent = append(
			screenContent,
			go3270.Field{Row: i+2, Col: 0, Content: formatters.FormatUnixTimestamp(v.Timestamp), Color: go3270.White},
			go3270.Field{Row: i+2, Col: 22, Content: formatters.FormatLatitude(v.Latitude), Color: formatters.Format3270Color(v.Latitude)},
			go3270.Field{Row: i+2, Col: 35, Content: formatters.FormatLongitude(v.Longitude), Color: formatters.Format3270Color(v.Longitude)},
			go3270.Field{Row: i+2, Col: 48, Content: formatters.FormatNumber("%4.1f", v.SpeedOverGround), Color: formatters.Format3270Color(v.SpeedOverGround)},
			go3270.Field{Row: i+2, Col: 55, Content: formatters.FormatNumber("%05.1f", v.CourseOverGround), Color: formatters.Format3270Color(v.CourseOverGround)},
			go3270.Field{Row: i+2, Col: 63, Content: formatters.FormatNumber("%05.1f", v.SpeedOverWater), Color: formatters.Format3270Color(v.SpeedOverWater)},
			go3270.Field{Row: i+2, Col: 73, Content: formatters.FormatNumber("%05.1f", v.MagneticBearing), Color: formatters.Format3270Color(v.MagneticBearing)},
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
		{Row: 0, Col: 0, Content: "Weather", Color: go3270.Green},
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
			go3270.Field{Row: i+2, Col: 0, Content: formatters.FormatUnixTimestamp(v.Timestamp), Color: go3270.White},
			go3270.Field{Row: i+2, Col: 23, Content: formatters.FormatNumber("%5.1f", v.AirTemperature), Color: formatters.Format3270Color(v.AirTemperature)},
			go3270.Field{Row: i+2, Col: 28, Content: formatters.FormatNumber("%5.1f", v.WaterTemperature), Color: formatters.Format3270Color(v.WaterTemperature)},
			go3270.Field{Row: i+2, Col: 33, Content: formatters.FormatNumber("%4.0f", v.Pressure), Color: formatters.Format3270Color(v.Pressure)},
			go3270.Field{Row: i+2, Col: 38, Content: formatters.FormatNumber("%5.1f", v.WindSpeed), Color: formatters.Format3270Color(v.WindSpeed)},
			go3270.Field{Row: i+2, Col: 43, Content: formatters.FormatNumber("%5.1f", v.AirTemperature), Color: formatters.Format3270Color(v.AirTemperature)},
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
		{Row: 0, Col: 0, Content: "Battery",  Color: go3270.Green},
		// Header
		{Row: 1, Col: 0, Content: "Time UTC"},
		{Row: 1, Col: 23, Content: "%"},
		{Row: 1, Col: 28, Content: "%/hr"},
	}

	for i, v := range result {
		screenContent = append(
			screenContent,
			go3270.Field{Row: i+2, Col: 0, Content: formatters.FormatUnixTimestamp(v.Timestamp), Color: go3270.White},
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