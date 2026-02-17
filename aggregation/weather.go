package aggregation

import (
	"datalogger/channels"
	"datalogger/database"
)

func SummarizeWeather(messages []channels.WeatherMessage, ts int64) database.Weather {
	weather := database.Weather{}
	weather.Timestamp = ts

	var airTempSum float64
	var airTempCnt int
	var pressureSum float64
	var pressureCnt int
	var humiditySum uint8
	var humidityCnt int

	for _, m := range messages {
		if m.Temperature != nil {
			airTempSum += float64(*m.Temperature)
			airTempCnt++
		}
		if m.Pressure != nil {
			pressureSum += float64(*m.Pressure)
			pressureCnt++
		}
		if m.Humidity != nil {
			humiditySum += *m.Humidity
			humidityCnt++
		}

		// TODO other
	}

	var airTempAvg *float64
	if airTempCnt > 0 {
		v := airTempSum / float64(airTempCnt)
		airTempAvg = &v
	}

	var pressureAvg *float64
	if pressureCnt > 0 {
		v := pressureSum / float64(pressureCnt)
		pressureAvg = &v
	}

	var humidityAvg *uint8
	if humidityCnt > 0 {
		v := humiditySum / uint8(humidityCnt)
		humidityAvg = &v
	}

	weather.AirTemperature = airTempAvg
	weather.Pressure = pressureAvg
	weather.Humidity = humidityAvg

	// TODO other

	return weather
}
