package repositories

import (
	"strconv"

	"github.com/Renan-Parise/codium/entities"
	"github.com/Renan-Parise/codium/errors"
)

type MockUserRepository struct {
	Users map[string]entities.User
}

func NewMockUserRepository() UserRepository {
	return &MockUserRepository{
		Users: make(map[string]entities.User),
	}
}

func (m *MockUserRepository) DeactivateUser(ID int) error {
	panic("unimplemented")
}

func (m *MockUserRepository) FindByID(id int) (*entities.User, error) {
	stringID := strconv.Itoa(id)
	user, exists := m.Users[stringID]
	if !exists {
		return nil, errors.NewQueryError("user not found")
	}

	return &user, nil
}

func (m *MockUserRepository) DeleteInactiveUsers() error {
	panic("unimplemented")
}

func (m *MockUserRepository) FindByEmail(email string) (*entities.User, error) {
	user, exists := m.Users[email]
	if !exists {
		return nil, errors.NewQueryError("user not found")
	}
	return &user, nil
}

func (m *MockUserRepository) Create(user entities.User) error {
	if _, exists := m.Users[user.Username]; exists {
		return errors.NewQueryError("user already exists")
	}
	m.Users[user.Username] = user
	return nil
}

func (m *MockUserRepository) Update(ID int, user entities.User) error {
	if _, exists := m.Users[user.Username]; !exists {
		return errors.NewQueryError("user not found")
	}
	m.Users[user.Username] = user
	return nil
}
