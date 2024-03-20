package routes

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/logs"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/utils"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/exp/zapslog"
)

func UpdateUserRoute(c *gin.Context) {
	logger := logs.CreateLogger()

	defer logger.Sync()

	sl := slog.New(zapslog.NewHandler(logger.Core(), nil))

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
		sl.Info(
			"incoming request",
			slog.String("method", "PUT"),
			slog.String("path", "/v1/user/self"),
			slog.Int("status", 204),
		)
	}

}
