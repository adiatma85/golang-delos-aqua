package helpers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetClientIP gets the correct IP for the end client instead of the proxy
// Reference --> https://github.com/sbecker/gin-api-demo/blob/07f9a9242f743fc51ae4a046ee58e12627bad571/util/log.go
func GetClientIP(c *gin.Context) string {
	// first check the X-Forwarded-For header
	requester := c.Request.Header.Get("X-Forwarded-For")
	// if empty, check the Real-IP header
	if len(requester) == 0 {
		requester = c.Request.Header.Get("X-Real-IP")
	}
	// if the requester is still empty, use the hard-coded address from the socket
	if len(requester) == 0 {
		requester = c.Request.RemoteAddr
	}

	// if requester is a comma delimited list, take the first one
	// (this happens when proxied via elastic load balancer then again through nginx)
	if strings.Contains(requester, ",") {
		requester = strings.Split(requester, ",")[0]
	}

	return requester
}

// Helper function to change id from string param to uint
// (according to gorm model)
func ParseUint(id string) (uint, error) {
	parsedUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsedUint), nil
}
