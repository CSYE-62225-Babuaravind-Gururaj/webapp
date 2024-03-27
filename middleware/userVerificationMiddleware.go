package middleware

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		authUser, ok := user.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user data"})
			return
		}

		var verification models.VerifyUser
		if err := database.DB.Where("username = ?", authUser.Username).First(&verification).Error; err != nil || !verification.EmailVerified {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Email not verified"})
			return
		}

		c.Next()
	}
}
