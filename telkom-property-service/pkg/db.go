package pkg

import (
	"fmt"
	"log"
	"os"
	"telkom-property-service/src/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || password == "" || name == "" || port == "" {
		return nil, fmt.Errorf("database configuration is missing or incomplete")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Starting database migrations...")
	err = db.AutoMigrate(
		&model.Alat{},
		&model.Raw{},
		&model.Valid{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform database migrations: %w", err)
	}

	log.Println("Database connection established and migrations completed successfully.")
	return db, nil
}
