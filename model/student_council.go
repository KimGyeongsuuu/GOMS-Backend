package model

import (
	"context"

	"github.com/google/uuid"
)

type StudentCouncilUseCase interface {
	CreateOuting(ctx context.Context) (uuid.UUID, error)
}

type OutingUUIDRepository interface {
	CreateOutingUUID(ctx context.Context) (uuid.UUID, error)
}
