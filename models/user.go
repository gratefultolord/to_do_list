package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey" example:"1"`
	Username  string         `json:"username" binding:"required" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" binding:"required,email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"password" binding:"required" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" swaggerignore:"true"`
}
