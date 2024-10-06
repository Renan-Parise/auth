package main

import (
	"log"

	"github.com/Renan-Parise/codium/database"
	"github.com/Renan-Parise/codium/routes"
	"github.com/Renan-Parise/codium/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. is it missing?")
	}

	utils.InitLogger()
	utils.InitElasticAPM()

	database.GetDBInstance()
	router := routes.SetupRouter()

	router.Run("127.0.0.1:8181")
}
