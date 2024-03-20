package middleware

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/utils"
	"encoding/base64"
	"net/http"
	"strings"

	"cloud-proj/health-check/logs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func BasicAuth() gin.HandlerFunc {

	logger := logs.CreateLogger()

	defer logger.Sync()

	// sl := slog.New(zapslog.NewHandler(logger.Core(), nil))

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error("Authorization header is required",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("header", authHeader),
				zap.Int("status", 401),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Extract credentials
		payload := strings.SplitN(authHeader, " ", 2)
		if len(payload) != 2 || payload[0] != "Basic" {
			logger.Error("Invalid Authorization header format",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("header", authHeader),
				zap.Int("status", 401),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(payload[1])
		if err != nil {
			logger.Error("Invalid Authorization header encoding",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", 401),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header encoding"})
			return
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			logger.Error("Invalid authentication credential format",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", 401),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization credential format"})
			return
		}

		username, password := parts[0], parts[1]

		// Verify the username and password
		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			logger.Error("Password authentication failed",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", 401),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed. Please check the username and password"})
			return
		}

		if !utils.CheckPasswordHash(password, user.Password) {
			logger.Error("Password authentication failed",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", 401),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed. Please check the username and password"})
			return
		}

		// Proceed if authentication is successful
		c.Set("user", user)
		c.Next()
	}
}
