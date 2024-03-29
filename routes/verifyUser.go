package routes

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VerifyUserRoute(c *gin.Context) {
	token := c.Query("token")

	log.Printf("Token print: %s", token)

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token query parameter is required"})
		return
	}

	var verifyUser models.VerifyUser
	err := database.DB.Where("id = ?", token).First(&verifyUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("VerifyUser entry not found for token: %s", token)
			c.JSON(http.StatusNotFound, gin.H{"error": "Verification record not found"})
			return
		}
		log.Printf("Failed to fetch VerifyUser entry for the token: %s, Error: %v", token, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch verification record"})
		return
	}

	twoMinutesFromNow := time.Now().Add(-2 * time.Minute) // Subtracting time for comparison with past
	if verifyUser.EmailTriggerTime.Before(twoMinutesFromNow) {
		c.JSON(http.StatusGone, gin.H{"error": "Verification link has expired"})
		return
	}

	if verifyUser.EmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already verified"})
		return
	}

	verifyUser.EmailVerified = true
	if err := database.DB.Save(&verifyUser).Error; err != nil {
		log.Printf("Failed to update verification status for token: %s, Error: %v", token, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verification status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
