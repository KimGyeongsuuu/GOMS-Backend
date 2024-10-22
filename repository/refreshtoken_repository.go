package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"encoding/json"
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

func (repository *RefreshTokenRepository) CreateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error {

	tokenJSON, err := json.Marshal(refreshToken)
	if err != nil {
		return err
	}

	expiration := time.Duration(refreshToken.ExpiredAt) * time.Second
	err = repository.rdb.Set(ctx, refreshToken.RefreshToken, tokenJSON, expiration).Err()
	return err
}
