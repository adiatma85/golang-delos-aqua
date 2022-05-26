package repository

import "github.com/adiatma85/golang-rest-template-api/internal/pkg/models"

var recordApiRepository RecordApiRepositoryInterface

type RecordApiRepository struct {
}

type RecordApiRepositoryInterface interface {
	Create(record models.RecordApi)
	GetAll() (*[]models.RecordApi, error)
	GetByModel(where models.RecordApi) (*models.RecordApi, error)
	UpdateCount(record *models.RecordApi) error
}

// Func to return instance of Record Api Interface
func GetRecordApiRepository() RecordApiRepositoryInterface {
	if recordApiRepository == nil {
		recordApiRepository = &RecordApiRepository{}
	}
	return recordApiRepository
}

// Func to Create Api Record
func (repo *RecordApiRepository) Create(record models.RecordApi) {
	Create(&record)
}

// Func to Get All Record
func (repo *RecordApiRepository) GetAll() (*[]models.RecordApi, error) {
	var records []models.RecordApi
	err := Find(&models.RecordApi{}, &records, []string{}, "request_path asc")
	return &records, err
}

// Func to Get from Model
func (repo *RecordApiRepository) GetByModel(where models.RecordApi) (*models.RecordApi, error) {
	var recordApi models.RecordApi
	_, err := First(&where, &recordApi, []string{})
	if err != nil {
		return nil, err
	}
	return &recordApi, err
}

// Func to Update the Count
func (repo *RecordApiRepository) UpdateCount(record *models.RecordApi) error {
	record.Count++
	return Save(record)
}
