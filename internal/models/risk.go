package models

import "time"

type RiskData struct {
	ID            uint      `gorm:"primaryKey"`
	EventID       string    `gorm:"index" json:"-"`
	RecordingTime time.Time `gorm:"type:timestamp(6)" json:"recording_time"`
	Prob          float64   `json:"prob"`
}
