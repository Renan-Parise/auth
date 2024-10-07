package repositories

import (
	"database/sql"
	"time"

	"github.com/Renan-Parise/codium/database"
	"github.com/Renan-Parise/codium/entities"
	"github.com/Renan-Parise/codium/errors"
	"github.com/Renan-Parise/codium/utils"
)

type UserRepository interface {
	FindByID(id int) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Create(user entities.User) error
	Update(user entities.User) error
	DeactivateUser(userID int) error
	DeleteInactiveUsers() error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindByID(id int) (*entities.User, error) {
	db := database.GetDBInstance()
	user := &entities.User{}
	query := "SELECT id, username, email, password, active FROM users WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewQueryError("User not found")
		}
		return nil, errors.NewQueryError(err.Error())
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	db := database.GetDBInstance()
	user := &entities.User{}
	query := "SELECT id, username, password, active FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Password, &user.Active)
	if err != nil {
		return nil, errors.NewQueryError(err.Error())
	}
	return user, nil
}

func (r *userRepository) Create(user entities.User) error {
	db := database.GetDBInstance()
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	_, err := db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to create user in repository method Create: ", err)

		return errors.NewQueryError(err.Error())
	}
	return nil
}

func (r *userRepository) Update(user entities.User) error {
	db := database.GetDBInstance()
	query := "UPDATE users SET username = ?, password = ? WHERE email = ?"
	_, err := db.Exec(query, user.Username, user.Password, user.Email)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to update user in repository method Update: ", err)

		return errors.NewQueryError(err.Error())
	}
	return nil
}

func (r *userRepository) DeactivateUser(userID int) error {
	db := database.GetDBInstance()
	query := "UPDATE users SET active = ?, deactivatedAt = ? WHERE id = ?"
	_, err := db.Exec(query, false, time.Now(), userID)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to deactivate user in repository method DeactivateUser: ", err)
		return errors.NewQueryError(err.Error())
	}
	return nil
}

func (r *userRepository) DeleteInactiveUsers() error {
	db := database.GetDBInstance()
	query := "DELETE FROM users WHERE is_active = ? AND deactivated_at <= ?"
	cutoffDate := time.Now().AddDate(0, 0, -30)
	result, err := db.Exec(query, false, cutoffDate)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to delete inactive users in repository method DeleteInactiveUsers: ", err)
		return errors.NewQueryError(err.Error())
	}
	rowsAffected, _ := result.RowsAffected()
	utils.GetLogger().Infof("Deleted %d inactive users.", rowsAffected)
	return nil
}
