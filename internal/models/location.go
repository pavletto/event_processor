package models

import "time"

type LocationData struct {
	ID            uint      `gorm:"primaryKey"`
	SourceID      string    `gorm:"index" json:"-"`
	RecordingTime time.Time `gorm:"type:timestamp(6)" json:"recording_time"`
	Altitude      *float64  `json:"altitude"`
	Latitude      *float64  `json:"latitude"`
	Longitude     *float64  `json:"longitude"`
	SpeedKMH      *float64  `json:"speed_kmh"`
	SpeedKnots    *float64  `json:"speed_knots"`
	SpeedMPH      *float64  `json:"speed_mph"`
}
