package handler

import (
	"errors"
	"fmt"
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
		farmHandler = &FarmHandler{
			FarmRepository: repository.GetFarmRepository(),
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

	farmRepo := handler.FarmRepository
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
	farmRepo := handler.FarmRepository

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
	farmRepo := handler.FarmRepository

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

	farmRepo := handler.FarmRepository

	// Check whether "ID" is specified in payload or not
	if updateFarmRequest.ID == 0 {
		// Not specified, so create it
		farmModel := &models.Farm{}

		smapping.FillStruct(farmModel, smapping.MapFields(&updateFarmRequest))

		// Check whether there is error when creating
		// Need new "small functional" so it does not duplicate
		if newFarm, err := farmRepo.Create(*farmModel); err != nil {
			response := response.BuildFailedResponse("failed to add new farm due to internal server error", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		} else {
			response := response.BuildSuccessResponse("success add new farm instance to database", newFarm)
			c.JSON(http.StatusOK, response)
		}
	} else {
		// Specified so update it
		existedFarm, err := farmRepo.GetById(fmt.Sprint(updateFarmRequest.ID))

		// If specified resource does not exist
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			response := response.BuildFailedResponse("failed to fetch data due to no record found with specified id", err.Error())
			c.AbortWithStatusJSON(http.StatusNotFound, response)
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
}

// HandlerFunc to Delete
func (handler *FarmHandler) Delete(c *gin.Context) {
	farmRepo := handler.FarmRepository
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
