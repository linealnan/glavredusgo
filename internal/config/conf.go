package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	VkApiToken string
}

func InitWithDotEnv() *AppConfig {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := new(AppConfig)
	conf.VkApiToken = os.Getenv("VK_API_TOKEN")

	return conf
}
