package model

import (
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Late struct {
	LateID    primitive.ObjectID `bson:"_id,omitempty"`
	AccountID primitive.ObjectID `bson:"account_id"`
	CreatedAt time.Time
}

type LateUseCase interface {
	GetTop3LateStudent(ctx context.Context) ([]output.LateOutput, error)
}

type LateRepository interface {
	FindTop3ByOrderByAccountDesc(ctx context.Context) ([]Late, error)
	FindLateByCreatedAt(ctx context.Context, date time.Time) ([]Late, error)
}
