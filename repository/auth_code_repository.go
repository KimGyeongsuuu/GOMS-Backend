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
