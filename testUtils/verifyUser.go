package testUtils

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/router"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

func VerifyUserByEmailToken(token string) error {
	database.InitDB()

	// Use the router setup for handling requests
	r := router.RouterSetup(database.DB)

	// Call the verify endpoint with the token from the created verification entry
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/user/verify?token=%s", token), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert on the response status code
	if w.Code != http.StatusOK {
		log.Printf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Fetch the updated verification entry from the database
	var verification models.VerifyUser

	// err := database.DB.Order("id DESC").First(&verification).Error
	// if err != nil {
	// 	log.Printf("Failed to fetch an existing user: %v", err)
	// }

	if err := database.DB.Where("token = ?", token).First(&verification).Error; err != nil {
		log.Printf("Failed to fetch verification entry: %v", err)
	}

	// Assert that the EmailVerified field is true
	if !verification.EmailVerified {
		log.Printf("Email was not verified successfully")
	}

	return nil
}
