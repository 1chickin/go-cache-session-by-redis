package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	SessionToken string `gorm:"unique;not null"`
	UserID       uint
	ExpiresAt    time.Time
}
