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
		authRoutes.POST("/fa/confirm", authController.ConfirmTwoFA)
		authRoutes.POST("/password/recover", authController.InitiatePasswordRecovery)
		authRoutes.POST("/password/reset", authController.ResetPassword)

		authRoutes.PUT("/update", middlewares.AuthMiddleware(), authController.Update)
		authRoutes.DELETE("/deactivate", middlewares.AuthMiddleware(), authController.Deactivate)
		authRoutes.POST("/fa/toggle", middlewares.AuthMiddleware(), authController.ToggleTwoFA)
		authRoutes.POST("/fa/confirm-toggle", middlewares.AuthMiddleware(), authController.ConfirmToggleTwoFA)
	}

	pingController := controllers.NewPingController()
	router.GET("/ping", pingController.Ping)

	return router
}
