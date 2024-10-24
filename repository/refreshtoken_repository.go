package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RefreshTokenRepository struct {
	rdb *redis.Client
}

func NewRefreshTokenRepository(rdb *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		rdb: rdb,
	}
}

func (repository *RefreshTokenRepository) SaveRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error {
	if repository.rdb == nil {
		return fmt.Errorf("redis client is nil")
	}
	if ctx == nil {
		return fmt.Errorf("context is nil")
	}

	tokenJSON, err := json.Marshal(refreshToken)
	if err != nil {
		return err
	}

	key := "refresh:token:" + refreshToken.RefreshToken
	expiration := time.Duration(refreshToken.ExpiredAt) * time.Second

	err = repository.rdb.Set(ctx, key, tokenJSON, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}
func (repository *RefreshTokenRepository) FindRefreshTokenByRefreshToken(ctx context.Context, refreshToken string) (*model.RefreshToken, error) {
	if repository.rdb == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	key := "refresh:token:" + refreshToken

	tokenJSON, err := repository.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, err
	}

	var storedToken model.RefreshToken
	err = json.Unmarshal([]byte(tokenJSON), &storedToken)
	if err != nil {
		return nil, err
	}

	return &storedToken, nil
}

func (repository *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error {
	err := repository.rdb.Del(ctx, refreshToken.RefreshToken).Err()
	if err != nil {
		return fmt.Errorf("failed to delete refresh token from Redis: %v", err)
	}
	return nil
}
