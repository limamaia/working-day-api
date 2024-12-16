package services

import (
	"errors"
	"fmt"
	"working-day-api/internal/domain"
	"working-day-api/internal/helpers"
)

type UserService struct {
	UserRepo UserRepository
	Hasher   helpers.PasswordHasher
}

func (s *UserService) GetUser(requestedUserID string, loggedUserID uint, loggedUserRole string) (*domain.User, error) {
	if loggedUserRole != "manager" && requestedUserID != fmt.Sprintf("%d", loggedUserID) {
		return nil, errors.New("you do not have permission to view this user")
	}

	user, err := s.UserRepo.FindByID(requestedUserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Password = ""

	return user, nil
}

func (s *UserService) CreateUser(user *domain.User) error {

	existingUser, err := s.UserRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("a user with this email already exists")
	}

	hashedPassword, err := s.Hasher.HashPassword(user.Password)
	if err != nil {
		return errors.New("error hashing password")
	}
	user.Password = hashedPassword

	if err := s.UserRepo.Create(user); err != nil {
		return errors.New("error creating user")
	}

	return nil
}
