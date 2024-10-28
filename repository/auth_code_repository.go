package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthCodeRepository struct {
	rdb *redis.Client
}

func NewAuthCodeRepository(rdb *redis.Client) *AuthCodeRepository {
	return &AuthCodeRepository{
		rdb: rdb,
	}
}

func (repository *AuthCodeRepository) SaveAuthCode(ctx context.Context, authCode *model.AuthCode) error {
	key := "authcode:" + authCode.Email

	authCodeJson, err := json.Marshal(authCode)
	if err != nil {
		return fmt.Errorf("failed to marshal authCode: %v", err)
	}

	expiration := time.Until(authCode.ExpiredAt)

	err = repository.rdb.Set(ctx, key, authCodeJson, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in Redis: %v", err)
	}

	return nil
}

func (repository *AuthCodeRepository) FindByEmail(ctx context.Context, email string) (*model.AuthCode, error) {
	key := "authcode:" + email

	authCodeJson, err := repository.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get value from Redis: %v", err)
	}

	var authCode model.AuthCode
	if err := json.Unmarshal([]byte(authCodeJson), &authCode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal authentication: %v", err)
	}

	return &authCode, nil
}
