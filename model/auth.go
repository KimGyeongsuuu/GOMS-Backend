package model

import (
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
)

type RefreshToken struct {
	RefreshToken string
	AccountID    uint64
	ExpiredAt    int64
}

type AuthUseCase interface {
	SignUp(ctx context.Context, input *input.SignUpInput) error
	SignIn(ctx context.Context, input *input.SignInInput) (output.TokenOutput, error)
}

type RefreshTokenRepository interface {
	SaveRefreshToken(ctx context.Context, refreshToken *RefreshToken) error
}
