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

var pondHandler *PondHandler

type PondHandler struct {
	PondRepository repository.PondRepositoryInterface
}

type PondHandlerInterface interface {
	CreatePond(c *gin.Context)
	GetAllPond(c *gin.Context)
	GetById(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// Func to get Pond Handler instance
func GetPondHandler() PondHandlerInterface {
	if pondHandler == nil {
		pondHandler = &PondHandler{
			PondRepository: repository.GetPondRepository(),
		}
	}
	return pondHandler
}

// HandlerFunc to Create Pond (POST)
func (handler *PondHandler) CreatePond(c *gin.Context) {
	var createPondRequest validator.CreatePondRequest
	err := c.ShouldBind(&createPondRequest)

	// Bad Request
	if err != nil {
		response := response.BuildFailedResponse("failed to add new pond due to bad request", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	pondRepo := handler.PondRepository
	pondModel := &models.Pond{}

	// smapping the struct
	smapping.FillStruct(pondModel, smapping.MapFields(&createPondRequest))

	if newPond, err := pondRepo.Create(*pondModel); err != nil {
		response := response.BuildFailedResponse("failed to add new farm due to internal server error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	} else {
		response := response.BuildSuccessResponse("success add new farm instance to database", newPond)
		c.JSON(http.StatusOK, response)
		return
	}
}

// HandlerFunc to Get All
func (handler *PondHandler) GetAllPond(c *gin.Context) {
	pondRepo := handler.PondRepository

	ponds, err := pondRepo.GetAll()

	// Internal server error
	if err != nil {
		response := response.BuildFailedResponse("failed to fetch data due to internal server error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	// Error when no record found
	if len(*ponds) == 0 {
		response := response.BuildFailedResponse("failed to fetch data due to no data row found", "no record found")
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	// Response
	response := response.BuildSuccessResponse("success to fetch data", ponds)
	c.JSON(http.StatusOK, response)
}

// HandlerFunc to Get By Id
func (handler *PondHandler) GetById(c *gin.Context) {
	pondRepo := handler.PondRepository

	pond, err := pondRepo.GetById(c.Param("pondId"))

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

	response := response.BuildSuccessResponse("success to fetch data", pond)
	c.JSON(http.StatusOK, response)
}

// HandlerFunc to Update
func (handler *PondHandler) Update(c *gin.Context) {
	var updatePondRequest validator.UpdatePondRequest
	err := c.ShouldBind(&updatePondRequest)

	if err != nil {
		response := response.BuildFailedResponse("failed to update new farm due to bad request", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	pondRepo := handler.PondRepository

	if updatePondRequest.ID == 0 {
		pondModel := &models.Pond{}
		smapping.FillStruct(pondModel, smapping.MapFields(&updatePondRequest))

		// Check whether there is error when creating
		// Need new "small functional" so it does not duplicate
		if newPond, err := pondRepo.Create(*pondModel); err != nil {
			response := response.BuildFailedResponse("failed to add new farm due to internal server error", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		} else {
			response := response.BuildSuccessResponse("success add new farm instance to database", newPond)
			c.JSON(http.StatusOK, response)
			return
		}
	} else {
		// Specified, so update it
		existedPond, err := pondRepo.GetById(fmt.Sprint(updatePondRequest.ID))

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			response := response.BuildFailedResponse("failed to fetch data due to no record found with specified id", err.Error())
			c.AbortWithStatusJSON(http.StatusNotFound, response)
			return
		}

		smapping.FillStruct(existedPond, smapping.MapFields(&updatePondRequest))
		err = pondRepo.Update(existedPond)

		if err != nil {
			response := response.BuildFailedResponse("failed to update a farm", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// HandlerFunc to Delete
func (handler *PondHandler) Delete(c *gin.Context) {
	pondRepo := handler.PondRepository
	existedPond, err := pondRepo.GetById(c.Param("pondId"))

	// If error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		failedResponse := response.BuildFailedResponse("failed to fetch data due to no record found", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, failedResponse)
		return
	}

	err = pondRepo.Delete(existedPond)

	if err != nil {
		response := response.BuildFailedResponse("failed to delete a pond", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
