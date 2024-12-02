package models

import "time"

type AccelerometerData struct {
	ID            uint      `gorm:"primaryKey"`
	EventID       string    `gorm:"index" json:"-"`
	RecordingTime time.Time `gorm:"type:timestamp(6)" json:"recording_time"`
	X             float64   `json:"x"`
	Y             float64   `json:"y"`
	Z             float64   `json:"z"`
}
