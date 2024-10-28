package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthenticationRepository struct {
	rdb *redis.Client
}

func NewAuthenticationRepository(rdb *redis.Client) *AuthenticationRepository {
	return &AuthenticationRepository{
		rdb: rdb,
	}
}

func (repository *AuthenticationRepository) SaveAuthentication(ctx context.Context, authentication *model.Authentication) error {
	key := "authentication:" + authentication.Email

	authenticationJson, err := json.Marshal(authentication)
	if err != nil {
		return fmt.Errorf("failed to marshal authCode: %v", err)
	}

	expiration := time.Until(authentication.ExpiredAt)

	err = repository.rdb.Set(ctx, key, authenticationJson, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in Redis: %v", err)
	}

	return nil
}

func (repository *AuthenticationRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	key := "authentication:" + email

	exists, err := repository.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence in Redis: %v", err)
	}

	return exists > 0, nil
}

func (repository *AuthenticationRepository) FindByEmail(ctx context.Context, email string) (*model.Authentication, error) {
	key := "authentication:" + email

	authenticationJson, err := repository.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get value from Redis: %v", err)
	}

	var authentication model.Authentication
	if err := json.Unmarshal([]byte(authenticationJson), &authentication); err != nil {
		return nil, fmt.Errorf("failed to unmarshal authentication: %v", err)
	}

	return &authentication, nil
}
