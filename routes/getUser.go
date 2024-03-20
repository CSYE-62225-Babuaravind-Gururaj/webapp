package routes

import (
	"cloud-proj/health-check/logs"
	"cloud-proj/health-check/models"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap/exp/zapslog"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetUserRoute(c *gin.Context) {

	logger := logs.CreateLogger()

	defer logger.Sync()

	sl := slog.New(zapslog.NewHandler(logger.Core(), nil))

	log.Println("Entered GetUserSelf handler")

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
	sl.Info(
		"incoming request",
		slog.String("method", "GET"),
		slog.String("path", "/v1/user/self"),
		slog.Int("status", 200),
	)
}
