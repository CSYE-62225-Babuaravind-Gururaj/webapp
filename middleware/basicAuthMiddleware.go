package middleware

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/logs"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/utils"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
	logger := logs.CreateLogger()

	// Note: No need for logger.Sync() with zerolog as it writes directly.

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error().
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Str("header", authHeader).
				Int("status", http.StatusUnauthorized).
				Msg("Authorization header is required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Extract credentials
		payload := strings.SplitN(authHeader, " ", 2)
		if len(payload) != 2 || payload[0] != "Basic" {
			logger.Error().
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Str("header", authHeader).
				Int("status", http.StatusUnauthorized).
				Msg("Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(payload[1])
		if err != nil {
			logger.Error().
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Int("status", http.StatusUnauthorized).
				Msg("Invalid Authorization header encoding")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header encoding"})
			return
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			logger.Error().
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Int("status", http.StatusUnauthorized).
				Msg("Invalid authentication credential format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization credential format"})
			return
		}

		username, password := parts[0], parts[1]

		// Verify the username and password
		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			logger.Error().
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Int("status", http.StatusUnauthorized).
				Msg("Authentication failed. Please check the username and password")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed. Please check the username and password"})
			return
		}

		if !utils.CheckPasswordHash(password, user.Password) {
			logger.Error().
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Int("status", http.StatusUnauthorized).
				Msg("Password authentication failed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed. Please check the username and password"})
			return
		}

		// Proceed if authentication is successful
		c.Set("user", user)
		c.Next()
	}
}
