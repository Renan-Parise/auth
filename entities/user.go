package entities

import (
	"regexp"
	"time"

	"github.com/Renan-Parise/codium-auth/errors"
)

var ErrTwoFARequired = errors.NewServiceError("2FA required")

type User struct {
	ID                    int        `json:"id"`
	Username              string     `json:"username"`
	Email                 string     `json:"email"`
	Password              string     `json:"password"`
	Active                bool       `json:"active"`
	DeactivatedAt         *time.Time `json:"deactivatedAt"`
	Is2FAEnabled          bool       `json:"is2FAEnabled"`
	TwoFACode             *string    `json:"-"`
	TwoFACodeExpiresAt    *time.Time `json:"-"`
	PasswordRecoveryCode  *string    `json:"-"`
	RecoveryCodeExpiresAt *time.Time `json:"-"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return errors.NewValidationError("username", "username is required. please provide a valid username")
	}

	if u.Email == "" {
		return errors.NewValidationError("email", "email is required. please provide a valid email")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(u.Email) {
		return errors.NewValidationError("email", "email is invalid. please provide a valid email")
	}

	if u.Password == "" {
		return errors.NewValidationError("password", "password is required. please provide a valid password")
	}

	return nil
}
