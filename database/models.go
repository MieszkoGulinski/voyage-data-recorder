package database

import (
	"gorm.io/gorm"
)

// Add tables structure here

// All columns with sensor data be null if the relevant sensor fails,
// for this reason these columns must be read using pointers.
// Serializing them to JSON will convert them either to number or null.

// Timestamp stores Unix timestamp in seconds.
// GetTimestamp is needed in the queryWithPagination helper.

// Use float64 for all floating point numbers.

type Weather struct {
	Timestamp             int64    `gorm:"column:timestamp;primaryKey" json:"timestamp"`
	AirTemperature        *float64 `gorm:"column:air_temperature_c" json:"airTemperature"`
	WaterTemperature      *float64 `gorm:"column:water_temperature_c" json:"waterTemperature"`
	Pressure              *float64 `gorm:"column:pressure_hpa" json:"pressure"`
	Humidity              *uint8   `gorm:"column:humidity_percent" json:"humidity"`
	ApparentWindDirection *uint8   `gorm:"column:apparent_wind_direction_rhumb" json:"apparentWindDirection"`
	ApparentWindSpeed     *float64 `gorm:"column:apparent_wind_direction_speed_kt" json:"apparentWindSpeed"`
	ApparentWindGustSpeed *float64 `gorm:"column:apparent_wind_direction_gust_speed_kt" json:"apparentWindGustSpeed"`
	TrueWindDirection     *uint8   `gorm:"column:true_wind_direction_rhumb" json:"trueWindDirection"`
	TrueWindSpeed         *float64 `gorm:"column:true_wind_speed_kt" json:"trueWindSpeed"`
	TrueWindGustSpeed     *float64 `gorm:"column:true_wind_gust_speed_kt" json:"trueWindGustSpeed"`
}

func (Weather) TableName() string {
	return "weather"
}

func (r Weather) GetTimestamp() int64 {
	return r.Timestamp
}

type Position struct {
	Timestamp        int64    `gorm:"column:timestamp;primaryKey" json:"timestamp"`
	SourceId         uint8    `gorm:"column:source_id" json:"sourceId"` // 1 = GPS, 2 = manual
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

type BatteryStatus struct {
	Timestamp int64   `gorm:"column:timestamp;primaryKey" json:"timestamp"`
	Charge    uint8   `gorm:"column:charge_percent" json:"charge"`
	Voltage   float64 `gorm:"column:voltage_v" json:"voltage"`
}

func (BatteryStatus) TableName() string {
	return "battery_status"
}

func (r BatteryStatus) GetTimestamp() int64 {
	return r.Timestamp
}

type MotorStatus struct {
	Timestamp    int64 `gorm:"column:timestamp;primaryKey" json:"timestamp"`
	Temperature1 int8  `gorm:"column:temperature1_c" json:"temperature1"`
	Current1     int8  `gorm:"column:current1_a" json:"current1"`
	Temperature2 int8  `gorm:"column:temperature2_c" json:"temperature2"`
	Current2     int8  `gorm:"column:current2_a" json:"current2"`
	Pwm1         int8  `gorm:"column:pwm1" json:"pwm1"`
	Pwm2         int8  `gorm:"column:pwm2" json:"pwm2"`
}

func (MotorStatus) TableName() string {
	return "motor_status"
}

func (r MotorStatus) GetTimestamp() int64 {
	return r.Timestamp
}

type NavtexMessages struct {
	Timestamp int64  `gorm:"column:timestamp;primaryKey" json:"timestamp"`
	Text      string `gorm:"column:text" json:"text"`
}

func (NavtexMessages) TableName() string {
	return "navtex_messages"
}

func (r NavtexMessages) GetTimestamp() int64 {
	return r.Timestamp
}

type TextNotes struct {
	Timestamp int64  `gorm:"column:timestamp;primaryKey" json:"timestamp"`
	Text      string `gorm:"column:text" json:"text"`
}

func (TextNotes) TableName() string {
	return "text_notes"
}

func (r TextNotes) GetTimestamp() int64 {
	return r.Timestamp
}

// More tables may come, e.g. GPS accuracy/status

func RegenerateTables(db *gorm.DB) {
	db.AutoMigrate(&Weather{}, &Position{}, &BatteryStatus{}, &MotorStatus{}, &NavtexMessages{}, &TextNotes{})
}
