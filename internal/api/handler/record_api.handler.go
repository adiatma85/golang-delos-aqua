package handler

import (
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/gin-gonic/gin"
)

var recordApiHandler *RecordApiHandler

type RecordApiHandler struct {
	RecordApiRepository repository.RecordApiRepositoryInterface
}

type RecordApiHandlerInterface interface {
	GetAllRecord(c *gin.Context)
}

// Func to get Record Handler instance
func GetRecordApiHandler() RecordApiHandlerInterface {
	if recordApiHandler == nil {
		recordApiHandler = &RecordApiHandler{
			RecordApiRepository: repository.GetRecordApiRepository(),
		}
	}
	return recordApiHandler
}

// HandlerFunc to Get All
func (handler *RecordApiHandler) GetAllRecord(c *gin.Context) {
	recordApiRepo := handler.RecordApiRepository

	records, err := recordApiRepo.GetAll()

	// Internal server error
	if err != nil {
		response := response.BuildFailedResponse("failed to fetch data due to internal server error", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	// Error when no record found
	if len(*records) == 0 {
		response := response.BuildFailedResponse("failed to fetch data due to no data row found", "no record found")
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	// Response
	response := response.BuildSuccessResponse("success to fetch data", records)
	c.JSON(http.StatusOK, response)
}
