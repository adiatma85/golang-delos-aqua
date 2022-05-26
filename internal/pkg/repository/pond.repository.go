package repository

import (
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/pkg/helpers"
)

var pondRepository *PondRepository

type PondRepository struct {
}

type PondRepositoryInterface interface {
	Create(pond models.Pond) (models.Pond, error)
	GetAll() (*[]models.Pond, error)
	GetById(pondId string) (*models.Pond, error)
	Update(pond *models.Pond) error
	Delete(pond *models.Pond) error
}

// Func to return Pond Repository instance
func GetPondRepository() PondRepositoryInterface {
	if pondRepository == nil {
		pondRepository = &PondRepository{}
	}
	return pondRepository
}

// Func to Create Pond
func (repo *PondRepository) Create(pond models.Pond) (models.Pond, error) {
	err := Create(&pond)
	if err != nil {
		return models.Pond{}, err
	}
	return pond, nil
}

// Func to get All Pond without Pagination
func (repo *PondRepository) GetAll() (*[]models.Pond, error) {
	var ponds []models.Pond
	err := Find(&models.Pond{}, &ponds, []string{"Farm"}, "id asc")
	return &ponds, err
}

// Func to Get Pond by Id
func (repo *PondRepository) GetById(pondId string) (*models.Pond, error) {
	var pond models.Pond
	where := models.Pond{}
	where.ID, _ = helpers.ParseUint(pondId)
	_, err := First(&where, &pond, []string{"Farm"})
	if err != nil {
		return nil, err
	}
	return &pond, nil
}

// Func to Update Pond by Model defined in handler
func (repo *PondRepository) Update(pond *models.Pond) error {
	return Save(pond)
}

// Func to Delete Pond by Model defined in handler
func (repo *PondRepository) Delete(pond *models.Pond) error {
	_, err := DeleteByModel(pond)
	if err != nil {
		return err
	}
	return nil
}
