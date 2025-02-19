package main

import (
	"log"
	"os"
	"telkom-property-service/pkg"
	"telkom-property-service/src/api"
	"telkom-property-service/src/scheduler"

	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}
}

func main() {
	// Add more logging
	log.Printf("Starting application...")

	loadEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Using PORT: %s", port)

	// Log all environment variables (be careful with sensitive data)
	log.Printf("Environment variables loaded")

	db, err := pkg.InitDb()
	if err != nil {
		log.Printf("Warning: Database initialization failed: %v", err)
		db = nil
	} else {
		log.Printf("Database connection successful")
	}

	// Setup routes
	r := api.SetupRoutes(db)

	c := scheduler.StartScheduler()
	defer c.Stop()

	address := ":" + port
	log.Printf("Starting server on address: %s", address)

	if err := r.Run(address); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
