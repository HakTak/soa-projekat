package services

import (
	"auth/internal/models"
	"auth/internal/repositories"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Register(userRegister models.UserRegisterDTO) (*models.User, error) {

	// hash lozinke
	hashed, _ := bcrypt.GenerateFromPassword([]byte(userRegister.Password), 12)

	user := &models.User{
		Username: userRegister.Username,
		Email:    userRegister.Email,
		Password: string(hashed),
		Role:     userRegister.Role,
	}

	err := s.repo.Register(user)
	return user, err
}

func (s *UserService) Login(username, password string) (*models.User, error) {
	user, err := s.repo.GetByUsername(username)

	if err != nil {
		return nil, fmt.Errorf("wrong username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("wrong password")
	}

	return user, nil
}

func (s *UserService) GetUsersForAdmin() ([]models.UserNoPassDTO, error) {
	return s.repo.GetUsersForAdmin()
}

func (s *UserService) GetById(id string) (*models.UserNoPassDTO, error) {
	return s.repo.GetById(id)
}

func (s *UserService) BlockUser(id string) (*models.UserNoPassDTO, error) {
	return s.repo.BlockUser(id)
}
