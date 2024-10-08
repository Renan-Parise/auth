package routes

import (
	"github.com/Renan-Parise/codium-auth/controllers"
	"github.com/Renan-Parise/codium-auth/middlewares"
	"github.com/Renan-Parise/codium-auth/repositories"
	"github.com/Renan-Parise/codium-auth/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	userRepo := repositories.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
		authRoutes.PUT("/update", middlewares.AuthMiddleware(), authController.Update)
		authRoutes.DELETE("/deactivate", middlewares.AuthMiddleware(), authController.Deactivate)
	}

	pingController := controllers.NewPingController()
	router.GET("/ping", pingController.Ping)

	return router
}
