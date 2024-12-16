package services_test

import (
	"errors"
	"log"
	"testing"
	"time"
	"working-day-api/internal/domain"
	"working-day-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

type MockMessenger struct {
	mock.Mock
}

func resetMock(mock *mock.Mock) {
	mock.ExpectedCalls = nil
	mock.Calls = nil
}

func (m *MockTaskRepository) FindByID(id string) (*domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) FindAll() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskRepository) FindByUserID(userID uint) ([]domain.Task, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockMessenger) Publish(topic string, message []byte) error {
	args := m.Called(topic, message)
	return args.Error(0)
}

func (m *MockTaskRepository) Create(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Update(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func mockDecrypt(input string) (string, error) {
	log.Printf("Decrypt called with input: %s", input)
	if input == "error_summary" {
		return "", errors.New("decryption error")
	}
	return "Decrypted Summary", nil
}

var mockEncryptBehavior = func(input string) (string, error) {
	if input == "error_summary" {
		return "", errors.New("encryption error")
	}
	return "encrypted_" + input, nil
}

func mockEncrypt(input string) (string, error) {
	return mockEncryptBehavior(input)
}

func TestTaskService_GetTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	service := &services.TaskService{
		TaskRepo: mockRepo,
		Decrypt:  mockDecrypt,
	}

	t.Run("Success", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		mockTask := &domain.Task{
			ID:      1,
			UserID:  1,
			Summary: "encrypted_summary",
		}

		mockRepo.On("FindByID", "1").Return(mockTask, nil)

		task, err := service.GetTask("1", 1, "user")

		assert.NoError(t, err)
		assert.Equal(t, "Decrypted Summary", task.Summary)
		mockRepo.AssertCalled(t, "FindByID", "1")
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		mockRepo.On("FindByID", "2").Return(nil, errors.New("task not found"))

		task, err := service.GetTask("2", 1, "user")

		assert.EqualError(t, err, "task not found")
		assert.Nil(t, task)
		mockRepo.AssertCalled(t, "FindByID", "2")
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		mockTask := &domain.Task{
			ID:      1,
			UserID:  2,
			Summary: "encrypted_summary",
		}

		mockRepo.On("FindByID", "3").Return(mockTask, nil)

		task, err := service.GetTask("3", 1, "user")

		assert.EqualError(t, err, "you do not have permission to view this task")
		assert.Nil(t, task)
		mockRepo.AssertCalled(t, "FindByID", "3")
	})

	t.Run("DecryptionError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		mockTask := &domain.Task{
			ID:      1,
			UserID:  1,
			Summary: "error_summary",
		}

		mockRepo.On("FindByID", "4").Return(mockTask, nil)

		task, err := service.GetTask("4", 1, "user")

		assert.EqualError(t, err, "error decrypting summary")
		assert.Nil(t, task)
		mockRepo.AssertCalled(t, "FindByID", "4")
	})
}

