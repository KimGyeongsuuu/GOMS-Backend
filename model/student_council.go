package model

import (
	"GOMS-BACKEND-GO/model/data/output"
	"context"

	"github.com/google/uuid"
)

type StudentCouncilUseCase interface {
	CreateOuting(ctx context.Context) (uuid.UUID, error)
	FindAllAccount(ctx context.Context) ([]output.AccountOutput, error)
}

type OutingUUIDRepository interface {
	CreateOutingUUID(ctx context.Context) (uuid.UUID, error)
	ExistsByOutingUUID(ctx context.Context, outingUUID uuid.UUID) (bool, error)
}
