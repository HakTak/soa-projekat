package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"SHOPPING-CART/internal/model" // Import your models here

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	// 1. Read Environment Variables (passed from Docker Compose)
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Fallback/Default values (good for local testing without docker)
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}

	// 2. Construct the Connection String (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Belgrade",
		host, user, password, dbname, port)

	// 3. Connect to Database
	var db *gorm.DB
	var err error

	// Simple retry mechanism because sometimes the DB takes a second to be ready
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			break
		}

		log.Printf("Failed to connect to database. Retrying in 2 seconds... (%d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after retries: ", err)
	}

	log.Println("Successfully connected to database!")

	// 4. Auto Migrate (Create tables automatically)
	// This creates the 'carts', 'order_items', 'purchase_tokens' tables in Postgres
	err = db.AutoMigrate(
		&model.ShoppingCart{},
		&model.OrderItem{},
		&model.PurchaseToken{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}

	return db
}
