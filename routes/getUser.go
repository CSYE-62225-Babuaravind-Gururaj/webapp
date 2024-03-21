package routes

import (
	"cloud-proj/health-check/logs"
	"cloud-proj/health-check/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetUserRoute(c *gin.Context) {
	logger := logs.CreateLogger()

	// Note: No need for logger.Sync() with zerolog as it writes directly.

	logger.Info().Msg("Entered GetUserSelf handler")

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	authUser, ok := user.(models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Not Found"})
		return
	}

	response := struct {
		ID             uuid.UUID `json:"id"`
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		Username       string    `json:"username"`
		AccountCreated string    `json:"account_created"`
		AccountUpdated string    `json:"account_updated"`
	}{
		ID:             authUser.ID,
		FirstName:      authUser.FirstName,
		LastName:       authUser.LastName,
		Username:       authUser.Username,
		AccountCreated: authUser.AccountCreated.Format(time.RFC3339Nano),
		AccountUpdated: authUser.AccountUpdated.Format(time.RFC3339Nano),
	}

	c.JSON(http.StatusOK, response)

	// Log the incoming request with zerolog.
	logger.Info().
		Str("method", "GET").
		Str("path", "/v1/user/self").
		Int("status", http.StatusOK).
		Msg("incoming request")
}
