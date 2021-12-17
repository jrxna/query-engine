package routes

import (
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func CreateRouteMappings() {
	Engine = gin.Default()

	// Attach middleware for CORS and Auth here
	// Engine.Use(middleware.Cors())

	// v1 := Engine.Group("/v1")
	{
		// v1.POST("/queries/:id", controllers.GetQueryResult)
	}
}
