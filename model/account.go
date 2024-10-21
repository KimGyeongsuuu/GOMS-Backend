package model

import (
	"GOMS-BACKEND-GO/model/data/constant"
	"context"
	"time"
)

type Account struct {
	ID         uint64
	Email      string
	Password   string
	Grade      int
	Name       string
	Gender     constant.Gender
	Major      constant.Major
	ProfileURL *string
	Authority  constant.Authority
	CreatedAt  time.Time
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *Account) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*Account, error)
}
