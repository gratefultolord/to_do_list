package models

import (
	"gorm.io/gorm"
	"time"
)

type Status string

const (
	StatusStarted   Status = "начата"
	StatusCompleted Status = "завершена"
)

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"id""`
	Title       string `gorm:"not null"`
	Description string
	Status      Status `gorm:"default:'начата'"`
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index" swaggerignore:"true""`
}
