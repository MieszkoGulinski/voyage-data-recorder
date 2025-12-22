package formatters

import (
	"fmt"
	"math"
	"time"

	"github.com/racingmars/go3270"
)

func FormatUnixTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	return t.Format("2006-01-02 15:04:05")
}

func FormatLatitude(ptr *float64) string {
	if ptr == nil {
		return "----"
	}
	value := *ptr
	abs := math.Abs(value)
	degrees := int(abs)
	minutes := (abs - float64(degrees)) * 60

	var hemisphere string
	if value >= 0 {
		hemisphere = "N"
	} else {
		hemisphere = "S"
	}

	return fmt.Sprintf("%d %05.2f %s", degrees, minutes, hemisphere)
}

func FormatLongitude(ptr *float64) string {
	if ptr == nil {
		return "----"
	}
	value := *ptr
	abs := math.Abs(value)
	degrees := int(abs)
	minutes := (abs - float64(degrees)) * 60
	
	var hemisphere string
	if value >= 0 {
		hemisphere = "E"
	} else {
		hemisphere = "W"
	}

	return fmt.Sprintf("%d %05.2f %s", degrees, minutes, hemisphere)
}

type Numeric interface {
	~int | ~int64 | ~uint8 | ~uint16 | ~float64
}

func FormatNumber[T Numeric](template string, ptr *T) string {
	if ptr == nil {
		return "----"
	}
	return fmt.Sprintf(template, *ptr)
}

func Format3270Color[T Numeric](ptr *T) go3270.Color {
	if ptr == nil {
		// Data is null - it means that a sensor failed, for this reason we display red
		return go3270.Red
	}

	return go3270.DefaultColor
}

func Format3270ColorWarningDanger(ptr *float64, warningThreshold float64, dangerThreshold float64) go3270.Color {
	if ptr == nil {
		return go3270.Red
	}
	if *ptr > dangerThreshold {
		return go3270.Red
	}
	if *ptr > warningThreshold {
		return go3270.Yellow
	}

	return go3270.DefaultColor
}
