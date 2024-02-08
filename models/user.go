package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key;"`
	Username       string    `gorm:"unique;not null"`
	Password       string
	FirstName      string
	LastName       string
	AccountCreated time.Time `gorm:"autoCreateTime"`
	AccountUpdated time.Time `gorm:"autoUpdateTime"`
}
