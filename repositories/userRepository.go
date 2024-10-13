package repositories

import (
	"time"

	"github.com/Renan-Parise/codium-auth/database"
	"github.com/Renan-Parise/codium-auth/entities"
	"github.com/Renan-Parise/codium-auth/errors"
	"github.com/Renan-Parise/codium-auth/utils"
)

type UserRepository interface {
	FindByID(id int) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Create(user entities.User) error
	Update(ID int, user entities.User) error
	DeactivateUser(ID int) error
	DeleteInactiveUsers() error
	UpdateTwoFACode(user *entities.User) error
	UpdateTwoFASettings(user *entities.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindByID(id int) (*entities.User, error) {
	db := database.GetDBInstance()
	user := &entities.User{}
	query := "SELECT id, username, email, password, active, isTwoFAEnabled, twoFACode, twoFACodeExpiration FROM users WHERE id = ?"

	var twoFACodeExpiresAtStr string

	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.Is2FAEnabled,
		&user.TwoFACode,
		&twoFACodeExpiresAtStr,
	)
	if err != nil {
		return nil, errors.NewQueryError(err.Error())
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", twoFACodeExpiresAtStr)
	if err != nil {
		return nil, errors.NewQueryError("invalid expiration time format: " + err.Error())
	}
	user.TwoFACodeExpiresAt = &parsedTime

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	db := database.GetDBInstance()
	user := &entities.User{}
	query := "SELECT id, username, email, password, active, isTwoFAEnabled, twoFACode, twoFACodeExpiration FROM users WHERE email = ?"

	var twoFACodeExpiresAtStr string

	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.Is2FAEnabled,
		&user.TwoFACode,
		&twoFACodeExpiresAtStr,
	)
	if err != nil {
		return nil, errors.NewQueryError(err.Error())
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", twoFACodeExpiresAtStr)
	if err != nil {
		return nil, errors.NewQueryError("invalid expiration time format: " + err.Error())
	}
	user.TwoFACodeExpiresAt = &parsedTime

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

func (r *userRepository) Update(ID int, user entities.User) error {
	db := database.GetDBInstance()
	query := "UPDATE users SET username = ? WHERE id = ?"
	_, err := db.Exec(query, user.Username, ID)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to update user in repository method Update: ", err)

		return errors.NewQueryError(err.Error())
	}
	return nil
}

func (r *userRepository) DeactivateUser(ID int) error {
	db := database.GetDBInstance()
	query := "UPDATE users SET active = ?, deactivatedAt = ? WHERE id = ?"
	_, err := db.Exec(query, false, time.Now(), ID)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to deactivate user in repository method DeactivateUser: ", err)
		return errors.NewQueryError(err.Error())
	}
	return nil
}

func (r *userRepository) DeleteInactiveUsers() error {
	db := database.GetDBInstance()
	query := "DELETE FROM users WHERE active = ? AND deactivatedAt <= ?"
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

func (r *userRepository) UpdateTwoFACode(user *entities.User) error {
	db := database.GetDBInstance()
	query := "UPDATE users SET twoFACode = ?, twoFACodeExpiration = ? WHERE id = ?"
	_, err := db.Exec(query, user.TwoFACode, user.TwoFACodeExpiresAt, user.ID)
	if err != nil {
		return errors.NewQueryError(err.Error())
	}
	return nil
}

func (r *userRepository) UpdateTwoFASettings(user *entities.User) error {
	db := database.GetDBInstance()
	query := "UPDATE users SET isTwoFAEnabled = ? WHERE id = ?"
	_, err := db.Exec(query, user.Is2FAEnabled, user.ID)
	if err != nil {
		return errors.NewQueryError(err.Error())
	}
	return nil
}
