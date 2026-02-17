package writer

import (
	"context"
	"datalogger/channels"
	"datalogger/database"
	"time"

	"github.com/stratoberry/go-gpsd"
	"gorm.io/gorm"
)

func StartWriter(ctx context.Context, db *gorm.DB, channelsSet channels.ChannelsSet, diagnostics bool) error {
	sensorTicker := time.NewTicker(1 * time.Minute)
	defer sensorTicker.Stop()

	// Positions need to be flushed more frequently than other sensors
	// on abrupt changes of course or speed
	positionTicker := time.NewTicker(10 * time.Second)
	defer positionTicker.Stop()

	buffers := &channels.BuffersSet{}
	initializeBuffers(buffers)

	for {
		select {
		case <-ctx.Done():
			return nil // shutdown
		case weather := <-channelsSet.WeatherCh:
			buffers.WeatherBuffer = append(buffers.WeatherBuffer, weather)
		case battery := <-channelsSet.BatteryCh:
			buffers.BatteryBuffer = append(buffers.BatteryBuffer, battery)
		case compass := <-channelsSet.CompassCh:
			buffers.CompassBuffer = append(buffers.CompassBuffer, compass)
		case gps := <-channelsSet.GPSCh:
			buffers.GPSBuffer = append(buffers.GPSBuffer, gps)
		// The following readouts are not buffered and processed, but directly written to DB
		case position := <-channelsSet.PositionCh:
			writePosition(db, position)
		case navtex := <-channelsSet.NavtexCh:
			writeNavtex(db, navtex)
		case textNote := <-channelsSet.TextNoteCh:
			writeTextNote(db, textNote)
		case <-sensorTicker.C:
			if err := summarizeAndSaveSensors(buffers, db, diagnostics); err != nil {
				return err
			}
		}
	}
}

// initializeBuffers creates preallocated buffers for sensor readouts
// TODO: make it a method of BuffersSet
// TODO: adjust sizes to expected data rates
func initializeBuffers(buffersSet *channels.BuffersSet) {
	buffersSet.WeatherBuffer = make([]channels.WeatherMessage, 0, 100)
	buffersSet.BatteryBuffer = make([]channels.BatteryMessage, 0, 100)
	buffersSet.CompassBuffer = make([]channels.CompassMessage, 0, 100)
	buffersSet.GPSBuffer = make([]gpsd.TPVReport, 0, 100)
}

// clearBuffers clears the buffers and preallocates slices for sensor readouts
// TODO: make it a method of BuffersSet
func clearBuffers(buffersSet *channels.BuffersSet) {
	buffersSet.WeatherBuffer = buffersSet.WeatherBuffer[:0]
	buffersSet.BatteryBuffer = buffersSet.BatteryBuffer[:0]
	buffersSet.CompassBuffer = buffersSet.CompassBuffer[:0]
	buffersSet.GPSBuffer = buffersSet.GPSBuffer[:0]
}

func writePosition(db *gorm.DB, position channels.PositionMessage) {
	db.Create(&database.Position{
		Timestamp: time.Now().Unix(),
		SourceId:  2,
		Latitude:  &position.Latitude,
		Longitude: &position.Longitude,
	})
}

func writeNavtex(db *gorm.DB, message string) {
	db.Create(&database.NavtexMessages{
		Timestamp: time.Now().Unix(),
		Text:      message,
	})
}

func writeTextNote(db *gorm.DB, textNote string) {
	db.Create(&database.TextNotes{
		Timestamp: time.Now().Unix(),
		Text:      textNote,
	})
}
