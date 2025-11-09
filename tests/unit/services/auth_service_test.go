package services

import (
	"testing"

	"github.com/AndikaPrasetia/pos-cafee/internal/models"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepo is a mock implementation of UserRepo interface
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUser(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) CreateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) UpdateUser(user *models.User) (*models.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) UpdateUserPassword(userID, hashedPassword string) error {
	args := m.Called(userID, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateUserStatus(userID string, isActive bool) error {
	args := m.Called(userID, isActive)
	return args.Error(0)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	authService := NewAuthService(mockUserRepo, "test-secret-key", 24*60*60) // 24 hours expiry

	testUser := &models.User{
		ID:        "test-id",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "$2a$10$N9qo8uLOickgx2ZMRZoMye.IjdQcVrLgKmKxqYv4qP7LKJLzDJoFu", // bcrypt hash for "password"
		Role:      types.UserRoleCashier,
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	loginData := &models.UserLogin{
		Username: "testuser",
		Password: "password",
	}

	mockUserRepo.On("GetUserByUsername", "testuser").Return(testUser, nil)

	result, err := authService.Login(loginData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	authService := NewAuthService(mockUserRepo, "test-secret-key", 24*60*60)

	loginData := &models.UserLogin{
		Username: "nonexistent",
		Password: "password",
	}

	mockUserRepo.On("GetUserByUsername", "nonexistent").Return((*models.User)(nil), nil) // Return nil user without error

	result, err := authService.Login(loginData)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
	assert.Nil(t, result)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Register_NewUser(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	authService := NewAuthService(mockUserRepo, "test-secret-key", 24*60*60)

	registerData := &models.UserRegister{
		Username:  "newuser",
		Email:     "newuser@example.com",
		Password:  "Password123!",
		Role:      types.UserRoleCashier,
		FirstName: "New",
		LastName:  "User",
	}

	testUser := &models.User{
		ID:        "generated-id",
		Username:  registerData.Username,
		Email:     registerData.Email,
		Password:  "$2a$10$N9qo8uLOickgx2ZMRZoMye.IjdQcVrLgKmKxqYv4qP7LKJLzDJoFu", // bcrypt hash
		Role:      registerData.Role,
		FirstName: registerData.FirstName,
		LastName:  registerData.LastName,
		IsActive:  true,
	}

	mockUserRepo.On("GetUserByUsername", "newuser").Return((*models.User)(nil), nil) // No existing user
	mockUserRepo.On("GetUserByEmail", "newuser@example.com").Return((*models.User)(nil), nil) // No existing email
	mockUserRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(testUser, nil)

	result, err := authService.Register(registerData)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Data)

	mockUserRepo.AssertExpectations(t)
}