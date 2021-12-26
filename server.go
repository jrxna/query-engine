package main

import (
	"query-engine/config"
	"query-engine/database"
	"query-engine/routes"
)

func main() {
	config.Init()
	client, ctx := database.Connect()
	defer client.Disconnect(ctx)

	routes.CreateRouteMappings(client)
	routes.Engine.Run(":8080")
}
