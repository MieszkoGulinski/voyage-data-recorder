package database

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateDatabaseReaderConnection() *gorm.DB {
	dsn := "file:db.sqlite?mode=ro&_journal_mode=WAL&_busy_timeout=5000"
	
	for {
		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}

		log.Println("Error opening DB:", err)
		log.Println("Will retry in 5 s")
		time.Sleep(5 * time.Second)
	}
}

func CreateDatabaseWriterConnection() *gorm.DB {
	dsn := "file:db.sqlite?_journal_mode=WAL&_busy_timeout=5000"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type WithTimestamp interface {
	GetTimestamp() int64
}

func QueryWithPagination[T WithTimestamp](
	db *gorm.DB,
	model *T,
	lastTimestamp int64,
	limit int,
) (
	result []T,
	newLastTimestamp int64,
	nextPageExists bool,
	err error,
) {

	query := db.Model(model)

	if lastTimestamp != 0 {
		query = query.Where("timestamp < ?", lastTimestamp)
	}

	err = query.
		Order("timestamp DESC").
		Limit(limit).
		Find(&result).
		Error

	nextPageExists = len(result) == limit
	
	if (len(result) > 0) {
		newLastTimestamp = result[len(result)-1].GetTimestamp()
	}

	return
}