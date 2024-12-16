package services

import "working-day-api/internal/domain"

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	Create(user *domain.User) error
}

type JWTService interface {
	GenerateToken(userID uint, role string) (string, error)
	ParseToken(token string) (*domain.Claim, error)
}
