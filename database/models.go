package database

// Add tables structure here

// All columns with sensor data be null if the relevant sensor fails,
// for this reason these columns must be read using pointers.
// Serializing them to JSON will convert them either to number or null.

// GetTimestamp is needed in the queryWithPagination helper.

type Weather struct {
	Timestamp             int64    `gorm:"column:timestamp;primaryKey" json:"timestamp"`
	AirTemperature        *float64 `gorm:"column:air_temperature_c" json:"airTemperature"`
	WaterTemperature      *float64 `gorm:"column:water_temperature_c" json:"waterTemperature"`
	Pressure              *float64 `gorm:"column:pressure_hpa" json:"pressure"`
	Sunlight              *uint16  `gorm:"column:sunlight_lux" json:"sunlight"`
	ApparentWindDirection *uint16  `gorm:"column:apparent_wind_direction_deg" json:"apparentWindDirection"`
	ApparentWindSpeeed    *float64 `gorm:"column:apparent_wind_direction_speed_ms" json:"apparentWindSpeed"`
	WindDirection         *uint16  `gorm:"column:true_wind_direction_deg" json:"trueWindDirection"`
	WindSpeed             *float64 `gorm:"column:true_wind_speed_ms" json:"trueWindSpeed"`
	RollAmplitude         *uint8   `gorm:"column:roll_amplitude_deg" json:"rollAmplitude"`
	PitchAmplitude        *uint8   `gorm:"column:pitch_amplitude_deg" json:"pitchAmplitude"`
}

func (Weather) TableName() string {
	return "weather"
}

func (r Weather) GetTimestamp() int64 {
	return r.Timestamp
}

type Position struct {
	Timestamp        int64    `gorm:"column:timestamp;primaryKey" json:"timestamp"`
	SourceId         uint8    `gorm:"column:source_id" json:"sourceId"`
	Latitude         *float64 `gorm:"column:latitude_deg" json:"latitude"`
	Longitude        *float64 `gorm:"column:longitude_deg" json:"lognitude"`
	SpeedOverGround  *float64 `gorm:"column:speed_over_ground_kt" json:"speedOverGround"`
	CourseOverGround *float64 `gorm:"column:course_over_ground_deg" json:"courseOverGround"`
	MagneticBearing  *float64 `gorm:"column:magnetic_bearing_deg" json:"magneticBearing"`
	SpeedOverWater   *float64 `gorm:"column:speed_over_water_kt" json:"speedOverWater"`
}

func (Position) TableName() string {
	return "positions"
}

func (r Position) GetTimestamp() int64 {
	return r.Timestamp
}

type Battery struct {
	Timestamp  int64   `gorm:"column:timestamp;primaryKey" json:"timestamp"`
	Percent    uint8   `gorm:"column:percent" json:"percent"`
	ChangeRate float64 `gorm:"column:change_rate" json:"changeRate"`
}

func (Battery) TableName() string {
	return "battery"
}

func (r Battery) GetTimestamp() int64 {
	return r.Timestamp
}

// More tables will come: GPS accuracy/status, electric motor status etc