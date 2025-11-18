package services

import (
	"auth/internal/models"
	"auth/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Register(username, email, password string, role models.Role) (*models.User, error) {

	// hash lozinke
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
		Role:     role,
		Blocked:  false,
	}

	err := s.repo.Create(user)
	return user, err
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}