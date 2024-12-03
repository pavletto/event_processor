package models

import (
	"time"
)

type Event struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SourceID  string    `gorm:"index" json:"source_id"`
	StartTime float32   `json:"start_time"`
	EndTime   float32   `json:"end_time"`
	Tag       string    `json:"tag"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
