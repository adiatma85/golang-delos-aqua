package models

import "gorm.io/gorm"

// Struct for Pond Models
type Pond struct {
	gorm.Model
	Name string `gorm:"type:varchar(100)" json:"name"`
	Farm Farm   `gorm:"foreignkey:FarmId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
}
