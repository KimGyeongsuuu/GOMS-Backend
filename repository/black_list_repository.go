package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type BlackListRepository struct {
	rdb *redis.Client
}

func NewBlackListRepository(rdb *redis.Client) *BlackListRepository {
	return &BlackListRepository{
		rdb: rdb,
	}
}

func (repository *BlackListRepository) SaveBlackList(ctx context.Context, blackList *model.BlackList) error {

	key := "blacklist:" + strconv.FormatUint(blackList.AccountID, 10)

	blackListJson, err := json.Marshal(blackList.AccountID)
	if err != nil {
		return fmt.Errorf("failed to marshal blackList: %v", err)
	}

	expiration := time.Duration(blackList.ExpiredAt) * time.Second

	fmt.Printf("Storing Key: %s, Value: %s\n", key, string(blackListJson))

	err = repository.rdb.Set(ctx, key, blackListJson, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in Redis: %v", err)
	}

	return nil

}
