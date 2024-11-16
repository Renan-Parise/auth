package controllers

import (
	"fmt"
	"net/http"

	"github.com/Renan-Parise/auth/entities"
	"github.com/Renan-Parise/auth/services"
	"github.com/Renan-Parise/auth/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(service services.AuthService) *AuthController {
	return &AuthController{authService: service}
}

func (ac *AuthController) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in controller method Login: ", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ac.authService.Login(credentials.Email, credentials.Password)
	if err != nil {
		if err == entities.ErrTwoFARequired {
			c.JSON(http.StatusAccepted, gin.H{"message": "2FA code sent to email"})
			return
		}
		utils.GetLogger().WithError(err).Error("Failed to login in controller method Login: ", err)

		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ac *AuthController) Register(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in controller method Register: ", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authService.Register(user)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to register in controller method Register: ", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func (ac *AuthController) Update(c *gin.Context) {
	ID, exists := c.Get("ID")
	if !exists {
		utils.GetLogger().Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in controller method Update: ", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authService.Update(ID.(int), user)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to update in controller method Update: ", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update successful"})
}

func (ac *AuthController) Deactivate(c *gin.Context) {
	ID, exists := c.Get("ID")
	if !exists {
		utils.GetLogger().Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := ac.authService.DeactivateAccount(ID.(int))
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to deactivate account in controller method Deactivate: ", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account deactivated successfully"})
}

func (ac *AuthController) ConfirmTwoFA(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in controller method ConfirmTwoFA: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ac.authService.VerifyTwoFACode(request.Email, request.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ac *AuthController) ToggleTwoFA(c *gin.Context) {
	fmt.Println("ToggleTwoFA")
	ID, exists := c.Get("ID")
	if !exists {
		utils.GetLogger().Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := ac.authService.GenerateAndSendTwoFACodeByID(ID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send 2FA code"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "2FA code sent to email"})
}

func (ac *AuthController) ConfirmToggleTwoFA(c *gin.Context) {
	ID, exists := c.Get("ID")
	if !exists {
		utils.GetLogger().Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in controller method ConfirmToggleTwoFA: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authService.ToggleTwoFA(ID.(int), request.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA setting updated"})
}

func (ac *AuthController) InitiatePasswordRecovery(c *gin.Context) {
	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in InitiatePasswordRecovery")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := ac.authService.InitiatePasswordRecovery(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password recovery email sent"})
}

func (ac *AuthController) ResetPassword(c *gin.Context) {
	var request struct {
		Email       string `json:"email"`
		Code        string `json:"code"`
		NewPassword string `json:"newPassword"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in ResetPassword")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := ac.authService.ResetPassword(request.Email, request.Code, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password has been reset successfully"})
}
