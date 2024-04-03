package routes

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func VerifyUserRoute(c *gin.Context) {
	// var user models.User
	// err := database.DB.Order("id DESC").First(&user).Error
	// if err != nil {
	// 	log.Printf("Failed to fetch an existing user: %v", err)
	// }

	token := c.Query("token")

	log.Printf("Token print: %s", c.Query("token"))

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token query parameter is required"})
		return
	}

	// Fetch the associated VerifyUser entry to get the token
	var verifyUser models.VerifyUser
	err := database.DB.Where("token = ?", token).First(&verifyUser).Error
	if err != nil {
		log.Printf("Failed to fetch VerifyUser entry for the user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Expiry check
	twoMinutesFromNow := verifyUser.EmailTriggerTime.Add(2 * time.Minute)
	if time.Now().After(twoMinutesFromNow) {
		c.JSON(http.StatusGone, gin.H{"error": "Verification link has expired"})
		return
	}
	// if verifyUser.EmailTriggerTime.After(twoMinutesFromNow) {
	// 	c.JSON(http.StatusGone, gin.H{"error": "Verification link has expired"})
	// 	return
	// }

	// Email already verified
	if verifyUser.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already verified"})
		return
	}

	// Verification
	verifyUser.EmailVerified = true
	if err := database.DB.Save(&verifyUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification status", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
