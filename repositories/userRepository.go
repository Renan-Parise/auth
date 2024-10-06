package repositories

import (
	"database/sql"

	"github.com/Renan-Parise/codium/database"
	"github.com/Renan-Parise/codium/entities"
	"github.com/Renan-Parise/codium/errors"
	"github.com/Renan-Parise/codium/utils"
)

type UserRepository interface {
	FindByEmail(email string) (*entities.User, error)
	Create(user entities.User) error
	Update(user entities.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	db := database.GetDBInstance()
	user := &entities.User{}
	query := "SELECT id, username, email, password FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewQueryError("User not found")
		}
		return nil, errors.NewQueryError(err.Error())
	}
	return user, nil
}

func (r *userRepository) Create(user entities.User) error {
	db := database.GetDBInstance()
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	_, err := db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to create user in repository method Create.")

		return errors.NewQueryError(err.Error())
	}
	return nil
}

func (r *userRepository) Update(user entities.User) error {
	db := database.GetDBInstance()
	query := "UPDATE users SET username = ?, password = ? WHERE email = ?"
	_, err := db.Exec(query, user.Username, user.Password, user.Email)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to update user in repository method Update.")

		return errors.NewQueryError(err.Error())
	}
	return nil
}
