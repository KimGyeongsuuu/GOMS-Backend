package repository

import (
	"GOMS-BACKEND-GO/global/config"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type OutingUUIDRepository struct {
	rdb          *redis.Client
	outingConfig *config.OutingConfig
}

func NewOutingUUIDRepository(rdb *redis.Client, outingConfig *config.OutingConfig) *OutingUUIDRepository {
	return &OutingUUIDRepository{
		rdb:          rdb,
		outingConfig: outingConfig,
	}
}
func (repository *OutingUUIDRepository) CreateOutingUUID(ctx context.Context) (uuid.UUID, error) {
	if repository.rdb == nil {
		return uuid.UUID{}, fmt.Errorf("redis client is nil")
	}
	if ctx == nil {
		return uuid.UUID{}, fmt.Errorf("context is nil")
	}

	outingUUID := uuid.New()
	key := "outing:uuid:" + outingUUID.String()
	value := outingUUID[:]

	fmt.Printf("Storing Key: %s, Value: %v\n", key, value)

	expiration := time.Duration(repository.outingConfig.OutingExp) * time.Second

	err := repository.rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return uuid.UUID{}, err
	}

	return outingUUID, nil
}

func (repository *OutingUUIDRepository) ExistsByOutingUUID(ctx context.Context, outingUUID uuid.UUID) (bool, error) {
	key := "outing:uuid:" + outingUUID.String()

	exists, err := repository.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists == 1, nil

}
