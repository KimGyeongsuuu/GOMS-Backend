package model

import (
	"GOMS-BACKEND-GO/model/data/output"
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
	FindAllOutingStudent(ctx context.Context) ([]output.OutingStudentOutput, error)
	CountAllOutingStudent(ctx context.Context) (int, error)
	SearchOutingStudent(ctx context.Context, name string) ([]output.OutingStudentOutput, error)
}

type OutingRepository interface {
	SaveOutingStudnet(ctx context.Context, outing *Outing) error
	ExistsOutingByAccountID(ctx context.Context, accountID uint64) (bool, error)
	DeleteOutingByAccountID(ctx context.Context, accountID uint64) error
	FindAllOuting(ctx context.Context) ([]Outing, error)
	FindByOutingAccountNameContaining(ctx context.Context, name string) ([]Outing, error)
}
