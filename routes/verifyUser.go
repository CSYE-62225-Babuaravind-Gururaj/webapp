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

	userName := "john.doe@example.com"

	// Fetch the associated VerifyUser entry to get the token (ID in this case).
	var verifyUser models.VerifyUser
	err := database.DB.Where("username = ?", userName).First(&verifyUser).Error
	if err != nil {
		log.Printf("Failed to fetch VerifyUser entry for the user: %v", err)
	}

	if verifyUser.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token query parameter is required"})
		return
	}

	var verificationEntry models.VerifyUser

	if err := database.DB.Where("id = ?", verifyUser.ID).First(&verificationEntry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Verification record not found"})
		return
	}

	// Expiry check
	if time.Since(verificationEntry.EmailTriggerTime) > 2*time.Minute {
		c.JSON(http.StatusGone, gin.H{"error": "Verification link has expired"})
		return
	}

	// Email already verified
	if verificationEntry.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already verified"})
		return
	}

	// Verification
	verificationEntry.EmailVerified = true
	if err := database.DB.Save(&verificationEntry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