func TestTaskService_GetTasks(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	service := &services.TaskService{
		TaskRepo: mockRepo,
		Decrypt:  mockDecrypt,
	}

	t.Run("ManagerGetsAllTasks", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		mockTasks := []domain.Task{
			{ID: 1, Summary: "encrypted_summary_1"},
			{ID: 2, Summary: "encrypted_summary_2"},
		}
		mockRepo.On("FindAll").Return(mockTasks, nil)

		tasks, err := service.GetTasks(0, "manager")

		assert.NoError(t, err)
		assert.Len(t, tasks, 2)
		assert.Equal(t, "Decrypted Summary", tasks[0].Summary)
		assert.Equal(t, "Decrypted Summary", tasks[1].Summary)
		mockRepo.AssertCalled(t, "FindAll")
	})

	t.Run("UserGetsOwnTasks", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		mockTasks := []domain.Task{
			{ID: 3, UserID: 1, Summary: "encrypted_summary_3"},
		}
		mockRepo.On("FindByUserID", uint(1)).Return(mockTasks, nil)

		tasks, err := service.GetTasks(1, "user")

		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, "Decrypted Summary", tasks[0].Summary)
		mockRepo.AssertCalled(t, "FindByUserID", uint(1))
	})

	t.Run("ErrorFetchingTasks", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		mockRepo.On("FindByUserID", uint(2)).Return([]domain.Task{}, errors.New("repository error"))

		tasks, err := service.GetTasks(2, "user")

		assert.EqualError(t, err, "error fetching tasks")
		assert.Nil(t, tasks)
		mockRepo.AssertCalled(t, "FindByUserID", uint(2))
	})

	t.Run("NoTasksFound", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		mockRepo.On("FindByUserID", uint(3)).Return([]domain.Task{}, nil)

		tasks, err := service.GetTasks(3, "user")

		assert.NoError(t, err)
		assert.Len(t, tasks, 0)
		mockRepo.AssertCalled(t, "FindByUserID", uint(3))
	})

	t.Run("DecryptionError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		mockTasks := []domain.Task{
			{ID: 4, UserID: 4, Summary: "error_summary"},
		}

		mockRepo.On("FindByUserID", uint(4)).Return(mockTasks, nil)

		tasks, err := service.GetTasks(4, "user")

		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, "Error decrypting summary", tasks[0].Summary)
		mockRepo.AssertCalled(t, "FindByUserID", uint(4))
	})
}

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	mockMessenger := new(MockMessenger)

	service := &services.TaskService{
		TaskRepo:  mockRepo,
		Messenger: mockMessenger,
		Decrypt:   mockDecrypt,
		Encrypt:   mockEncrypt,
	}

	t.Run("Success", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		resetMock(&mockMessenger.Mock)

		mockTask := &domain.Task{
			ID:      1,
			Task:    "Test Task",
			Summary: "plain_summary",
		}

		mockRepo.On("Create", mockTask).Return(nil)
		mockMessenger.On("Publish", "tasks", mock.Anything).Return(nil)

		err := service.CreateTask(mockTask, 1)

		assert.NoError(t, err)
		assert.Equal(t, "encrypted_plain_summary", mockTask.Summary)
		mockRepo.AssertCalled(t, "Create", mockTask)
		mockMessenger.AssertCalled(t, "Publish", "tasks", mock.Anything)
	})

	t.Run("EncryptionError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		resetMock(&mockMessenger.Mock)

		mockTask := &domain.Task{
			ID:      2,
			Task:    "Test Task",
			Summary: "error_summary",
		}

		err := service.CreateTask(mockTask, 1)

		assert.EqualError(t, err, "error encrypting summary")
		mockRepo.AssertNotCalled(t, "Create", mockTask)
		mockMessenger.AssertNotCalled(t, "Publish", "tasks", mock.Anything)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		resetMock(&mockMessenger.Mock)

		mockTask := &domain.Task{
			ID:      3,
			Task:    "Test Task",
			Summary: "plain_summary",
		}

		mockRepo.On("Create", mockTask).Return(errors.New("repository error"))

		err := service.CreateTask(mockTask, 1)

		assert.EqualError(t, err, "error creating task")
		mockRepo.AssertCalled(t, "Create", mockTask)
		mockMessenger.AssertNotCalled(t, "Publish", "tasks", mock.Anything)
	})

	t.Run("PublishError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		resetMock(&mockMessenger.Mock)

		mockTask := &domain.Task{
			ID:      4,
			Task:    "Test Task",
			Summary: "plain_summary",
		}

		mockRepo.On("Create", mockTask).Return(nil)
		mockMessenger.On("Publish", "tasks", mock.Anything).Return(errors.New("publish error"))

		err := service.CreateTask(mockTask, 1)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", mockTask)
		mockMessenger.AssertCalled(t, "Publish", "tasks", mock.Anything)
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	service := &services.TaskService{
		TaskRepo: mockRepo,
		Encrypt:  mockEncrypt,
	}

	t.Run("Success", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "1"
		userID := uint(1)
		updatedSummary := "new_summary"

		existingTask := &domain.Task{
			ID:      1,
			UserID:  userID,
			Task:    "Old Task",
			Summary: "old_summary",
		}

		updatedTask := &domain.Task{
			Task:          "Updated Task",
			Summary:       updatedSummary,
			PerformedDate: time.Now(),
		}

		mockEncryptBehavior = func(input string) (string, error) {
			if input == updatedSummary {
				return "encrypted_summary", nil
			}
			return "", errors.New("unexpected input")
		}

		mockRepo.On("FindByID", taskID).Return(existingTask, nil)
		mockRepo.On("Update", mock.AnythingOfType("*domain.Task")).Return(nil)

		err := service.UpdateTask(taskID, userID, updatedTask)

		assert.NoError(t, err)
		assert.Equal(t, "Updated Task", existingTask.Task)
		assert.Equal(t, "encrypted_summary", existingTask.Summary)
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertCalled(t, "Update", existingTask)
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "2"
		userID := uint(1)

		updatedTask := &domain.Task{
			Task:          "Updated Task",
			Summary:       "new_summary",
			PerformedDate: time.Now(),
		}

		mockRepo.On("FindByID", taskID).Return(nil, errors.New("task not found"))

		err := service.UpdateTask(taskID, userID, updatedTask)

		assert.EqualError(t, err, "task not found")
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "3"
		userID := uint(2)

		existingTask := &domain.Task{
			ID:      3,
			UserID:  1,
			Task:    "Old Task",
			Summary: "old_summary",
		}

		updatedTask := &domain.Task{
			Task:          "Updated Task",
			Summary:       "new_summary",
			PerformedDate: time.Now(),
		}

		mockRepo.On("FindByID", taskID).Return(existingTask, nil)

		err := service.UpdateTask(taskID, userID, updatedTask)

		assert.EqualError(t, err, "you do not have permission to update this task")
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
	})

	t.Run("EncryptionError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "4"
		userID := uint(1)
		updatedSummary := "error_summary"

		existingTask := &domain.Task{
			ID:      4,
			UserID:  userID,
			Task:    "Old Task",
			Summary: "old_summary",
		}

		updatedTask := &domain.Task{
			Task:          "Updated Task",
			Summary:       updatedSummary,
			PerformedDate: time.Now(),
		}

		mockEncryptBehavior = func(input string) (string, error) {
			if input == updatedSummary {
				return "", errors.New("encryption error")
			}
			return "encrypted_" + input, nil
		}

		mockRepo.On("FindByID", taskID).Return(existingTask, nil)

		err := service.UpdateTask(taskID, userID, updatedTask)

		assert.EqualError(t, err, "error encrypting summary")
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
	})

	t.Run("UpdateError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "5"
		userID := uint(1)

		existingTask := &domain.Task{
			ID:      5,
			UserID:  userID,
			Task:    "Old Task",
			Summary: "old_summary",
		}

		updatedTask := &domain.Task{
			Task:          "Updated Task",
			Summary:       "new_summary",
			PerformedDate: time.Now(),
		}

		mockRepo.On("FindByID", taskID).Return(existingTask, nil)
		mockRepo.On("Update", mock.Anything).Return(errors.New("repository error"))

		err := service.UpdateTask(taskID, userID, updatedTask)

		assert.EqualError(t, err, "error updating task")
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertCalled(t, "Update", existingTask)
	})
}

