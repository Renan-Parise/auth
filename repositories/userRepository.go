package repositories

import (
	"database/sql"
	"time"

	"github.com/Renan-Parise/auth/database"
	"github.com/Renan-Parise/auth/entities"
	"github.com/Renan-Parise/auth/errors"
	"github.com/Renan-Parise/auth/utils"
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
	UpdatePasswordRecoveryCode(user *entities.User) error
	UpdatePassword(user *entities.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindByID(id int) (*entities.User, error) {
	db := database.GetDBInstance()
	user := &entities.User{}
	query := "SELECT id, username, email, password, active, isTwoFAEnabled, twoFACode, twoFACodeExpiration, passwordRecoveryCode, recoveryCodeExpiration FROM users WHERE id = ?"

	var twoFACodeExpiresAtStr sql.NullString
	var recoveryCodeExpiresAtStr sql.NullString

	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.Is2FAEnabled,
		&user.TwoFACode,
		&twoFACodeExpiresAtStr,
		&user.PasswordRecoveryCode,
		&recoveryCodeExpiresAtStr,
	)
	if err != nil {
		return nil, errors.NewQueryError(err.Error())
	}

	if twoFACodeExpiresAtStr.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", twoFACodeExpiresAtStr.String)
		if err != nil {
			return nil, errors.NewQueryError("invalid expiration time format: " + err.Error())
		}
		user.TwoFACodeExpiresAt = &parsedTime
	} else {
		user.TwoFACodeExpiresAt = nil
	}

	if recoveryCodeExpiresAtStr.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", recoveryCodeExpiresAtStr.String)
		if err != nil {
			return nil, errors.NewQueryError("invalid expiration time format: " + err.Error())
		}
		user.RecoveryCodeExpiresAt = &parsedTime
	} else {
		user.RecoveryCodeExpiresAt = nil
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	db := database.GetDBInstance()
	user := &entities.User{}
	query := "SELECT id, username, email, password, active, isTwoFAEnabled, twoFACode, twoFACodeExpiration, passwordRecoveryCode, recoveryCodeExpiration FROM users WHERE email = ?"

	var twoFACodeExpiresAt sql.NullString
	var recoveryCodeExpiresAt sql.NullString

	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Active,
		&user.Is2FAEnabled,
		&user.TwoFACode,
		&twoFACodeExpiresAt,
		&user.PasswordRecoveryCode,
		&recoveryCodeExpiresAt,
	)
	if err != nil {
		return nil, errors.NewQueryError(err.Error())
	}

	if twoFACodeExpiresAt.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", twoFACodeExpiresAt.String)
		if err != nil {
			return nil, errors.NewQueryError("invalid expiration time format: " + err.Error())
		}
		user.TwoFACodeExpiresAt = &parsedTime
	} else {
		user.TwoFACodeExpiresAt = nil
	}

	if recoveryCodeExpiresAt.Valid {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", recoveryCodeExpiresAt.String)
		if err != nil {
			return nil, errors.NewQueryError("invalid expiration time format: " + err.Error())
		}
		user.RecoveryCodeExpiresAt = &parsedTime
	} else {
		user.RecoveryCodeExpiresAt = nil
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

func (r *userRepository) UpdatePasswordRecoveryCode(user *entities.User) error {
	db := database.GetDBInstance()
	query := "UPDATE users SET passwordRecoveryCode = ?, recoveryCodeExpiration = ? WHERE id = ?"
	_, err := db.Exec(query, user.PasswordRecoveryCode, user.RecoveryCodeExpiresAt, user.ID)
	if err != nil {
		return errors.NewQueryError(err.Error())
	}
	return nil
}

func (r *userRepository) UpdatePassword(user *entities.User) error {
	db := database.GetDBInstance()
	query := "UPDATE users SET password = ?, passwordRecoveryCode = NULL, recoveryCodeExpiration = NULL WHERE id = ?"
	_, err := db.Exec(query, user.Password, user.ID)
	if err != nil {
		return errors.NewQueryError(err.Error())
	}
	return nil
}
