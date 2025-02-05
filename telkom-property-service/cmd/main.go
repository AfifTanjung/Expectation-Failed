package main

import (
	"log"
	"os"
	"telkom-property-service/src/api"
	"telkom-property-service/src/scheduler"
)

func main() {
	//Scheduler untuk API data Validation
	c := scheduler.StartScheduler()
	defer c.Stop()
	r := api.SetupRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
		log.Printf("Using default port: %s", port)
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
