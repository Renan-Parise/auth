package entities

import (
	"regexp"

	"github.com/Renan-Parise/codium-auth/errors"
)

type Email struct {
	Address string `json:"address"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (e *Email) Validate() error {
	if e.Address == "" {
		return errors.NewValidationError("address", "address is required. please provide a valid address")
	}

	if regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(e.Address) {
		return errors.NewValidationError("address", "address is invalid. please provide a valid address")
	}

	if e.Subject == "" {
		return errors.NewValidationError("subject", "subject is required. please provide a valid subject")
	}

	if e.Body == "" {
		return errors.NewValidationError("body", "body is required. please provide a valid body")
	}

	return nil
}
