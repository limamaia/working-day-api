package repositories

import "working-day-api/internal/domain"

type UserRepository interface {
	FindByID(id string) (*domain.User, error)
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}
