package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	key := "blacklist:" + blackList.AccountID.Hex()

	blackListJson, err := json.Marshal(blackList)
	if err != nil {
		return fmt.Errorf("failed to marshal blackList: %v", err)
	}

	expiration := time.Duration(blackList.ExpiredAt) * time.Second

	err = repository.rdb.Set(ctx, key, blackListJson, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in Redis: %v", err)
	}

	return nil
}

func (repository *BlackListRepository) FindBlackListByAccountID(ctx context.Context, accountID primitive.ObjectID) (*model.BlackList, error) {
	key := "blacklist:" + accountID.Hex()

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
	key := "blacklist:" + blackList.AccountID.Hex()

	err := repository.rdb.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete value from Redis: %v", err)
	}

	return nil

}

func (repository *BlackListRepository) ExistsByAccountID(ctx context.Context, accountID primitive.ObjectID) (bool, error) {
	key := "blacklist:" + accountID.Hex()

	count, err := repository.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence in Redis: %v", err)
	}

	return count > 0, nil
}
