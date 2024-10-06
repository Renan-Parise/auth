package routes

import (
	"github.com/Renan-Parise/codium/controllers"
	"github.com/Renan-Parise/codium/middlewares"
	"github.com/Renan-Parise/codium/repositories"
	"github.com/Renan-Parise/codium/services"

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
		authRoutes.PUT("/update", authController.Update)
		authRoutes.DELETE("/deactivate", middlewares.AuthMiddleware(), authController.Deactivate)
	}

	pingController := controllers.NewPingController()
	router.GET("/ping", pingController.Ping)

	return router
}
