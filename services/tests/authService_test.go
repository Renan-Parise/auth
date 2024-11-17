package services

import (
	"testing"

	"github.com/Renan-Parise/auth/entities"
	"github.com/Renan-Parise/auth/errors"
	"github.com/Renan-Parise/auth/services"
	"github.com/stretchr/testify/assert"
)

type mockUserRepository struct {
	users map[string]entities.User
}

func (m *mockUserRepository) UpdatePasswordRecoveryCode(user *entities.User) error {
	panic("unimplemented")
}

func (m *mockUserRepository) UpdatePassword(user *entities.User) error {
	panic("unimplemented")
}

func (m *mockUserRepository) FindByID(id int) (*entities.User, error) {
	panic("unimplemented")
}

func (m *mockUserRepository) DeactivateUser(ID int) error {
	panic("unimplemented")
}

func (m *mockUserRepository) DeleteInactiveUsers() error {
	panic("unimplemented")
}

func (m *mockUserRepository) UpdateTwoFACode(user *entities.User) error {
	panic("unimplemented")
}

func (m *mockUserRepository) UpdateTwoFASettings(user *entities.User) error {
	panic("unimplemented")
}

func (m *mockUserRepository) FindByEmail(email string) (*entities.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, errors.NewQueryError("user not found")
	}
	return &user, nil
}

func (m *mockUserRepository) Create(user entities.User) error {
	if _, exists := m.users[user.Username]; exists {
		return errors.NewQueryError("user already exists")
	}
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepository) Update(ID int, user entities.User) error {
	if _, exists := m.users[user.Username]; !exists {
		return errors.NewQueryError("user not found")
	}
	m.users[user.Username] = user
	return nil
}

func TestRegister(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]entities.User)}
	service := services.NewAuthService(repo, nil)

	user := entities.User{
		Username: "testuser",
		Password: "password123",
		Email:    "testuser@example.com",
	}

	err := service.Register(user)
	assert.Nil(t, err)

	err = service.Register(user)
	assert.NotNil(t, err)
}

func TestLogin(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]entities.User)}
	service := services.NewAuthService(repo, nil)

	user := entities.User{
		Username: "testuser",
		Password: "password123",
	}

	err := service.Register(user)
	assert.Nil(t, err)

	token, err := service.Login("testuser", "password123")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	_, err = service.Login("testuser", "wrongpassword")
	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]entities.User)}
	service := services.NewAuthService(repo, nil)

	user := entities.User{
		Username: "testuser",
		Password: "password123",
	}

	err := service.Register(user)
	assert.Nil(t, err)

	ID := 0

	user.Password = "newpassword123"
	err = service.Update(ID, user)
	assert.Nil(t, err)

	_, err = service.Login("testuser", "password123")
	assert.NotNil(t, err)

	token, err := service.Login("testuser", "newpassword123")
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}
