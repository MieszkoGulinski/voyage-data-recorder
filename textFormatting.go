package main

import (
	"fmt"
	"math"
	"time"
)

func formatUnixTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0).UTC()
	return t.Format("2006-01-02 15:04:05")
}

func formatLatitude(ptr *float64) string {
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

func formatLongitude(ptr *float64) string {
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

func formatNumber[T Numeric](template string, ptr *T) string {
	if ptr == nil {
		return "----"
	}
	return fmt.Sprintf(template, *ptr)
}
