package routes

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/logs"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUserRoute(c *gin.Context) {
	logger := logs.CreateLogger()

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUser, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	//declare and fetch verifyUser data
	var userVerify models.VerifyUser
	if err := database.DB.Where("username = ?", authUser.Username).First(&userVerify).Error; err != nil {
		// If there's an error fetching the record, it could mean the record doesn't exist
		logger.Error().Err(err).Msg("Failed to fetch verification status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user status"})
		return
	}

	if !userVerify.EmailVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Email not verified"})
		return
	}

	var updateReq models.UpdateUser
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if updateReq.FirstName == "" && updateReq.LastName == "" && updateReq.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields entered"})
		return
	}

	if !utils.ValidateName(updateReq.FirstName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a valid FirstName"})
		return
	} else {
		authUser.FirstName = updateReq.FirstName
	}

	if !utils.ValidateName(updateReq.LastName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a valid LastName"})
		return
	} else {
		authUser.LastName = updateReq.LastName
	}

	if updateReq.UserName != authUser.Username {
		c.Status(http.StatusUnauthorized)
	} else {
		if !utils.ValidatePassword(updateReq.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a valid Password"})
			return
		} else {
			hashedPassword, err := utils.HashPassword(updateReq.Password)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to update user"})
				return
			}
			authUser.Password = hashedPassword
		}

		// Saving to DB
		result := database.DB.Save(&authUser)
		if result.Error != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to update user"})
			return
		}

		c.Status(http.StatusNoContent)

		// Log the successful user update
		logger.Info().
			Str("method", "PUT").
			Str("path", "/v1/user/self").
			Int("status", http.StatusNoContent).
			Msg("User updated successfully")
	}
}
