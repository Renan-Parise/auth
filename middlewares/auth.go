package middlewares

import (
	"net/http"
	"strings"

	"github.com/Renan-Parise/codium/repositories"
	"github.com/Renan-Parise/codium/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authentication token"})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token format"})
			return
		}

		tokenString := tokenParts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			return
		}

		IDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			return
		}
		ID := int(IDFloat)

		userRepo := repositories.NewUserRepository()
		user, err := userRepo.FindByID(ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": "user not found"})
			return
		}
		if !user.Active {
			c.AbortWithStatusJSON(http.StatusLocked, gin.H{"error": "user account is deactivated"})
			return
		}

		c.Set("ID", ID)
		c.Next()
	}
}
