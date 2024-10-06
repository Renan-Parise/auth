package entities

import (
	"regexp"
	"time"

	"github.com/Renan-Parise/codium/errors"
)

type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Active        bool      `json:"active"`
	DeactivatedAt time.Time `json:"deactivatedAt"`
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
