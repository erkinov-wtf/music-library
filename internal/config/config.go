package config

import (
	"log"
	"os"
)

var Cfg Config

type Config struct {
	General  General
	Database Database
}

type General struct {
	Env  string
	Port string
}

type Database struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Timezone string
}

func MustLoad() {
	Cfg = Config{
		General: General{
			Env:  getEnv("ENV"),
			Port: getEnv("PORT"),
		},
		Database: Database{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			Name:     getEnv("DB_NAME"),
			User:     getEnv("DB_USERNAME"),
			Password: getEnv("DB_PASSWORD"),
		},
	}

	log.Printf("Configuration loaded: %v", Cfg)
}

func getEnv(key string) string {
	return os.Getenv(key)
}
