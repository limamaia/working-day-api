package main

import (
	"log"
	"working-day-api/config"
	_ "working-day-api/docs"
	"working-day-api/internal/container"
	"working-day-api/router"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	container := container.NewContainer(config)

	router.LoadRoutes(config, container)
}
