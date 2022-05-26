package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/adiatma85/golang-rest-template-api/internal/api/router/v1"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/db"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/adiatma85/golang-rest-template-api/test"
	"github.com/adiatma85/golang-rest-template-api/test/fixtures"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type FarmHandlerSuite struct {
	suite.Suite
	Router *gin.Engine
}

func TestFarmHandler(t *testing.T) {
	suite.Run(t, new(FarmHandlerSuite))
	defer test.TearDownHelper()
}

// Function to initialize the test suite
func (suite *FarmHandlerSuite) SetupSuite() {
	// Initialize Configuration
	test.SetupInitialize("../../.env")
	db.SetupTestingDb(test.Host, test.Username, test.Password, test.Port, test.Database)

	// Initialize Router for testing
	suite.Router = v1.Setup()
}

// Function to Create new Farm
func (suite *FarmHandlerSuite) TestCreateFarm_Positive() {
	a := suite.Assert()

	newBody := models.Farm{
		Name: "new one",
	}

	requestBody, err := json.Marshal(newBody)
	if err != nil {
		a.Error(err)
	}

	req, w := createFarm(suite.Router, bytes.NewBuffer(requestBody))
	a.Equal(http.MethodPost, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")
}

// Function to Get All but return success because there is exist record
func (suite *FarmHandlerSuite) TestGetAllFarm_Positive() {
	farm, err := insertFarm()
	a := suite.Assert()

	a.NotNil(farm, "fail to insert resource")
	a.NoError(err, "fail to insert resource")

	req, w := getAllFarmRequest(suite.Router)
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request code error")
}

// Function to Get By Id but return success because there is exist record
func (suite *FarmHandlerSuite) TestGetById_Positive() {
	farm, err := insertFarm()
	a := suite.Assert()

	a.NotNil(farm, "fail to insert resource")
	a.NoError(err, "fail to insert resource")

	req, w := getFarmByIdRequest(suite.Router, 1)
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request code error")

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		a.Error(err)
	}

	actual := response.Response{}
	if err := json.Unmarshal(body, &actual); err != nil {
		a.Error(err)
	}

	a.Nil(actual.Errors, "error at response")
	a.Equal("success to fetch data", actual.Message, "response message is different than supposed to be")
}

// Function to Get By Id but return not found because there is no record
func (suite *FarmHandlerSuite) TestGetById_Negative() {
	req, w := getFarmByIdRequest(suite.Router, 1000)
	a := suite.Assert()
	a.Equal(http.MethodGet, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request code error")
}

// Functon to Update an Existing Resource
func (suite *FarmHandlerSuite) TestUpdate_Existing() {
	farm, err := insertFarm()
	a := suite.Assert()

	a.NotNil(farm, "fail to insert resource")
	a.NoError(err, "fail to insert resource")

	updateBody := models.Farm{
		Model: gorm.Model{
			ID: farm.ID,
		},
		Name: "new one",
	}

	requestBody, err := json.Marshal(updateBody)
	if err != nil {
		a.Error(err)
	}

	req, w := updateFarm(suite.Router, bytes.NewBuffer(requestBody))
	a.Equal(http.MethodPut, req.Method, "HTTP request method error")
	a.Equal(http.StatusNoContent, w.Code, "HTTP request status code error")
}

// Function to Update Non-Existing Resource
// Therefore, it will create new Resource
func (suite *FarmHandlerSuite) TestUpdate_NonExisting() {
	a := suite.Assert()

	updateBody := models.Farm{
		Name: "new one",
	}

	requestBody, err := json.Marshal(updateBody)
	if err != nil {
		a.Error(err)
	}

	req, w := updateFarm(suite.Router, bytes.NewBuffer(requestBody))
	a.Equal(http.MethodPut, req.Method, "HTTP request method error")
	a.Equal(http.StatusOK, w.Code, "HTTP request status code error")
}

// Function to Delete By Id but return Positive
func (suite *FarmHandlerSuite) TestDeleteById_Positive() {
	farm, err := insertFarm()
	a := suite.Assert()

	a.NotNil(farm, "fail to insert resource")
	a.NoError(err, "fail to insert resource")

	req, w := deleteFarmByIdRequest(suite.Router, farm.ID)
	a.Equal(http.MethodDelete, req.Method, "HTTP request method error")
	a.Equal(http.StatusNoContent, w.Code, "HTTP request code error")
}

// Function to Delete By Id but return negative
func (suite *FarmHandlerSuite) TestDeleteById_Negative() {
	a := suite.Assert()
	req, w := deleteFarmByIdRequest(suite.Router, 1000)
	a.Equal(http.MethodDelete, req.Method, "HTTP request method error")
	a.Equal(http.StatusNotFound, w.Code, "HTTP request code error")
}

// Helper function to createFarm
func createFarm(r *gin.Engine, body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(http.MethodPost, "/api/v1/farm", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

// Helper function getAllFarm
func getAllFarmRequest(r *gin.Engine) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/farm", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

// Helper function getById
func getFarmByIdRequest(r *gin.Engine, farmId uint) (*http.Request, *httptest.ResponseRecorder) {
	url := fmt.Sprintf("/api/v1/farm/%d", farmId)
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
func updateFarm(r *gin.Engine, body *bytes.Buffer) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(http.MethodPut, "/api/v1/farm", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

// Helper function deleteById
func deleteFarmByIdRequest(r *gin.Engine, farmId uint) (*http.Request, *httptest.ResponseRecorder) {
	url := fmt.Sprintf("/api/v1/far/%d", farmId)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

// Helper function insertFarm
func insertFarm() (models.Farm, error) {
	farmRepo := repository.GetFarmRepository()
	return farmRepo.Create(fixtures.WillBeFarm)
}
