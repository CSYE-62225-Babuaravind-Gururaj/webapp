package tests

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/router"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyUser(t *testing.T) {
	database.InitDB()

	var user models.User
	err := database.DB.Order("id DESC").First(&user).Error
	if err != nil {
		t.Fatalf("Failed to fetch an existing user: %v", err)
	}

	// Fetch the associated VerifyUser entry to get the token (ID in this case).
	var verifyUser models.VerifyUser
	err = database.DB.Where("username = ?", user.Username).First(&verifyUser).Error
	if err != nil {
		t.Fatalf("Failed to fetch VerifyUser entry for the user: %v", err)
	}

	// Use the router setup for handling requests
	r := router.RouterSetup(database.DB)

	// Call the verify endpoint with the token from the created verification entry
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/user/verify?token=%s", verifyUser.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert on the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Fetch the updated verification entry from the database
	var verification models.VerifyUser
	if err := database.DB.Where("id = ?", verifyUser.ID).First(&verification).Error; err != nil {
		t.Fatalf("Failed to fetch verification entry: %v", err)
	}

	// Assert that the EmailVerified field is true
	if !verification.EmailVerified {
		t.Errorf("Email was not verified successfully")
	}
}
