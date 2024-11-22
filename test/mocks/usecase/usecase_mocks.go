package mocks

import (
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"

	"github.com/stretchr/testify/mock"
)

type AuthUseCase struct {
	mock.Mock
}

func (m *AuthUseCase) SignUp(ctx context.Context, input input.SignUpInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *AuthUseCase) SignIn(ctx context.Context, input input.SignInInput) (output.TokenOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(output.TokenOutput), args.Error(1)
}

func (m *AuthUseCase) TokenReissue(ctx context.Context, refreshToken string) (output.TokenOutput, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(output.TokenOutput), args.Error(1)
}

func (m *AuthUseCase) SendAuthEmail(ctx context.Context, input input.SendEmaiInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *AuthUseCase) VerifyAuthCode(ctx context.Context, email string, authCode string) error {
	args := m.Called(ctx, email, authCode)
	return args.Error(0)
}
