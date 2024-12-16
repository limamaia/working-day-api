package repositories

import (
	"working-day-api/database"
	"working-day-api/internal/domain"
)

type UserRepositoryImpl struct{}

func (r *UserRepositoryImpl) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := database.DB.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmailAll(email string) (*domain.User, error) {
	var user domain.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) FindByID(id string) (*domain.User, error) {
	var user domain.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Create(user *domain.User) error {
	return database.DB.Create(user).Error
}
