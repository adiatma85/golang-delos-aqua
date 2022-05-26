package middleware

import (
	"errors"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/pkg/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Middleware to record api to database
// Reference --> https://github.com/sbecker/gin-api-demo/blob/master/middleware/json_logger.go
func RecordApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Do the next first to get the status
		c.Next()

		// Search if it is exist or not
		recordRepo := repository.GetRecordApiRepository()
		whereRecord := models.RecordApi{
			RequestPath: c.Request.RequestURI,
			UserAgent:   helpers.GetClientIP(c),
			Status:      c.Writer.Status(),
			Referer:     c.Request.Referer(),
		}

		existedRecord, err := recordRepo.GetByModel(whereRecord)

		// If not exist, make new entry
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			whereRecord.Count = 1
			recordRepo.Create(whereRecord)
			return
		}

		// If exist, update the count of the access trafic
		recordRepo.UpdateCount(existedRecord)
	}
}
