package seeder

import (
	"datalogger/database"
	"log"
	"math/rand/v2"
	"time"

	"gorm.io/gorm"
)

func FillDatabase(db *gorm.DB, options SeederOptions) {
	if !options.InsertData {
		return
	}

	// At first, seed the positions table
	lat := 54.72
	
	for i := range options.SamplesCount {
		timestamp := time.Now().Unix() + int64(i*120)
		// Update position table
		sog := 1.0 + rand.Float64() * 0.5
		cog := 0.0 + rand.Float64() * 0.5
		bearing := 0.0 + rand.Float64() * 1
		sow := 1.0 + rand.Float64() * 0.5
		lon := 16.85 + rand.Float64() * 0.001

		positionRow := database.Position{
			Timestamp: timestamp,
			SourceId:  1,
			Latitude:  &lat,
			Longitude: &lon,
			SpeedOverGround: &sog,
			CourseOverGround: &cog,
			MagneticBearing: &bearing,
			SpeedOverWater: &sow,
    }

		err := db.Create(&positionRow).Error
		if err != nil {
			log.Fatal(err)
		}

		lat = lat + sog / (60 * 60) // realistic change in latitude assuming COG 0 degrees

		// Update weather table
		airTemp := 10 + rand.Float64() * 10
		waterTemp := 5 + rand.Float64() * 10
		pressure := 960 + rand.Float64() * 80 // 960 - 1040 hPa
		var awd uint8 = 0
		aws := 0 + rand.Float64()
		awgs := aws + rand.Float64()
		var twd uint8 = 1
		tws := aws + sog
		twgs := tws + rand.Float64()

		weatherRow:= database.Weather{
			Timestamp: timestamp,
			AirTemperature: &airTemp,
			WaterTemperature: &waterTemp,
			Pressure: &pressure,
			ApparentWindDirection: &awd,
			ApparentWindSpeed: &aws,
			ApparentWindGustSpeed: &awgs,
			TrueWindDirection: &twd,
			TrueWindSpeed: &tws,
			TrueWindGustSpeed: &twgs,
		}
		err = db.Create(&weatherRow).Error
		if err != nil {
			log.Fatal(err)
		}

		// Update battery table
		batteryRow := database.Battery {
			Timestamp: timestamp,
			Charge: 99,
			Voltage: 24 + rand.Float64(),
		}
		err = db.Create(&batteryRow).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	// TODO nulls and warning/danger values
}