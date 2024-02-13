package tests

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/router"
	"cloud-proj/health-check/testUtils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetUser(t *testing.T) {

	os.Setenv("DBHOST", "localhost")
	os.Setenv("DBPORT", "5432")
	os.Setenv("DBUSER", "postgres")
	os.Setenv("DBPASS", "root")
	os.Setenv("DBNAME", "userdb")

	database.InitDB()

	// Create a test user
	testUser := testUtils.CreateTestUser()

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
