package models

import (
	"time"

	"gorm.io/gorm"
)

type VerifyUser struct {
	gorm.Model
	ID               string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username         string `gorm:"type:varchar(100);uniqueIndex"`
	Token            string `gorm:"type:varchar(100);uniqueIndex"`
	EmailTriggerTime time.Time
	EmailVerified    bool `gorm:"default:false"`
}
