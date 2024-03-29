package testUtils

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/utils"
	"log"
	"time"
)

func CreateTestUser() (models.User, models.VerifyUser, error) {
	database.InitDB()

	pass, err := utils.HashPassword("Abcd1234")
	if err != nil {
		log.Printf("Failed to hash password: %v", err)

	}

	user := models.User{
		FirstName:      "John",
		LastName:       "Doe",
		Username:       "john.doe@example.com",
		Password:       pass,
		AccountCreated: time.Now(),
		AccountUpdated: time.Now(),
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		log.Printf("Failed to create test user: %v", result.Error)
	}

	verifyUser := models.VerifyUser{
		Username:         user.Username,
		EmailTriggerTime: time.Now(),
		EmailVerified:    false,
		Token:            user.ID.String(),
	}
	verifyResult := database.DB.Create(&verifyUser)
	if verifyResult.Error != nil {
		log.Printf("Failed to create test user: %v", verifyResult.Error)
	}

	return user, verifyUser, nil
}
