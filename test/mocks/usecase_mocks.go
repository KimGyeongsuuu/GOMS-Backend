package mocks

import (
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) SignUp(ctx context.Context, input input.SignUpInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockAuthUseCase) SignIn(ctx context.Context, input input.SignInInput) (output.TokenOutput, error) {
	args := m.Called(ctx, input)

	if tokenOutput, ok := args.Get(0).(output.TokenOutput); ok {
		return tokenOutput, nil
	}

	return output.TokenOutput{}, args.Error(0)
}

func (m *MockAuthUseCase) TokenReissue(ctx context.Context, refreshToken string) (output.TokenOutput, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(output.TokenOutput), args.Error(1)
}

func (m *MockAuthUseCase) SendAuthEmail(ctx context.Context, input input.SendEmaiInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockAuthUseCase) VerifyAuthCode(ctx context.Context, email string, authCode string) error {
	args := m.Called(ctx, email, authCode)
	return args.Error(0)
}
