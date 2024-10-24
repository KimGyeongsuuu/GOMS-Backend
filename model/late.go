package model

import (
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"time"
)

type Late struct {
	LateID    uint64
	AccountID uint64
	Account   *Account `gorm:"foreignKey:AccountID"`
	CreatedAt time.Time
}
type LateUseCase interface {
	GetTop3LateStudent(ctx context.Context) ([]output.LateTop3Output, error)
}

type LateRepository interface {
	FindTop3ByOrderByAccountDesc(ctx context.Context) ([]Late, error)
}
