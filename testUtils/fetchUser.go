// In file: testUtils/utils.go

package testUtils

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"

	"github.com/rs/zerolog/log"
)

func GetTestUserFromDB() (models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", "john.doe@example.com").First(&user).Error; err != nil {
		log.Error().Err(err).Msg("Failed to fetch test user")
		return models.User{}, err
	}
	return user, nil
}
