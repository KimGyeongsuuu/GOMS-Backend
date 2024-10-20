package model

import (
	"context"
	"time"
)

type Account struct {
	Email      string    `gorm:"type:varchar(40);not null"`
	Password   string    `gorm:"not null"`
	Grade      int       `gorm:"not null"`
	Name       string    `gorm:"type:varchar(10);not null"`
	Gender     Gender    `gorm:"type:varchar(10);not null"`
	Major      Major     `gorm:"type:varchar(20);not null"`
	ProfileURL *string   `gorm:"type:text"`
	Authority  Authority `gorm:"type:varchar(20);not null"`
	CreatedAt  time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *Account) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type Authority string

const (
	ROLE_STUDENT         Authority = "ROLE_STUDENT"
	ROLE_STUDENT_COUNCIL Authority = "ROLE_STUDENT_COUNCIL"
)

type Gender string

const (
	MAN   Gender = "MAN"
	WOMAN Gender = "WOMAN"
)

type Major string

const (
	SW_DEVELOP Major = "SW_DEVELOP"
	SMART_IOT  Major = "SMART_IOT"
	AI         Major = "AI"
)
