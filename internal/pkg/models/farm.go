package models

import "gorm.io/gorm"

// Struct for Farm Models
type Farm struct {
	gorm.Model
	Name  string `gorm:"type:varchar(100)" json:"name"`
	Ponds []Pond `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
