package middlewares

import (
	"net/http"

	"github.com/Renan-Parise/codium/repositories"
	"github.com/Renan-Parise/codium/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authentication token"})
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.GetLogger().WithError(err).Error("Invalid authentication token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			return
		}
		userID := int(userIDFloat)

		userRepo := repositories.NewUserRepository()
		user, err := userRepo.FindByID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		if !user.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user account is deactivated"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
