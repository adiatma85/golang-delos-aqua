package handler

import (
	"errors"
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/validator"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
	"gorm.io/gorm"
)

var farmHandler *FarmHandler

type FarmHandler struct {
	FarmRepository repository.FarmRepositoryInterface
}

type FarmHandlerInterface interface {
	CreateFarm(c *gin.Context)
	GetAllFarm(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// Func to get Farm Handler instance
func GetFarmHandler() FarmHandlerInterface {
	if farmHandler == nil {
		farmRepo := repository.GetFarmRepository()
		farmHandler = &FarmHandler{
			farmRepo,
		}
	}
	return farmHandler
}

// HandlerFunc to Create Farm (POST)
func (handler *FarmHandler) CreateFarm(c *gin.Context) {
	var createFarmRequest validator.CreateFarmRequest
	err := c.ShouldBind(&createFarmRequest)

	// Bad Request
	if err != nil {
		response := response.BuildFailedResponse("failed to add new farm due to bad request", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	farmRepo := repository.GetFarmRepository()
	farmModel := &models.Farm{}

	// smapping the struct
	smapping.FillStruct(farmModel, smapping.MapFields(&createFarmRequest))

	if newFarm, err := farmRepo.Create(*farmModel); err != nil {
		response := response.BuildFailedResponse("failed to add new farm due to internal server error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	} else {
		response := response.BuildSuccessResponse("success add new farm instance to database", newFarm)
		c.JSON(http.StatusOK, response)
		return
	}
}

// HandlerFunc to Get All
func (handler *FarmHandler) GetAllFarm(c *gin.Context) {
	farmRepo := repository.GetFarmRepository()

	farms, err := farmRepo.GetAll()

	// Internal server error
	if err != nil {
		response := response.BuildFailedResponse("failed to fetch data due to internal server error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	// Error when no record found
	if len(*farms) == 0 {
		response := response.BuildFailedResponse("failed to fetch data due to no data row found", "no record found")
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	// Response
	response := response.BuildSuccessResponse("success to fetch data", farms)
	c.JSON(http.StatusOK, response)
}

// HandlerFunc to Get By Id
func (handler *FarmHandler) GetById(c *gin.Context) {
	farmRepo := repository.GetFarmRepository()

	farm, err := farmRepo.GetById(c.Param("farmId"))

	if err != nil {
		var failedResponse response.Response
		switch {
		// Case error not found
		case errors.Is(err, gorm.ErrRecordNotFound):
			failedResponse = response.BuildFailedResponse("failed to fetch data due to no record found", err.Error())
			c.AbortWithStatusJSON(http.StatusNotFound, failedResponse)
		// Case error on internal server error
		default:
			failedResponse = response.BuildFailedResponse("failed to fetch data", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, failedResponse)
		}
		return
	}

	response := response.BuildSuccessResponse("success to fetch data", farm)
	c.JSON(http.StatusOK, response)
}

// Handlerfunc to Update
func (handler *FarmHandler) Update(c *gin.Context) {
	var updateFarmRequest validator.UpdateFarmRequest
	err := c.ShouldBind(&updateFarmRequest)

	if err != nil {
		response := response.BuildFailedResponse("failed to update new farm due to bad request", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	farmRepo := repository.GetFarmRepository()

	// check whether the farm is exist or not
	existedFarm, err := farmRepo.GetById(c.Param("farmId"))

	// If farm does not exist, create it
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// Create
		// Need more improvise to make it
		farmModel := &models.Farm{}
		smapping.FillStruct(farmModel, smapping.MapFields(&updateFarmRequest))
		if newFarm, err := farmRepo.Create(*farmModel); err != nil {
			response := response.BuildFailedResponse("failed to add new farm due to internal server error", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		} else {
			response := response.BuildSuccessResponse("success add new farm instance to database", newFarm)
			c.JSON(http.StatusOK, response)
		}
		return
	}

	smapping.FillStruct(existedFarm, smapping.MapFields(&updateFarmRequest))

	err = farmRepo.Update(existedFarm)

	if err != nil {
		response := response.BuildFailedResponse("failed to update a farm", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// HandlerFunc to Delete
func (handler *FarmHandler) Delete(c *gin.Context) {
	farmRepo := repository.GetFarmRepository()
	existedFarm, err := farmRepo.GetById(c.Param("farmId"))

	// If error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		failedResponse := response.BuildFailedResponse("failed to fetch data due to no record found", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, failedResponse)
		return
	}

	err = farmRepo.Delete(existedFarm)

	if err != nil {
		response := response.BuildFailedResponse("failed to delete a farm", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
