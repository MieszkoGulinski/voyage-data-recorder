package main

import (
	"datalogger/database"
	"log"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
)

func main() {
	db := database.CreateDatabaseWriterConnection()

	shouldInsertNull := false
	promptInsertNull := &survey.Confirm{
		Message: "Insert NULLs to emulate sensor faults?",
	}
	survey.AskOne(promptInsertNull, &shouldInsertNull)

	shouldInsertOutOfBounds := false
	promptInsertOutOfBounds := &survey.Confirm{
		Message: "Insert values in warning/danger levels?",
	}
	survey.AskOne(promptInsertOutOfBounds, &shouldInsertOutOfBounds)

	// At first, seed the positions table
	
	lat := 54.51
	lon := 18.56
	timestamp := time.Now().Unix()
	sog := 1.1
	cog := 90.0
	bearing := 88.0
	sow := 1.2

	positionsCount := 40
	for i := range positionsCount {
		pos := database.Position{
        Timestamp: time.Now().Unix() + int64(i),
        SourceId:  1,
        Latitude:  &lat,
        Longitude: &lon,
        SpeedOverGround: &sog,
				CourseOverGround: &cog,
				MagneticBearing: &bearing,
				SpeedOverWater: &sow,
    }

		err := db.Create(&pos).Error
		if err != nil {
			log.Fatal(err)
		}

		lon = lon + 0.01
		timestamp = timestamp + 120 // 2 minutes
	}

	// TODO more tables
	// TODO nulls and warning/danger values
}
