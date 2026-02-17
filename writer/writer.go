package writer

import (
	"context"
	"time"

	"gorm.io/gorm"
)

func StartWriter(ctx context.Context, db *gorm.DB, channelsSet *ChannelsSet) error {
	interval := 1 * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	buffers := BuffersSet{}

	for {
		select {
		case <-ctx.Done():
			return nil // shutdown
		case weather := <-channelsSet.WeatherCh:
			buffers.WeatherBuffer = append(buffers.WeatherBuffer, weather)
		case gps := <-channelsSet.GPSCh:
			buffers.GPSBuffer = append(buffers.GPSBuffer, gps)
		case navtex := <-channelsSet.NavtexCh:
			buffers.NavtexBuffer = append(buffers.NavtexBuffer, navtex)
		case position := <-channelsSet.PositionCh:
			buffers.PositionBuffer = append(buffers.PositionBuffer, position)
		case textNote := <-channelsSet.TextNoteCh:
			buffers.TextNoteBuffer = append(buffers.TextNoteBuffer, textNote)
		case <-ticker.C:
			// TODO summarize data
		}
	}
}
