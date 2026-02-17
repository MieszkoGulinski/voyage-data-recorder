package writer

import (
	"datalogger/aggregation"
	"datalogger/channels"
	"time"

	"gorm.io/gorm"
)

func summarizeAndSaveSensors(buffersSet *channels.BuffersSet, db *gorm.DB, diagnostics bool) error {
	weather := aggregation.SummarizeWeather(buffersSet.WeatherBuffer, time.Now().Unix())
	// ... more types ...

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&weather).Error; err != nil {
			return err
		}
		// ... more types ...
		return nil
	})
	if err != nil {
		return err
	}

	// Done - buffers can be cleared
	clearBuffers(buffersSet)
	return nil
}

func summarizeAndSavePositions(buffersSet *channels.BuffersSet, db *gorm.DB, diagnostics bool) error {
	// TODO
	return nil
}
