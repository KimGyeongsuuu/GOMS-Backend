package mocks

import (
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// token
type MockGenerateTokenAdapter struct {
	mock.Mock
}

func NewGenerateTokenAdapter(t *testing.T) *MockGenerateTokenAdapter {
	return &MockGenerateTokenAdapter{
		Mock: mock.Mock{},
	}
}

func (m *MockGenerateTokenAdapter) GenerateToken(ctx context.Context, accountId primitive.ObjectID, authority constant.Authority) (output.TokenOutput, error) {
	args := m.Called(ctx, accountId, authority)
	result, ok := args.Get(0).(output.TokenOutput)
	if !ok {
		return output.TokenOutput{}, args.Error(1)
	}
	return result, args.Error(1)
}

type MockPasswordUtil struct {
	mock.Mock
}

func NewPasswordUtil(t *testing.T) *MockPasswordUtil {
	return &MockPasswordUtil{
		Mock: mock.Mock{},
	}
}

func (m *MockPasswordUtil) EncodePassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordUtil) IsPasswordMatch(rawPassword string, encodedPassword string) (bool, error) {
	args := m.Called(rawPassword, encodedPassword)
	return args.Bool(0), args.Error(1)
}
