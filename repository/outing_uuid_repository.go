package repository

import (
	"GOMS-BACKEND-GO/global/config"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type OutingUUIDRepository struct {
	rdb              *redis.Client
	outingProperties *config.OutingProperties
}

func NewOutingUUIDRepository(rdb *redis.Client, outingProperties *config.OutingProperties) *OutingUUIDRepository {
	return &OutingUUIDRepository{
		rdb:              rdb,
		outingProperties: outingProperties,
	}
}

func (repository *OutingUUIDRepository) CreateOutingUUID(ctx context.Context) (uuid.UUID, error) {
	newUUID := uuid.New()

	err := repository.rdb.Set(ctx, "outing:uuid:"+newUUID.String(), newUUID[:], 0).Err()
	if err != nil {
		return uuid.UUID{}, err
	}

	return newUUID, nil
}
