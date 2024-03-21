package database

import (
	"cloud-proj/health-check/models"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"),
	)
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		//Inbuilt Logging capacity provide by GORM. Use this to temporarily log the information
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DB")
	}

	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create extension 'uuid-ossp'")
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to auto-migrate")
	}

}
