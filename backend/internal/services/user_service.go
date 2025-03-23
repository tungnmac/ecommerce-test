package services

import (
	"ecommerce-test/internal/models"
	"ecommerce-test/internal/repository"
	"errors"
)

type UserService struct {
	Repo repository.UserRepository
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
