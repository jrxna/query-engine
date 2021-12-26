package main

import (
	"query-engine/config"
	"query-engine/routes"
)

func main() {
	config.Init()

	routes.CreateRouteMappings()
	routes.Engine.Run(":8080")
}
