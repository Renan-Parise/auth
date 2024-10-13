package services

import (
	"time"

	"github.com/Renan-Parise/codium-auth/entities"
	"github.com/Renan-Parise/codium-auth/errors"
	"github.com/Renan-Parise/codium-auth/repositories"
	"github.com/Renan-Parise/codium-auth/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (string, error)
	Register(user entities.User) error
	Update(ID int, user entities.User) error
	DeactivateAccount(ID int) error
	GenerateAndSendTwoFACode(user *entities.User) error
	VerifyTwoFACode(email, code string) (string, error)
	GenerateAndSendTwoFACodeByID(userID int) error
	ToggleTwoFA(userID int, code string) error
	InitiatePasswordRecovery(email string) error
	ResetPassword(email, code, newPassword string) error
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{userRepo: repo}
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.NewServiceError("authentication failed because user does not exist")
	}

	if !user.Active {
		return "", errors.NewServiceError("authentication failed because account is deactivated")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.NewServiceError("authentication failed because password is incorrect")
	}

	if user.Is2FAEnabled {
		err := s.GenerateAndSendTwoFACode(user)
		if err != nil {
			return "", errors.NewServiceError("failed to send 2FA code")
		}
		return "", entities.ErrTwoFARequired
	}

	return utils.GenerateToken(user.ID)
}

func (s *authService) Register(user entities.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	_, err := s.userRepo.FindByEmail(user.Email)
	if err == nil {
		return errors.NewServiceError("user already exists. please login or use another email")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewServiceError("failed to hash password. please try again")
	}

	user.Password = string(hashedPassword)

	err = s.userRepo.Create(user)
	if err != nil {
		return errors.NewServiceError("failed to register user. please try again")
	}

	return nil
}

func (s *authService) Update(ID int, user entities.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewServiceError("failed to hash password. please try again")
	}

	user.Password = string(hashedPassword)

	err = s.userRepo.Update(ID, user)
	if err != nil {
		return errors.NewServiceError("failed to update user. please try again")
	}

	return nil
}

func (s *authService) DeactivateAccount(ID int) error {
	err := s.userRepo.DeactivateUser(ID)
	if err != nil {
		return errors.NewServiceError("failed to deactivate account")
	}
	return nil
}

func (s *authService) GenerateAndSendTwoFACode(user *entities.User) error {
	code := utils.GenerateCode(6)

	user.TwoFACode = &code
	expirationTime := time.Now().Add(5 * time.Minute)
	user.TwoFACodeExpiresAt = &expirationTime

	err := s.userRepo.UpdateTwoFACode(user)
	if err != nil {
		return err
	}

	err = s.sendTwoFACodeEmail(user.Email, code)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) VerifyTwoFACode(email string, code string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.NewServiceError("user not found")
	}

	if *user.TwoFACode != code || time.Now().After(*user.TwoFACodeExpiresAt) {
		return "", errors.NewServiceError("invalid or expired 2FA code")
	}

	user.TwoFACode = nil
	user.TwoFACodeExpiresAt = nil

	err = s.userRepo.UpdateTwoFACode(user)
	if err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID)
}

func (s *authService) GenerateAndSendTwoFACodeByID(userID int) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.NewServiceError("user not found")
	}

	return s.GenerateAndSendTwoFACode(user)
}

func (s *authService) ToggleTwoFA(userID int, code string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.NewServiceError("user not found")
	}

	if *user.TwoFACode != code || time.Now().After(*user.TwoFACodeExpiresAt) {
		return errors.NewServiceError("invalid or expired 2FA code")
	}

	if user.Is2FAEnabled {
		user.Is2FAEnabled = false
	} else {
		user.Is2FAEnabled = true
	}

	err = s.userRepo.UpdateTwoFASettings(user)
	if err != nil {
		return errors.NewServiceError("failed to update 2FA settings")
	}

	return nil
}

func (s *authService) InitiatePasswordRecovery(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.NewServiceError("user not found")
	}

	code := utils.GenerateCode(6)
	user.PasswordRecoveryCode = &code
	expirationTime := time.Now().Add(30 * time.Minute)
	user.RecoveryCodeExpiresAt = &expirationTime

	err = s.userRepo.UpdatePasswordRecoveryCode(user)
	if err != nil {
		return errors.NewServiceError("failed to save recovery code")
	}

	err = s.sendPasswordRecoveryEmail(user.Email, code)
	if err != nil {
		return errors.NewServiceError("failed to send recovery email")
	}

	return nil
}

func (s *authService) ResetPassword(email, code, newPassword string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.NewServiceError("user not found")
	}

	if user.PasswordRecoveryCode == nil || *user.PasswordRecoveryCode != code || time.Now().After(*user.RecoveryCodeExpiresAt) {
		return errors.NewServiceError("invalid or expired recovery code")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewServiceError("failed to hash new password")
	}

	user.Password = string(hashedPassword)
	user.PasswordRecoveryCode = nil
	user.RecoveryCodeExpiresAt = nil

	err = s.userRepo.UpdatePassword(user)
	if err != nil {
		return errors.NewServiceError("failed to update password")
	}

	return nil
}
