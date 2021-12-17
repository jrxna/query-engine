package main

import (
	"query-engine/routes"
)

func main() {
	routes.CreateRouteMappings()
	routes.Engine.Run(":8080")
}
