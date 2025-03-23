package services

import (
	"ecommerce-test/internal/models"
	"ecommerce-test/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo repository.UserRepository
}

// RegisterUser hashes password and saves user to database
func (s *AuthService) RegisterUser(user models.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	// Save user to database
	return s.UserRepo.CreateUser(user)
}

// AuthenticateUser verifies email and password
func (s *AuthService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
