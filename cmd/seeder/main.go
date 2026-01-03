package main

import (
	"datalogger/database"
	"log"
	"math/rand/v2"
	"os"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
)

func main() {
	shouldInsertData := false
	shouldInsertNull := false
	shouldInsertOutOfBounds := false
	samplesCount := 0

	promptInsertData := &survey.Confirm{
		Message: "Insert example data?",
	}
	survey.AskOne(promptInsertData, &shouldInsertData)

	if (shouldInsertData) {
		promptInsertNull := &survey.Confirm{
			Message: "Insert NULLs to emulate sensor faults?",
		}
		survey.AskOne(promptInsertNull, &shouldInsertNull)

		
		promptInsertOutOfBounds := &survey.Confirm{
			Message: "Insert values in warning/danger levels?",
		}
		survey.AskOne(promptInsertOutOfBounds, &shouldInsertOutOfBounds)

		promptSamplesCount := &survey.Input{
			Message: "Samples count:",
		}
		survey.AskOne(promptSamplesCount, &samplesCount)
	}

	oldDB := database.CreateDatabaseWriterConnection()

	// Checkpoint WAL to flush all data
	if err := oldDB.Exec("PRAGMA wal_checkpoint(FULL);").Error; err != nil {
		log.Fatal(err)
	}

	// Backup the current database
	filename := time.Now().UTC().Format("20060102-150405") + ".sqlite"
	err := oldDB.Exec("VACUUM INTO '" + filename + "'").Error
	if err != nil {
		log.Fatal(err)
	}

  // Close and remove the current database
	sqlDB, err := oldDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.Close()

	os.Remove("db.sqlite")
	os.Remove("db.sqlite-wal")
	os.Remove("db.sqlite-shm")

	// Regenerate the database

	db := database.CreateDatabaseWriterConnection()
	db.AutoMigrate(&database.Weather{}, &database.Position{}, &database.Battery{})

	if (!shouldInsertData) {
		return
	}

	// At first, seed the positions table
	lat := 54.72
	
	for i := range samplesCount {
		// Update position table
		sog := 1.0 + rand.Float64() * 0.5
		cog := 0.0 + rand.Float64() * 0.5
		bearing := 0.0 + rand.Float64() * 1
		sow := 1.0 + rand.Float64() * 0.5
		lon := 16.85 + rand.Float64() * 0.001

		pos := database.Position{
			Timestamp: time.Now().Unix() + int64(i*60),
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

		lat = lat + sog / (60 * 60) // realistic change in latitude assuming COG 0 degrees

		// Update 
	}

	// TODO nulls and warning/danger values
}
