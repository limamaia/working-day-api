package services_test

import (
	"errors"
	"testing"
	"working-day-api/internal/domain"
	"working-day-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := &services.UserService{
		UserRepo: mockRepo,
	}

	t.Run("SuccessForManager", func(t *testing.T) {
		mockUser := &domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		}

		mockRepo.On("FindByID", "1").Return(mockUser, nil)

		user, err := service.GetUser("1", 2, "manager")

		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)
		mockRepo.AssertCalled(t, "FindByID", "1")
	})

	t.Run("SuccessForSelf", func(t *testing.T) {
		mockUser := &domain.User{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		}

		mockRepo.On("FindByID", "1").Return(mockUser, nil)

		user, err := service.GetUser("1", 1, "user")

		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)
		mockRepo.AssertCalled(t, "FindByID", "1")
	})

	t.Run("PermissionDenied", func(t *testing.T) {
		user, err := service.GetUser("2", 1, "user")

		assert.EqualError(t, err, "you do not have permission to view this user")
		assert.Nil(t, user)
		mockRepo.AssertNotCalled(t, "FindByID")
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.On("FindByID", "3").Return(nil, errors.New("user not found"))

		user, err := service.GetUser("3", 1, "manager")

		assert.EqualError(t, err, "user not found")
		assert.Nil(t, user)
		mockRepo.AssertCalled(t, "FindByID", "3")
	})
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockHasher := new(MockPasswordHasher)

	service := &services.UserService{
		UserRepo: mockRepo,
		Hasher:   mockHasher,
	}

	t.Run("Success", func(t *testing.T) {
		mockUser := &domain.User{
			ID:       1,
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "plain_password",
		}

		mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)
		mockHasher.On("HashPassword", "plain_password").Return("hashed_password", nil)
		mockRepo.On("Create", mockUser).Return(nil)

		err := service.CreateUser(mockUser)

		assert.NoError(t, err)
		assert.Equal(t, "hashed_password", mockUser.Password)
		mockRepo.AssertCalled(t, "FindByEmail", "test@example.com")
		mockHasher.AssertCalled(t, "HashPassword", "plain_password")
		mockRepo.AssertCalled(t, "Create", mockUser)
	})

	t.Run("DuplicateEmailError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		resetMock(&mockHasher.Mock)
		existingUser := &domain.User{
			ID:    1,
			Email: "duplicate@example.com",
		}
		mockUser := &domain.User{
			ID:       2,
			Name:     "New User",
			Email:    "duplicate@example.com",
			Password: "plain_password",
		}

		mockRepo.On("FindByEmail", "duplicate@example.com").Return(existingUser, nil)

		err := service.CreateUser(mockUser)

		assert.EqualError(t, err, "a user with this email already exists")
		mockRepo.AssertCalled(t, "FindByEmail", "duplicate@example.com")
		mockHasher.AssertNotCalled(t, "HashPassword", mock.Anything)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("HashingError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		resetMock(&mockHasher.Mock)
		mockUser := &domain.User{
			ID:       3,
			Name:     "Test User",
			Email:    "test2@example.com",
			Password: "plain_password",
		}

		mockRepo.On("FindByEmail", "test2@example.com").Return(nil, nil)
		mockHasher.On("HashPassword", "plain_password").Return("", errors.New("hashing error"))

		err := service.CreateUser(mockUser)

		assert.EqualError(t, err, "error hashing password")
		mockRepo.AssertCalled(t, "FindByEmail", "test2@example.com")
		mockHasher.AssertCalled(t, "HashPassword", "plain_password")
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		resetMock(&mockRepo.Mock)
		resetMock(&mockHasher.Mock)
		mockUser := &domain.User{
			ID:       4,
			Name:     "Test User",
			Email:    "test3@example.com",
			Password: "plain_password",
		}

		mockRepo.On("FindByEmail", "test3@example.com").Return(nil, nil)
		mockHasher.On("HashPassword", "plain_password").Return("hashed_password", nil)
		mockRepo.On("Create", mockUser).Return(errors.New("repository error"))

		err := service.CreateUser(mockUser)

		assert.EqualError(t, err, "error creating user")
		mockRepo.AssertCalled(t, "FindByEmail", "test3@example.com")
		mockHasher.AssertCalled(t, "HashPassword", "plain_password")
		mockRepo.AssertCalled(t, "Create", mockUser)
	})
}
