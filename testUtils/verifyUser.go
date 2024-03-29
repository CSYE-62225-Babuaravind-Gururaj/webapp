package testUtils

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/router"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func VerifyUserByEmailToken(token string) error {
	database.InitDB()

	// Use the router setup for handling requests
	r := router.RouterSetup(database.DB)

	// Call the verify endpoint with the token from the created verification entry
	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/user/verify?token=%s", token), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert on the response status code
	if w.Code != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", w.Code)
	}

	// Fetch the updated verification entry from the database
	var verification models.VerifyUser
	if err := database.DB.Where("id = ?", token).First(&verification).Error; err != nil {
		return fmt.Errorf("failed to fetch verification entry: %w", err)
	}

	// Assert that the EmailVerified field is true
	if !verification.EmailVerified {
		return fmt.Errorf("email was not verified successfully")
	}

	return nil
}
