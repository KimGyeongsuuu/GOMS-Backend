package jwt

import (
	"GOMS-BACKEND-GO/global/config"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/output"
	"GOMS-BACKEND-GO/repository"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	jwtConfig   *config.JWTConfig
	rdb         *redis.Client
	refreshRepo *repository.RefreshTokenRepository
}

type GenerateToken interface {
	GenerateToken(ctx context.Context, accountId uint64, authority constant.Authority) (output.TokenOutput, error)
}

type ParseToken interface {
	ParseRefreshToken(refreshToken string) (string, error)
}

func NewToken(jwtConfig *config.JWTConfig, rdb *redis.Client, refreshRepo *repository.RefreshTokenRepository) *Token {
	return &Token{
		jwtConfig:   jwtConfig,
		rdb:         rdb,
		refreshRepo: refreshRepo,
	}
}

func (token *Token) GenerateToken(ctx context.Context, accountId uint64, authority constant.Authority) (output.TokenOutput, error) {
	accessToken, err := token.generateAccessToken(accountId, authority)
	if err != nil {
		return output.TokenOutput{}, err
	}
	refreshToken, err := token.generateRefreshToken(ctx, accountId)
	if err != nil {
		return output.TokenOutput{}, err
	}

	accessTokenExp := time.Now().Add(time.Duration(token.jwtConfig.AccessExp) * time.Second)
	refreshTokenExp := time.Now().Add(time.Duration(token.jwtConfig.RefreshExp) * time.Second)

	return output.TokenOutput{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccessTokenExp:  accessTokenExp.Format("2006-01-02 15:04:05"),
		RefreshTokenExp: refreshTokenExp.Format("2006-01-02 15:04:05"),
		Authority:       authority,
	}, nil
}

func (token *Token) generateAccessToken(accountId uint64, authority constant.Authority) (string, error) {
	claims := jwt.MapClaims{
		"sub":       accountId,
		"authority": authority,
		"exp":       time.Now().Add(time.Duration(token.jwtConfig.AccessExp) * time.Second).Unix(),
		"iat":       time.Now().Unix(),
	}
	token22 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token22.SignedString([]byte(token.jwtConfig.AccessSecret))
	if err != nil {
		return "", errors.New("failed to sign access token")
	}
	return signedToken, nil
}

func (token *Token) generateRefreshToken(ctx context.Context, accountId uint64) (string, error) {
	claims := jwt.MapClaims{
		"sub": accountId,
		"exp": time.Now().Add(time.Duration(token.jwtConfig.RefreshExp) * time.Second).Unix(),
		"iat": time.Now().Unix(),
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := tokenClaims.SignedString([]byte(token.jwtConfig.RefreshSecret))
	if err != nil {
		return "", errors.New("failed to sign refresh token")
	}

	refreshToken := &model.RefreshToken{
		RefreshToken: signedToken,
		AccountID:    accountId,
		ExpiredAt:    token.jwtConfig.RefreshExp,
	}

	token.refreshRepo.SaveRefreshToken(ctx, refreshToken)
	return signedToken, nil
}

func (token *Token) ParseRefreshToken(refreshToken string) (string, error) {
	if strings.HasPrefix(refreshToken, "Bearer ") {
		return strings.TrimPrefix(refreshToken, "Bearer "), nil
	}
	return "", errors.New("invalid refresh token format")
}
