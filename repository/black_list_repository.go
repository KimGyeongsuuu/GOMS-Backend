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

	blackListJson, err := json.Marshal(blackList)
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

func (repository *BlackListRepository) FindBlackListByAccountID(ctx context.Context, accountID uint64) (*model.BlackList, error) {
	key := "blacklist:" + strconv.FormatUint(accountID, 10)

	blackListJson, err := repository.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get value from Redis: %v", err)
	}

	var blackList model.BlackList
	err = json.Unmarshal([]byte(blackListJson), &blackList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal blackList: %v", err)
	}

	return &blackList, nil
}
func (repository *BlackListRepository) DeleteBlackList(ctx context.Context, blackList *model.BlackList) error {
	key := "blacklist:" + strconv.FormatUint(blackList.AccountID, 10)

	err := repository.rdb.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete value from Redis: %v", err)
	}

	fmt.Printf("Deleted Key: %s\n", key)
	return nil

}
