package routes

import (
	"query-engine/controller"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var Engine *gin.Engine

func CreateRouteMappings(client *mongo.Client) {
	Engine = gin.Default()

	queryController := new(controller.QueryController)
	queryController.Database = client

	Engine.POST("/v1/query-result", queryController.GetQueryResult)
}
