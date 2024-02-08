package routes

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUserRoute(c *gin.Context) {
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

	var updateReq models.UpdateUser
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if updateReq.UserName != authUser.Username {
		c.Status(http.StatusUnauthorized)
	} else {
		if updateReq.Password != "" {
			hashedPassword, err := utils.HashPassword(updateReq.Password)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to update user"})
				return
			}
			authUser.Password = hashedPassword
		}

		if updateReq.FirstName != "" {
			authUser.FirstName = updateReq.FirstName
		}

		if updateReq.LastName != "" {
			authUser.LastName = updateReq.LastName
		}

		// Saving to DB
		result := database.DB.Save(&authUser)
		if result.Error != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Failed to update user"})
			return
		}

		c.Status(http.StatusNoContent)
	}

}
