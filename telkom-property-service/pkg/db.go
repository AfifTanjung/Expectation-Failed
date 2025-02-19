package pkg

import (
	"fmt"
	"log"
	"os"
	"telkom-property-service/src/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb() (*gorm.DB, error) {
	var dsn string

	// Ambil environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")
	isLocal := os.Getenv("LOCAL_ENV")

	// Validasi konfigurasi database
	if dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatalf("Database configuration is missing! Please set DB_USER, DB_PASSWORD, DB_NAME.")
	}

	// Konfigurasi koneksi database
	if instanceConnectionName != "" && !isLocalSet(isLocal) {
		// Koneksi melalui Unix Socket (Cloud SQL)
		dsn = fmt.Sprintf(
			"host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable",
			instanceConnectionName, dbUser, dbPassword, dbName,
		)
		log.Println("Using Unix Socket connection for Cloud SQL")
	} else if dbHost != "" {
		// Koneksi via TCP untuk Cloud Run
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
			dbHost, dbUser, dbPassword, dbName,
		)
		log.Println("Using standard TCP connection")
	} else if isLocalSet(isLocal) {
		// Koneksi ke database lokal jika ENV LOCAL_ENV diset
		dsn = fmt.Sprintf(
			"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
			dbUser, dbPassword, dbName,
		)
		log.Println("Using Localhost connection")
	} else {
		log.Fatalf("No valid database connection method found (set DB_HOST or INSTANCE_CONNECTION_NAME)")
	}

	// Coba koneksi database dengan retry
	var (
		db  *gorm.DB
		err error
	)
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			log.Println("Database connection established successfully")
			break
		}
		log.Printf("Database connection failed (attempt %d/%d): %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			time.Sleep(2 * time.Second) // Tunggu sebelum retry
		}
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	// Lakukan migrasi
	log.Println("Running database migrations...")
	err = db.AutoMigrate(
		&model.Alat{},
		&model.Raw{},
		&model.Valid{},
		&model.Maintenance{},
		&model.Admin{},
	)
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Println("Database connected and migrations completed successfully.")
	return db, nil
}

// Fungsi untuk mengecek jika environment LOCAL_ENV di set untuk pengembangan lokal
func isLocalSet(isLocal string) bool {
	return isLocal == "true"
}
