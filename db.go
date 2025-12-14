package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Add tables structure here
// GetTimestamp is needed in the queryWithPagination helper

type Weather struct {
	Timestamp        int64   `gorm:"column:timestamp;primaryKey"`
	AirTemperature   float64 `gorm:"column:air_temperature"`
	WaterTemperature float64 `gorm:"column:water_temperature"`
	Pressure         float64 `gorm:"column:pressure"`
	Sunlight         uint16  `gorm:"column:sunlight"`
	RawWindDirection uint16 `gorm:"column:raw_wind_direction"`
	RawWindSpeeed    float64 `gorm:"column:raw_wind_speed"`
	WindDirection    uint16 `gorm:"column:wind_direction"`
	WindSpeeed       float64 `gorm:"column:wind_speed"`
}

func (Weather) TableName() string {
	return "weather"
}

func (r Weather) GetTimestamp() int64 {
	return r.Timestamp
}

type Position struct {
	Timestamp  int64   `gorm:"column:timestamp;primaryKey"`
	Lat        float64 `gorm:"column:lat"`
	Lon        float64 `gorm:"column:lon"`
	Speed      float64 `gorm:"column:speed"`
	Direction  float64 `gorm:"column:direction"`
	MagBearing float64 `gorm:"column:mag_bearing"`
}

func (Position) TableName() string {
	return "positions"
}

func (r Position) GetTimestamp() int64 {
	return r.Timestamp
}

type Battery struct {
	Timestamp  int64   `gorm:"column:timestamp;primaryKey"`
	Percent    uint8   `gorm:"column:percent"`
	ChangeRate float64 `gorm:"column:change_rate"`
}

func (Battery) TableName() string {
	return "battery"
}

func (r Battery) GetTimestamp() int64 {
	return r.Timestamp
}

func createDatabaseConnection() *gorm.DB {
	dsn := "file:db.sqlite?mode=ro&_journal_mode=WAL&immutable=1"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}


// Reusable helper

type WithTimestamp interface {
	GetTimestamp() int64
}

func queryWithPagination[T WithTimestamp](
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