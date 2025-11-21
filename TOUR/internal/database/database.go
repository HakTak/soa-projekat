package database

import (
	"fmt"
	"log"
	"os"
	"tour/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Belgrade",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "super"),
		getEnv("DB_NAME", "tours"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(" Failed to connect to database:", err)
	}

	// Auto migrate tabele
	if err := db.AutoMigrate(&model.Tour{}); err != nil {
		log.Fatal(" AutoMigrate failed:", err)
	}

	DB = db
	log.Println(" Connected to PostgreSQL & migrated Tour model")
	return db
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
