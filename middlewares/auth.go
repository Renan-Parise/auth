package middlewares

import (
	"net/http"

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

		_, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.GetLogger().WithError(err).Error("Failed to validate token in middleware AuthMiddleware.")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			return
		}

		c.Next()
	}
}
