package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewRefreshTokenRepository(db *gorm.DB, rdb *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db:  db,
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
