package services

import (
	"errors"
	"working-day-api/internal/helpers"
)

type LoginService struct {
	UserRepo UserRepository
	JWT      JWTService
	Hasher   helpers.PasswordHasher
}

func (s *LoginService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := s.Hasher.CheckPasswordHash(password, user.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := s.JWT.GenerateToken(user.ID, user.Role.Slug)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
