package middleware

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/utils"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Extract credentials
		payload := strings.SplitN(authHeader, " ", 2)
		if len(payload) != 2 || payload[0] != "Basic" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(payload[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header encoding"})
			return
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization credential format"})
			return
		}

		username, password := parts[0], parts[1]

		// Verify the username and password
		var user models.User
		if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed. Please check the username and password"})
			return
		}

		if !utils.CheckPasswordHash(password, user.Password) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed. Please check the username and password"})
			return
		}

		// Proceed if authentication is successful
		c.Set("user", user)
		c.Next()
	}
}
