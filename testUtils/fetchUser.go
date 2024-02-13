// In file: testUtils/utils.go

package testUtils

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"log"
)

func GetTestUserFromDB() (models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", "john.doe@example.com").First(&user).Error; err != nil {
		log.Printf("Failed to fetch test user: %v", err)
		return models.User{}, err
	}
	return user, nil
}
