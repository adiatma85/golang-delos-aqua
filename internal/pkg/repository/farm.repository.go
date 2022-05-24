package repository

import (
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/pkg/helpers"
)

var farmRepository *FarmRepository

type FarmRepositoryInterface interface {
	Create(farm models.Farm) (models.Farm, error)
	GetAll() (*[]models.Farm, error)
	GetById(farmId string) (*models.Farm, error)
	Update(farm *models.Farm) error
	Delete(farm *models.Farm) error
}

type FarmRepository struct {
}

// Func to return instance of Farm Repository
func GetFarmRepository() FarmRepositoryInterface {
	if farmRepository == nil {
		farmRepository = &FarmRepository{}
	}
	return farmRepository
}

// Func to Create Farm
func (repo *FarmRepository) Create(farm models.Farm) (models.Farm, error) {
	err := Create(&farm)
	if err != nil {
		return models.Farm{}, err
	}
	return farm, nil
}

// Func to get All Farm without Pagination
func (repo *FarmRepository) GetAll() (*[]models.Farm, error) {
	var farms []models.Farm
	err := Find(&models.Farm{}, &farms, []string{}, "id asc")
	return &farms, err
}

// Func to get By Id
func (repo *FarmRepository) GetById(farmId string) (*models.Farm, error) {
	var farm models.Farm
	where := models.Farm{}
	where.ID, _ = helpers.ParseUint(farmId)
	_, err := First(&where, &farm, []string{})
	if err != nil {
		return nil, err
	}
	return &farm, nil
}

// Func to update farm according to model defined
func (repo *FarmRepository) Update(farm *models.Farm) error {
	return Save(farm)
}

// Func to delete farm according to model defined
func (repo *FarmRepository) Delete(farm *models.Farm) error {
	_, err := DeleteByModel(farm)
	if err != nil {
		return err
	}
	return nil
}
