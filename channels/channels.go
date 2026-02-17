package channels

import "github.com/stratoberry/go-gpsd"

// Channels for data from listeners to the summarizer
// Each channel contains an individual data packet type, already decoded from its original format
// Note that a given packet type can contain multiple individual values (e.g. weather packet
// contains temperature, pressure, humidity etc.)
// See docs/input-format.md for details.

type ChannelsSet struct {
	WeatherCh  chan WeatherMessage
	BatteryCh  chan BatteryMessage
	CompassCh  chan CompassMessage
	GPSCh      chan gpsd.TPVReport // received from GPS
	NavtexCh   chan string
	PositionCh chan PositionMessage // for manually added positions, not from GPS
	TextNoteCh chan string
}

func NewChannelsSet() ChannelsSet {
	// Although the summarizer stores received data in a slice, and summarizes it every
	// specified time (default 1 minute), buffered channels are needed to not block the
	// listeners when the summarizer is busy (e.g. writing to SQLite).
	return ChannelsSet{
		WeatherCh:  make(chan WeatherMessage, 10),
		BatteryCh:  make(chan BatteryMessage, 10),
		CompassCh:  make(chan CompassMessage, 10),
		GPSCh:      make(chan gpsd.TPVReport, 10),
		NavtexCh:   make(chan string, 10),
		PositionCh: make(chan PositionMessage, 10),
		TextNoteCh: make(chan string, 10),
	}
}

type BuffersSet struct {
	WeatherBuffer  []WeatherMessage
	BatteryBuffer  []BatteryMessage
	CompassBuffer  []CompassMessage
	GPSBuffer      []gpsd.TPVReport
	NavtexBuffer   []string
	PositionBuffer []PositionMessage
	TextNoteBuffer []string
}

// Formats for individual data packets
// Nil field indicates sensor fault

type WeatherMessage struct {
	Temperature  *float32 // C
	Pressure     *float32 // hPa
	AppWindSpeed *uint8   // kt
	AppWindDir   *uint8   // 1/32 of full circle
	Humidity     *uint8   // %
}

type CompassMessage struct {
	Heading *uint16 // degrees - compass has 1 deg accuracy
	// LATER add more fields
}
type BatteryMessage struct {
	Charge  *uint8   // %
	Voltage *float32 // V
	Current *float32 // A
}

// PositionMessage is used for manually added positions, not from GPS, and does not include course and speed information
type PositionMessage struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}
