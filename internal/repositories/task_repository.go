package repositories

import (
	"working-day-api/database"
	"working-day-api/internal/domain"
)

type TaskRepositoryImpl struct{}

func (r *TaskRepositoryImpl) FindByID(id string) (*domain.Task, error) {
	var task domain.Task
	if err := database.DB.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepositoryImpl) Create(task *domain.Task) error {
	return database.DB.Create(task).Error
}

func (r *TaskRepositoryImpl) Update(task *domain.Task) error {
	return database.DB.Save(task).Error
}

func (r *TaskRepositoryImpl) Delete(task *domain.Task) error {
	return database.DB.Delete(task).Error
}

func (r *TaskRepositoryImpl) FindAll() ([]domain.Task, error) {
	var tasks []domain.Task
	if err := database.DB.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepositoryImpl) FindByUserID(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	if err := database.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
