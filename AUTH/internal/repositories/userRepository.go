package repositories

import (
	"auth/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Register(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepository) GetUsersForAdmin() ([]models.UserNoPassDTO, error) {
	var result []models.UserNoPassDTO

	err := r.db.
		Model(&models.User{}).
		Select("id", "username", "email", "role", "blocked").
		Where("role IN ?", []models.Role{models.RoleGuide, models.RoleTourist}).
		Find(&result).Error

	return result, err
}
