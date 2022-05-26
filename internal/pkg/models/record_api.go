package models

import (
	"time"

	"gorm.io/gorm"
)

// Record Api Model for record the traffic
type RecordApi struct {
	ID          uint           `gorm:"primarykey"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	RequestPath string         `json:"request_path"`
	UserAgent   string         `json:"user_agent"`
	Status      int            `json:"status"`
	Referer     string         `json:"referer"`
	Count       uint           `json:"count"`
}
