package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"music-library/internal/config"
	"music-library/internal/models/song"
	"music-library/pkg/utils/logger"
)

var DB *gorm.DB

func DBConnect() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		config.Cfg.Database.Host,
		config.Cfg.Database.User,
		config.Cfg.Database.Password,
		config.Cfg.Database.Name,
		config.Cfg.Database.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	log.Print("DB setup successfully")

	if err = DB.AutoMigrate(&song.Song{}); err != nil {
		logger.Logger.Error(fmt.Sprintf("failed to migrate model: %v", err))
	}

	log.Print("Model loaded successfully")
	return
}
