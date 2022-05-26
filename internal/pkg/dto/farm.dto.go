package dto

import (
	"gorm.io/gorm"
)

type subPond struct {
	gorm.Model
	Name string `json:"name"`
}

type FarmResponseDto struct {
	gorm.Model
	Name  string    `json:"name"`
	Ponds []subPond `json:"ponds"`
}
