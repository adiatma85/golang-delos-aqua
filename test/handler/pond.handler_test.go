package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/adiatma85/golang-rest-template-api/internal/api/router/v1"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/db"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/test"
	"github.com/adiatma85/golang-rest-template-api/test/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PondHandlerSuite struct {
	suite.Suite
	Router *gin.Engine
}

func TestPondHandler(t *testing.T) {
	suite.Run(t, new(PondHandlerSuite))
}

// Function to initialize the test suite
func (suite *PondHandlerSuite) SetupSuite() {
	// Initialize Configuration
	test.SetupInitialize("../../.env")
	db.SetupTestingDb(test.Host, test.Username, test.Password, test.Port, test.Database)

	// Initialize Router for testing
	suite.Router = v1.Setup()
}

// Function to Create new pond
func (suite *PondHandlerSuite) TestCreatePond_Positive() {
	farm, _ := insertFarm()
	a := suite.Assert()

	newBody := models.Pond{
		Name:   "new one",
		FarmId: farm.ID,
	}

	requestBody, err := json.Marshal(newBody)
	if err != nil {
		a.Error(err)
	}

	req, w := createPond(suite.Router, bytes.NewBuffer(requestBody))
	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")
}

// Function to Get All but return success because there is exist record
func (suite *PondHandlerSuite) TestGetAllPond_Positive() {
	pond, err := insertPond()
	a := suite.Assert()

	a.NotNil(pond, "fail to insert resource")
	a.NoError(err, "fail to insert resource")

	req, w := getAllPondRequest(suite.Router)
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request code error")
}

// Function to Get All but return success because there is exist record
func (suite *PondHandlerSuite) TestGetAllPond_Negative() {
	a := suite.Assert()
	req, w := getAllPondRequest(suite.Router)
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request code error")
}

// Function to Get By Id but and return success because there is exist record
func (suite *PondHandlerSuite) TestGetById_Positive() {
	pond, err := insertPond()
	a := suite.Assert()

	a.NotNil(pond, "fail to insert resource")
	a.NoError(err, "fail to insert resource")

	req, w := getPondmByIdRequest(suite.Router, 1)
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request code error")
}

// Function to Get By Id but but return not found
func (suite *PondHandlerSuite) TestGetById_Negative() {
	a := suite.Assert()
	req, w := getPondmByIdRequest(suite.Router, 1000)
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request code error")
}

// Functon to Update an Existing Resource
func (suite *PondHandlerSuite) TestUpdate_Existing() {
	farm, _ := insertFarm()
	pond, err := insertPond()
	a := suite.Assert()

	a.NotNil(pond, "fail to insert resource")
	a.NoError(err, "fail to insert resource")

	updateBody := models.Pond{
		Model: gorm.Model{
			ID: pond.ID,
		},
		Name:   "new one edited",
		FarmId: farm.ID,
	}

	requestBody, err := json.Marshal(updateBody)
	if err != nil {
		a.Error(err)
	}

	req, w := updatePond(suite.Router, bytes.NewBuffer(requestBody))
	a.Equal(http.MethodPut, req.Method, "HTTP request method error")
	a.Equal(http.StatusNoContent, w.Code, "HTTP request status code error")
}

// Function to delete by id and return success
func (suite *PondHandlerSuite) TestDeleteById_Positive() {
	pond, err := insertPond()
	a := suite.Assert()

	a.NotNil(pond, "fail to insert resource")
	a.NoError(err, "fail to insert resource")

	req, w := deletePondByIdRequest(suite.Router, pond.ID)
	a.Equal(http.MethodDelete, req.Method, "HTTP request method error")
	a.Equal(http.StatusNoContent, w.Code, "HTTP request code error")
}

// Function to Update Non-Existing Resource
// Therefore, it will create new Resource
func (suite *PondHandlerSuite) TestUpdate_NonExisting() {
	farm, _ := insertFarm()
	a := suite.Assert()

	newBody := models.Pond{
		Name:   "new one",
		FarmId: farm.ID,
	}

	requestBody, err := json.Marshal(newBody)
	if err != nil {
		a.Error(err)
	}

	req, w := updatePond(suite.Router, bytes.NewBuffer(requestBody))
	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")
}

// Function to Delete By Id but return negative
func (suite *PondHandlerSuite) TestDeleteById_Negative() {
	a := suite.Assert()
	req, w := deletePondByIdRequest(suite.Router, 1000)
	a.Equal(http.MethodDelete, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request code error")
}

// Helper function createPond
func createPond(r *gin.Engine, body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(http.MethodPost, "/api/v1/pond", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

// Helper function getAllFarm
func getAllPondRequest(r *gin.Engine) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/pond", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

// Helper function getById
func getPondmByIdRequest(r *gin.Engine, pondId uint) (*http.Request, *httptest.ResponseRecorder) {
	url := fmt.Sprintf("/api/v1/pond/%d", pondId)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

// Helper function update
func updatePond(r *gin.Engine, body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(http.MethodPut, "/api/v1/pond", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

// Helper function deleteById
func deletePondByIdRequest(r *gin.Engine, pondId uint) (*http.Request, *httptest.ResponseRecorder) {
	url := fmt.Sprintf("/api/v1/pond/%d", pondId)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

// Helper function insertPond
func insertPond() (models.Pond, error) {
	pondRepo := repository.GetPondRepository()
	return pondRepo.Create(fixtures.WillBePond)
}
