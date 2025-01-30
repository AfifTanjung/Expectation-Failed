package main

import (
	"log"
	"os"
	"telkom-property-service/src/api"
)

func main() {
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
