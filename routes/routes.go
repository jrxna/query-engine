package routes

import (
	"query-engine/controller"

	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

func CreateRouteMappings() {
	Engine = gin.Default()

	// Attach middleware for CORS and Auth here
	// Engine.Use(middleware.Cors())

	Engine.POST("/v1/query-result", controller.GetQueryResult)
}
