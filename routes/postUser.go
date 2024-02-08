package routes

import (
	"cloud-proj/health-check/database"
	"cloud-proj/health-check/models"
	"cloud-proj/health-check/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUserRoute(c *gin.Context) {
	var jsonMap map[string]interface{}

	if err := c.ShouldBindJSON(&jsonMap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	expectedFields := []string{"first_name", "last_name", "password", "username"}

	for key := range jsonMap {
		if !utils.ContainsString(expectedFields, key) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected field in request"})
			return
		}
	}

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
		Username  string `json:"username"`
	}

	//Manually testing the fields as GORM does not provide internal JSON field check
	input.FirstName, _ = jsonMap["first_name"].(string)
	input.LastName, _ = jsonMap["last_name"].(string)
	input.Password, _ = jsonMap["password"].(string)
	input.Username, _ = jsonMap["username"].(string)

	hashedPassword, err := utils.HashPassword(input.Password)

	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "Error hashing password"})
		return
	}

	if !utils.ValidateEmail(input.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter Valid Email"})
		return
	}

	if !utils.ValidateName(input.FirstName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter Valid First Name"})
		return
	}

	if !utils.ValidateName(input.LastName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter Valid Last Name"})
		return
	}

	if !utils.ValidatePassword(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password should have a minimum length of 8 characters and should contain at least one letter and one number"})
		return
	}

	user := models.User{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Username:       input.Username,
		Password:       hashedPassword,
		AccountCreated: time.Now(),
		AccountUpdated: time.Now(),
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		} else if result.Error == gorm.ErrPrimaryKeyRequired {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Primary Key"})
		} else if result.Error == gorm.ErrInvalidField {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Field"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":              user.ID,
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
		"username":        user.Username,
		"account_created": user.AccountCreated,
		"account_updated": user.AccountUpdated,
	})

}
