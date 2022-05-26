package dto

import (
	"gorm.io/gorm"
)

type subFarm struct {
	gorm.Model
	Name string `json:"name"`
}

type PondResponseDto struct {
	gorm.Model
	Name string  `json:"name"`
	Farm subFarm `json:"farm"`
}
