package v1

import (
	"fmt"
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/internal/api/handler"
	"github.com/adiatma85/golang-rest-template-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// V1 Router
func Setup() *gin.Engine {
	app := gin.New()

	// Middlewares
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s \" \" %s\" \" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	app.Use(gin.Recovery())
	app.Use(middleware.CORS())
	app.Use(middleware.RecordApi())
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	// Routes for v1
	v1Route := app.Group("/api/v1")
	{
		v1Route.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "Welcome")
		})
	}

	// FarmGroup
	farmGroup := v1Route.Group("farm")
	farmHandler := handler.GetFarmHandler()
	{
		farmGroup.GET("", farmHandler.GetAllFarm)
		farmGroup.GET(":farmId", farmHandler.GetById)
		farmGroup.POST("", farmHandler.CreateFarm)
		// Non-standar PUT route according to requirement
		farmGroup.PUT("", farmHandler.Update)
		farmGroup.DELETE(":farmId", farmHandler.Delete)
	}

	// PondGroup
	pondGroup := v1Route.Group("pond")
	pondHandler := handler.GetPondHandler()
	{
		pondGroup.GET("", pondHandler.GetAllPond)
		pondGroup.GET(":pondId", pondHandler.GetById)
		pondGroup.POST("", pondHandler.CreatePond)
		// Non-standar PUT route according to requirement
		pondGroup.PUT("", pondHandler.Update)
		pondGroup.DELETE(":pondId", pondHandler.Delete)
	}

	return app
}
