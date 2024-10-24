package model

import (
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"

	"github.com/google/uuid"
)

type StudentCouncilUseCase interface {
	CreateOuting(ctx context.Context) (uuid.UUID, error)
	FindAllAccount(ctx context.Context) ([]output.AccountOutput, error)
	SearchAccount(ctx context.Context, accountInput *input.SearchAccountInput) ([]output.AccountOutput, error)
	UpdateAccountAuthority(ctx context.Context, authorityInput *input.UpdateAccountAuthorityInput) error
	AddBlackList(ctx context.Context, accountID uint64) error
	ExcludeBlackList(ctx context.Context, accountID uint64) error
}

type OutingUUIDRepository interface {
	CreateOutingUUID(ctx context.Context) (uuid.UUID, error)
	ExistsByOutingUUID(ctx context.Context, outingUUID uuid.UUID) (bool, error)
}
