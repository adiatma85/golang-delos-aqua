package fixtures

import "github.com/adiatma85/golang-rest-template-api/internal/pkg/models"

var Ponds []models.Pond = []models.Pond{
	{
		Name:   "Pond 1 in Farm 1",
		FarmId: 1,
	},
	{
		Name:   "Pond 2 in Farm 1",
		FarmId: 1,
	},
	{
		Name:   "Pond 3 in Farm 2",
		FarmId: 2,
	},
}

var WillBePond models.Pond = models.Pond{
	Name:   "Random Pond",
	FarmId: 1,
}
