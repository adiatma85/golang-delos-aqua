package repository

import "github.com/adiatma85/golang-rest-template-api/internal/pkg/models"

var farms []models.Farm = []models.Farm{
	{
		Name:  "Farm 1",
		Ponds: []models.Pond{},
	},
	{
		Name:  "Random Farm",
		Ponds: []models.Pond{},
	},
	{
		Name:  "Farm 2",
		Ponds: []models.Pond{},
	},
}

var willBeFarm models.Farm = models.Farm{
	Name:  "Random Farm",
	Ponds: []models.Pond{},
}
