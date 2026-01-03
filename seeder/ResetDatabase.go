package seeder

import (
	"datalogger/database"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

func ResetDatabase() *gorm.DB {
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

	return db
}