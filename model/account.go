package model

import (
	"time"
)

type Account struct {
	Email      string    `gorm:"type:varchar(40);not null"`
	Password   string    `gorm:"not null"`
	Grade      int       `gorm:"not null"`
	Name       string    `gorm:"type:varchar(10);not null"`
	Gender     string    `gorm:"type:varchar(10);not null"`
	Major      string    `gorm:"type:varchar(20);not null"`
	ProfileURL *string   `gorm:"type:text"`
	Authority  string    `gorm:"type:varchar(20);not null"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
