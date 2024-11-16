package services

import (
	"fmt"

	"github.com/Renan-Parise/auth/entities"
	"github.com/Renan-Parise/auth/utils"
)

func (s *authService) sendPasswordRecoveryEmail(email, code string) error {
	emailEntity := entities.Email{
		Address: email,
		Subject: "Password Recovery",
		Body:    fmt.Sprintf("Your password recovery code is: %s", code),
	}

	err := utils.SendEmail(emailEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) sendTwoFACodeEmail(email, code string) error {
	emailEntity := entities.Email{
		Address: email,
		Subject: "Your Two-Factor Authentication Code",
		Body:    fmt.Sprintf("Your Two-Factor Authentication code is: %s", code),
	}

	err := utils.SendEmail(emailEntity)
	if err != nil {
		return err
	}

	return nil
}