func TestTaskService_DeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)

	service := &services.TaskService{
		TaskRepo: mockRepo,
	}

	t.Run("Success", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "1"
		userID := uint(1)
		userRole := "user"

		existingTask := &domain.Task{
			ID:     1,
			UserID: userID,
			Task:   "Task to delete",
		}

		mockRepo.On("FindByID", taskID).Return(existingTask, nil)
		mockRepo.On("Delete", existingTask).Return(nil)

		err := service.DeleteTask(taskID, userID, userRole)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertCalled(t, "Delete", existingTask)
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "2"
		userID := uint(1)
		userRole := "user"

		mockRepo.On("FindByID", taskID).Return(nil, errors.New("task not found"))

		err := service.DeleteTask(taskID, userID, userRole)

		assert.EqualError(t, err, "task not found")
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "3"
		userID := uint(1)
		userRole := "user"

		existingTask := &domain.Task{
			ID:     3,
			UserID: 2,
			Task:   "Task to delete",
		}

		mockRepo.On("FindByID", taskID).Return(existingTask, nil)

		err := service.DeleteTask(taskID, userID, userRole)

		assert.EqualError(t, err, "you do not have permission to delete this task")
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
	})

	t.Run("ManagerCanDelete", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "4"
		userID := uint(1)
		userRole := "manager"

		existingTask := &domain.Task{
			ID:     4,
			UserID: 2,
			Task:   "Task to delete",
		}

		mockRepo.On("FindByID", taskID).Return(existingTask, nil)
		mockRepo.On("Delete", existingTask).Return(nil)

		err := service.DeleteTask(taskID, userID, userRole)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertCalled(t, "Delete", existingTask)
	})

	t.Run("DeleteError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		taskID := "5"
		userID := uint(1)
		userRole := "user"

		existingTask := &domain.Task{
			ID:     5,
			UserID: userID,
			Task:   "Task to delete",
		}

		mockRepo.On("FindByID", taskID).Return(existingTask, nil)
		mockRepo.On("Delete", existingTask).Return(errors.New("repository error"))

		err := service.DeleteTask(taskID, userID, userRole)

		assert.EqualError(t, err, "error deleting task")
		mockRepo.AssertCalled(t, "FindByID", taskID)
		mockRepo.AssertCalled(t, "Delete", existingTask)
	})
}
