package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	VkApiToken          string
	TelegramBotApiToken string
	DBPath              string
	BlevePath           string
}

func InitWithDotEnv() *AppConfig {
	// Пробуем загрузить .env, но не падаем при ошибке
	_ = godotenv.Load("../.env")

	conf := new(AppConfig)
	conf.VkApiToken = os.Getenv("VK_API_TOKEN")
	conf.TelegramBotApiToken = os.Getenv("TOKEN")

	// Пути к файлам данных (поддержка Docker)
	// По умолчанию используем локальные пути
	conf.DBPath = os.Getenv("DB_PATH")
	conf.BlevePath = os.Getenv("BLEVE_PATH")

	if conf.DBPath == "" {
		conf.DBPath = "../glavredus.db"
	}

	if conf.BlevePath == "" {
		conf.BlevePath = "../history.bleve"
	}

	// Если токен пустой, используем заглушку для тестирования
	if conf.VkApiToken == "" {
		log.Println("Warning: VK_API_TOKEN not set, using placeholder")
		conf.VkApiToken = "placeholder_token"
	}

	return conf
}

// GetAbsoluteDBPath возвращает абсолютный путь к базе данных
func (c *AppConfig) GetAbsoluteDBPath() (string, error) {
	if filepath.IsAbs(c.DBPath) {
		return c.DBPath, nil
	}
	return filepath.Abs(c.DBPath)
}

// GetAbsoluteBlevePath возвращает абсолютный путь к индексу Bleve
func (c *AppConfig) GetAbsoluteBlevePath() (string, error) {
	if filepath.IsAbs(c.BlevePath) {
		return c.BlevePath, nil
	}
	return filepath.Abs(c.BlevePath)
}
