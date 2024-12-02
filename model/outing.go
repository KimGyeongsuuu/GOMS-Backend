package model

import (
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Outing struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AccountID primitive.ObjectID `bson:"account_id"`
	CreatedAt time.Time
}
type OutingUseCase interface {
	OutingStudent(c *gin.Context, ctx context.Context, outingUUID uuid.UUID) error
	FindAllOutingStudent(ctx context.Context) ([]output.OutingStudentOutput, error)
	CountAllOutingStudent(ctx context.Context) (int, error)
	SearchOutingStudent(ctx context.Context, name string) ([]output.OutingStudentOutput, error)
}

type OutingRepository interface {
	SaveOutingStudent(ctx context.Context, outing *Outing) error
	ExistsOutingByAccountID(ctx context.Context, accountID primitive.ObjectID) (bool, error)
	DeleteOutingByAccountID(ctx context.Context, accountID primitive.ObjectID) error
	FindAllOuting(ctx context.Context) ([]Outing, error)
	FindByOutingAccountNameContaining(ctx context.Context, name string) ([]Outing, error)
}
