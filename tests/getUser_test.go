package tests

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/router"
	"cloud-proj/health-check/testUtils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetUser(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	database.InitDB()

	// Create a test user, _ here is used to ignore second variable returned as we don't need it
	testUser, _, err := testUtils.CreateTestUser()

	var user models.User
	err = database.DB.Order("id DESC").First(&user).Error
	if err != nil {
		t.Fatalf("Failed to fetch an existing user: %v", err)
	}

	// Fetch the associated VerifyUser entry to get the token (ID in this case).
	//var verifyUser models.VerifyUser
	err = database.DB.Where("username = ?", user.Username).First(&user).Error
	if err != nil {
		t.Fatalf("Failed to fetch VerifyUser entry for the user: %v", err)
	}

	// Set up your application router
	router := router.RouterSetup(database.DB)

	err = testUtils.VerifyUserByEmailToken(fmt.Sprint(user.ID))
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	}

	req, _ := http.NewRequest("GET", "/v1/user/self", nil)

	// Basic Auth header
	req.SetBasicAuth(testUser.Username, "Abcd1234")

	// Execute the request using a test server
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert on the response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

}
