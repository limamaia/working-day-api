package services_test

import (
	"errors"
	"testing"
	"working-day-api/internal/domain"
	"working-day-api/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID uint, role string) (string, error) {
	args := m.Called(userID, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ParseToken(token string) (*domain.Claim, error) {
	return nil, nil
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) CheckPasswordHash(password, hash string) error {
	args := m.Called(password, hash)
	return args.Error(0)
}

func (m *MockPasswordHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func TestLoginService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJWT := new(MockJWTService)
	mockHasher := new(MockPasswordHasher)

	service := &services.LoginService{
		UserRepo: mockRepo,
		JWT:      mockJWT,
		Hasher:   mockHasher,
	}

	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: "$2a$10$hashSimuladoParaTeste", // Hash simulado
		Role:     &domain.Role{Slug: "user"},
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockJWT.ExpectedCalls = nil
		mockHasher.ExpectedCalls = nil

		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
		mockHasher.On("CheckPasswordHash", "password", user.Password).Return(nil)
		mockJWT.On("GenerateToken", user.ID, "user").Return("mockToken", nil)

		token, err := service.Login("test@example.com", "password")

		assert.NoError(t, err)
		assert.Equal(t, "mockToken", token)
		mockRepo.AssertCalled(t, "FindByEmail", "test@example.com")
		mockHasher.AssertCalled(t, "CheckPasswordHash", "password", user.Password)
		mockJWT.AssertCalled(t, "GenerateToken", user.ID, "user")
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockJWT.ExpectedCalls = nil
		mockHasher.ExpectedCalls = nil

		mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.New("user not found"))

		token, err := service.Login("notfound@example.com", "password")

		assert.EqualError(t, err, "user not found")
		assert.Empty(t, token)
		mockRepo.AssertCalled(t, "FindByEmail", "notfound@example.com")
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockJWT.ExpectedCalls = nil
		mockHasher.ExpectedCalls = nil

		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
		mockHasher.On("CheckPasswordHash", "wrongpassword", user.Password).Return(errors.New("invalid credentials"))

		token, err := service.Login("test@example.com", "wrongpassword")

		assert.EqualError(t, err, "invalid credentials")
		assert.Empty(t, token)
		mockRepo.AssertCalled(t, "FindByEmail", "test@example.com")
		mockHasher.AssertCalled(t, "CheckPasswordHash", "wrongpassword", user.Password)
	})

	t.Run("TokenGenerationError", func(t *testing.T) {
		mockRepo.ExpectedCalls = nil
		mockJWT.ExpectedCalls = nil
		mockHasher.ExpectedCalls = nil

		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
		mockHasher.On("CheckPasswordHash", "password", user.Password).Return(nil)
		mockJWT.On("GenerateToken", user.ID, "user").Return("", errors.New("failed to generate token"))

		token, err := service.Login("test@example.com", "password")

		assert.EqualError(t, err, "failed to generate token")
		assert.Empty(t, token)
		mockRepo.AssertCalled(t, "FindByEmail", "test@example.com")
		mockHasher.AssertCalled(t, "CheckPasswordHash", "password", user.Password)
		mockJWT.AssertCalled(t, "GenerateToken", user.ID, "user")
	})
}
