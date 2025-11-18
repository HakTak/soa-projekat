package main

import (
	"log"
	stdhttp "net/http"
	"os"

	"stakeholders/internal/api/http"
	"stakeholders/internal/model"
	"stakeholders/internal/repository"
	"stakeholders/internal/service"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=stakeholders password=secret dbname=stakeholders port=5432 sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// Automigrate profile table (creates table if missing). Safe for dev; consider controlled migrations in prod.
	if err := db.AutoMigrate(&model.Profile{}); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	var count int64
	db.Model(&model.Profile{}).Count(&count)
	if count == 0 {
		seedProfiles(db)
	}

	repo := repository.NewGormProfileRepo(db)
	svc := service.NewProfileService(repo)
	h := http.NewHandler(svc)

	r := chi.NewRouter()
	r.Use(http.AuthMiddleware)
	h.RegisterRoutes(r)

	addr := ":8081"
	log.Printf("listening on %s", addr)
	if err := stdhttp.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func seedProfiles(db *gorm.DB) {
	profiles := []model.Profile{
		{
			UserID:         "123e4567-e89b-12d3-a456-426614174001",
			FirstName:      "Alice",
			LastName:       "Smith",
			ProfilePicture: "",
			Biography:      "Tour guide from NYC.",
			Motto:          "Adventure awaits!",
			Role:           model.RoleGuide,
			IsBlocked:      false,
		},
		{
			UserID:         "123e4567-e89b-12d3-a456-426614174002",
			FirstName:      "Bob",
			LastName:       "Brown",
			ProfilePicture: "",
			Biography:      "Tourist from CA.",
			Motto:          "Live, laugh, travel.",
			Role:           model.RoleTourist,
			IsBlocked:      false,
		},
		{
			UserID:         "123e4567-e89b-12d3-a456-426614174003",
			FirstName:      "Carol",
			LastName:       "White",
			ProfilePicture: "",
			Biography:      "Admin user.",
			Motto:          "Keeping things running.",
			Role:           model.RoleAdmin,
			IsBlocked:      false,
		},
	}

	result := db.Create(&profiles)
	if result.Error != nil {
		log.Fatalf("failed to seed profiles: %v", result.Error)
	}

	log.Println("Seeded profiles successfully")
}
