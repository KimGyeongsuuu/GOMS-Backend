package jwt

import (
	"GOMS-BACKEND-GO/global/config"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/output"
	"GOMS-BACKEND-GO/repository"
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

type GenerateTokenAdapter struct {
	jwtConfig   *config.JWTConfig
	rdb         *redis.Client
	refreshRepo *repository.RefreshTokenRepository
}

func NewGenerateTokenAdapter(jwtConfig *config.JWTConfig, rdb *redis.Client, refreshRepo *repository.RefreshTokenRepository) *GenerateTokenAdapter {
	return &GenerateTokenAdapter{
		jwtConfig:   jwtConfig,
		rdb:         rdb,
		refreshRepo: refreshRepo,
	}
}

func (adapter *GenerateTokenAdapter) GenerateToken(ctx context.Context, accountId uint64, authority constant.Authority) (output.TokenOutput, error) {
	accessToken, err := adapter.generateAccessToken(accountId, authority)
	if err != nil {
		return output.TokenOutput{}, err
	}
	refreshToken, err := adapter.generateRefreshToken(ctx, accountId)
	if err != nil {
		return output.TokenOutput{}, err
	}

	accessTokenExp := time.Now().Add(time.Duration(adapter.jwtConfig.AccessExp) * time.Second)
	refreshTokenExp := time.Now().Add(time.Duration(adapter.jwtConfig.RefreshExp) * time.Second)

	err = adapter.rdb.Set(context.Background(), refreshToken, accountId, time.Duration(adapter.jwtConfig.RefreshExp)*time.Second).Err()
	if err != nil {
		return output.TokenOutput{}, err
	}

	return output.TokenOutput{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessTokenExp:  accessTokenExp.Format("2006-01-02 15:04:05"),
		RefreshTokenExp: refreshTokenExp.Format("2006-01-02 15:04:05"),
		Authority:       authority,
	}, nil
}

func (adapter *GenerateTokenAdapter) generateAccessToken(accountId uint64, authority constant.Authority) (string, error) {
	claims := jwt.MapClaims{
		"sub":       accountId,
		"accountID": accountId,
		"authority": authority,
		"exp":       time.Now().Add(time.Duration(adapter.jwtConfig.AccessExp) * time.Second).Unix(),
		"iat":       time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(adapter.jwtConfig.AccessSecret))
	if err != nil {
		return "", errors.New("failed to sign access token")
	}
	return signedToken, nil
}

func (adapter *GenerateTokenAdapter) generateRefreshToken(ctx context.Context, accountId uint64) (string, error) {
	claims := jwt.MapClaims{
		"sub": accountId,
		"exp": time.Now().Add(time.Duration(adapter.jwtConfig.RefreshExp) * time.Second).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(adapter.jwtConfig.RefreshSecret))
	if err != nil {
		return "", errors.New("failed to sign refresh token")
	}

	refreshToken := &model.RefreshToken{
		RefreshToken: signedToken,
		AccountID:    accountId,
		ExpiredAt:    adapter.jwtConfig.RefreshExp,
	}

	adapter.refreshRepo.SaveRefreshToken(ctx, refreshToken)
	return signedToken, nil
}
