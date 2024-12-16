package services

import (
	"errors"
	"fmt"
	"log"
	"working-day-api/internal/domain"
	"working-day-api/internal/repositories"
	"working-day-api/messaging"
)

type TaskService struct {
	TaskRepo  repositories.TaskRepository
	Messenger messaging.Messenger
	Decrypt   func(encrypted string) (string, error)
	Encrypt   func(string) (string, error)
}

func (s *TaskService) GetTask(taskID string, userID uint, userRole string) (*domain.Task, error) {
	task, err := s.TaskRepo.FindByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	if task.UserID != userID && userRole != "manager" {
		return nil, errors.New("you do not have permission to view this task")
	}

	decryptedSummary, err := s.Decrypt(task.Summary)
	if err != nil {
		return nil, errors.New("error decrypting summary")
	}
	task.Summary = decryptedSummary

	return task, nil
}

func (s *TaskService) GetTasks(userID uint, userRole string) ([]domain.Task, error) {
	var tasks []domain.Task
	var err error

	if userRole == "manager" {
		tasks, err = s.TaskRepo.FindAll()
	} else {
		log.Printf("Calling FindByUserID")
		tasks, err = s.TaskRepo.FindByUserID(userID)
	}

	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		return nil, errors.New("error fetching tasks")
	}

	if len(tasks) == 0 {
		log.Printf("No tasks found")
		return tasks, nil
	}

	for i, task := range tasks {
		decryptedSummary, err := s.Decrypt(task.Summary)
		if err != nil {
			log.Printf("Decryption error for task ID: %d, Summary: %s, Error: %v", task.ID, task.Summary, err)
			tasks[i].Summary = "Error decrypting summary"
			continue
		}
		tasks[i].Summary = decryptedSummary
	}

	return tasks, nil
}

func (s *TaskService) CreateTask(task *domain.Task, userID uint) error {
	task.UserID = userID

	encryptedSummary, err := s.Encrypt(task.Summary)
	if err != nil {
		return errors.New("error encrypting summary")
	}
	task.Summary = encryptedSummary

	if err := s.TaskRepo.Create(task); err != nil {
		return errors.New("error creating task")
	}

	message := fmt.Sprintf("Task Created: ID=%d, Task=%s", task.ID, task.Task)
	if err := s.Messenger.Publish("tasks", []byte(message)); err != nil {
		log.Printf("Error publishing message to RabbitMQ: %v", err)
	}

	return nil
}

func (s *TaskService) UpdateTask(taskID string, userID uint, updatedTask *domain.Task) error {
	task, err := s.TaskRepo.FindByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	if task.UserID != userID {
		return errors.New("you do not have permission to update this task")
	}

	task.Task = updatedTask.Task
	task.PerformedDate = updatedTask.PerformedDate

	encryptedSummary, err := s.Encrypt(updatedTask.Summary)
	if err != nil {
		return errors.New("error encrypting summary")
	}
	task.Summary = encryptedSummary

	if err := s.TaskRepo.Update(task); err != nil {
		return errors.New("error updating task")
	}

	return nil
}

func (s *TaskService) DeleteTask(taskID string, userID uint, userRole string) error {
	task, err := s.TaskRepo.FindByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	if task.UserID != userID && userRole != "manager" {
		return errors.New("you do not have permission to delete this task")
	}

	if err := s.TaskRepo.Delete(task); err != nil {
		return errors.New("error deleting task")
	}

	return nil
}
