package tests

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/router"
	"cloud-proj/health-check/testUtils"
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

	// Set up your application router
	router := router.RouterSetup(database.DB)

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
