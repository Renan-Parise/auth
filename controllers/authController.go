package controllers

import (
	"net/http"

	"github.com/Renan-Parise/codium/entities"
	"github.com/Renan-Parise/codium/services"
	"github.com/Renan-Parise/codium/utils"
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
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in controller method Login.")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ac.authService.Login(credentials.Email, credentials.Password)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to login in controller method Login.")

		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ac *AuthController) Register(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in controller method Register.")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authService.Register(user)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to register in controller method Register.")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func (ac *AuthController) Update(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to bind JSON in controller method Update.")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authService.Update(user)
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to update in controller method Update.")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update successful"})
}
