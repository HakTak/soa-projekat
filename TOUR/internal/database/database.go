package database

import (
	"fmt"
	"log"
	"os"
	"time"
	"tour/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "super"),
		getEnv("DB_NAME", "tours"),
		getEnv("DB_PORT", "5432"),
	)

	var db *gorm.DB
	var err error
	maxRetries := 5

	// Petlja za ponovno povezivanje (Retry Logic)
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err == nil {
			// GORM nekad ne prijavljuje gresku odmah, pa radimo Ping da proverimo "zivu" vezu
			sqlDB, _ := db.DB()
			if errPing := sqlDB.Ping(); errPing == nil {
				log.Println("✅ Successfully connected to database!")
				break // Uspesno povezivanje, izlazimo iz petlje
			}
		}

		log.Printf("⏳ Failed to connect to database (attempt %d/%d). Retrying in 2 seconds... Error: %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	// Ako i posle 5 pokusaja nismo uspeli, tek onda rusimo aplikaciju
	if err != nil {
		log.Fatal("❌ Could not connect to database after multiple retries:", err)
	}

	// Auto migrate tabele
	log.Println("🛠 Running migrations...")
	if err := db.AutoMigrate(&model.Tour{}); err != nil {
		log.Fatal("❌ AutoMigrate Tour failed:", err)
	}

	if err := db.AutoMigrate(&model.Comment{}); err != nil {
		log.Fatal("❌ AutoMigrate Comment failed:", err)
	}

	DB = db
	log.Println("🚀 Connected to PostgreSQL & migrated models")
	return db
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
