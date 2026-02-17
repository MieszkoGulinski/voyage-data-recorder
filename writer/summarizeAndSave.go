package writer

import (
	"datalogger/aggregation"
	"datalogger/channels"
	"time"

	"gorm.io/gorm"
)

// summarizeAndSave saves buffered data to the database, and clears the buffers.
func summarizeAndSave(buffersSet *channels.BuffersSet, db *gorm.DB, diagnostics bool) error {
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
