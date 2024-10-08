package main

import (
	"log"

	"github.com/Renan-Parise/codium-auth/database"
	"github.com/Renan-Parise/codium-auth/repositories"
	"github.com/Renan-Parise/codium-auth/routes"
	"github.com/Renan-Parise/codium-auth/utils"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. is it missing?")
	}

	utils.InitLogger()
	utils.InitElasticAPM()

	database.GetDBInstance()

	c := cron.New()
	_, err = c.AddFunc("@weekly", func() {
		userRepo := repositories.NewUserRepository()
		err := userRepo.DeleteInactiveUsers()
		if err != nil {
			utils.GetLogger().WithError(err).Error("Failed to delete inactive users in cron job: ", err)
		}
	})
	if err != nil {
		utils.GetLogger().WithError(err).Error("Failed to schedule cron job: ", err)
	}
	c.Start()
	defer c.Stop()

	router := routes.SetupRouter()
	router.Run("127.0.0.1:8181")
}
