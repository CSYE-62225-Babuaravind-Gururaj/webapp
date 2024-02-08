package database

import (
	"cloud-proj/health-check/models"
	"fmt"
	"log"
	"os"

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
		log.Fatal("Failed to connect to DB: ", err)
	}

	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatalf("Failed to create extension 'uuid-ossp': %v", err)
	}

	log.Println("mirgating")
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Panic("Failed to auto-migrate:", err)
	}

	log.Println("Auto-migration done")
}
