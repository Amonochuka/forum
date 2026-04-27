package auth

import (
	"errors"
	"testing"
)

// MockRepo implements the UserRepo interface
type MockRepo struct {
	GetUserByEmailFunc func(email string) (User, error)
	CreateUserFunc     func(user User) error
}

func (m *MockRepo) GetUserByEmail(email string) (User, error) { return m.GetUserByEmailFunc(email) }
func (m *MockRepo) CreateUser(user User) error                { return m.CreateUserFunc(user) }
func (m *MockRepo) GetUserByID(id int) (User, error)          { return User{}, nil }

func TestRegister_EmailTaken(t *testing.T) {
	// 1. Setup: Define behavior for the mock
	mock := &MockRepo{
		GetUserByEmailFunc: func(email string) (User, error) {
			// Simulate that the email ALREADY exists in the DB
			return User{Email: email}, nil
		},
	}

	service := NewService(mock)

	// 2. Execute
	err := service.Register(User{Username: "test", Email: "taken@test.com", Password: "password123"})

	// 3. Verify: Did we get the expected error?
	if !errors.Is(err, ErrEmailTaken) {
		t.Errorf("expected ErrEmailTaken, got %v", err)
	}
}
