package model

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Outing struct {
	ID        uint64
	AccountID uint64 `gorm:"not null;index"`
	Account   *Account
	CreatedAt time.Time
}
type OutingUseCase interface {
	OutingStudent(c *gin.Context, ctx context.Context, outingUUID uuid.UUID) error
}

type OutingRepository interface {
	SaveOutingStudnet(ctx context.Context, outing *Outing) error
	ExistsOutingByAccountID(ctx context.Context, accountID uint64) (bool, error)
	DeleteOutingByAccountID(ctx context.Context, accountID uint64) error
}
