package repositories

import "working-day-api/internal/domain"

type TaskRepository interface {
	FindByID(id string) (*domain.Task, error)
	Create(task *domain.Task) error
	Update(task *domain.Task) error
	Delete(task *domain.Task) error
	FindAll() ([]domain.Task, error)
	FindByUserID(userID uint) ([]domain.Task, error)
}
