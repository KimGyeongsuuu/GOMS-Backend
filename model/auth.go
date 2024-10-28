package model

import (
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"time"
)

type RefreshToken struct {
	RefreshToken string
	AccountID    uint64
	ExpiredAt    int64
}

type AuthCode struct {
	Email     string
	AuthCode  string
	ExpiredAt time.Time
}

type Authentication struct {
	Email           string
	AttemptCount    int // 인증번호 전송 횟수
	AuthCodeCount   int // 인증번호 검증 횟수
	IsAuthenticated bool
	ExpiredAt       time.Time
}

type AuthUseCase interface {
	SignUp(ctx context.Context, input *input.SignUpInput) error
	SignIn(ctx context.Context, input *input.SignInInput) (output.TokenOutput, error)
	TokenReissue(ctx context.Context, refreshToken string) (output.TokenOutput, error)
	SendAuthEmail(ctx context.Context, input *input.SendEmaiInput) error
}

type RefreshTokenRepository interface {
	SaveRefreshToken(ctx context.Context, refreshToken *RefreshToken) error
	FindRefreshTokenByRefreshToken(ctx context.Context, RefreshToken string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, RefreshToken *RefreshToken) error
}

type AuthCodeRepository interface {
	SaveAuthCode(ctx context.Context, authCode *AuthCode) error
}

type AuthenticationRepository interface {
	SaveAuthentication(ctx context.Context, authentication *Authentication) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*Authentication, error)
}
