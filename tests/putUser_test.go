package tests

import (
	"bytes"
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/router"
	"cloud-proj/health-check/testUtils"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func TestUpdateUser(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	database.InitDB()

	testUser, err := testUtils.GetTestUserFromDB()
	if err != nil {
		t.Fatalf("Failed to fetch test user: %v", err)
	}

	updateInfo := models.UpdateUser{
		FirstName: "Joe First Second",
		LastName:  "Doe",
		Password:  "Abcd1234",
		UserName:  "john.doe@example.com",
	}

	router := router.RouterSetup(database.DB)

	payloadBytes, _ := json.Marshal(updateInfo)
	payload := bytes.NewReader(payloadBytes)

	req, _ := http.NewRequest("PUT", "/v2/user/self", payload)
	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(testUser.Username, "Abcd1234")

	// Execute the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert on the response status code
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}

	var updatedUser models.User
	if err := database.DB.Where("username = ?", testUser.Username).First(&updatedUser).Error; err != nil {
		t.Fatalf("Failed to fetch updated user: %v", err)
	}

	if updatedUser.FirstName != updateInfo.FirstName || updatedUser.LastName != updateInfo.LastName {
		log.Println(updatedUser.FirstName, updatedUser.LastName)
		log.Println(updatedUser.Username)
		t.Errorf("User was not updated correctly")
	}
}
