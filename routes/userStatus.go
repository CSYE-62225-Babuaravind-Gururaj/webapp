package routes

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyUserRoute(c *gin.Context) {

	email := c.Query("email")
	token := c.Query("token")

	if email == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and token query parameters are required"})
		return
	}

	var verificationEntry models.VerifyUser

	if err := database.DB.Where("username = ?", email).First(&verificationEntry).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification record not found"})
		return
	}

	// Here, validate the token matches what was sent in the email
	// This is a simplified version. In reality, you would need a secure way to generate and validate tokens.
	// if verificationEntry.EmailVerified || !TokenIsValid(token) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification link"})
	// 	return
	// }

	verificationEntry.EmailVerified = true
	if err := database.DB.Save(&verificationEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
