package models

import (
	"time"
)

type Source struct {
	ID             string              `gorm:"primaryKey" json:"id"`
	DeviceID       string              `gorm:"not null" json:"device_id"`
	EventID        string              `gorm:"not null" json:"event_id"`
	EventType      string              `gorm:"not null" json:"event_type"`
	SerialNumber   string              `json:"serial_number"`
	Message        string              `json:"message"`
	MaxProb        float64             `json:"max_prob"`
	OrganizationID string              `gorm:"not null" json:"organization_id"`
	Hidden         bool                `gorm:"default:false" json:"hidden"`
	StartTime      time.Time           `gorm:"type:timestamp(6)" json:"start_time"`
	EndTime        time.Time           `gorm:"type:timestamp(6)" json:"end_time"`
	IncidentTime   time.Time           `gorm:"type:timestamp(6)" json:"incident_time"`
	CreatedAt      time.Time           `gorm:"type:timestamp(6)" json:"created_at"`
	UpdatedAt      time.Time           `gorm:"type:timestamp(6)" json:"updated_at"`
	ThumbnailKey   string              `json:"thumbnail_key"`
	VideoKey       string              `json:"video_key"`
	Accelerometer  []AccelerometerData `gorm:"foreignKey:SourceID;references:ID" json:"accelerometer"`
	Location       []LocationData      `gorm:"foreignKey:SourceID;references:ID" json:"location"`
	Risk           []RiskData          `gorm:"foreignKey:SourceID;references:ID" json:"risk"`
	Events         []Event             `gorm:"foreignKey:SourceID;references:ID" json:"events"`
}
