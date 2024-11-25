package mocks

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

// Account
type MockAccountRepository struct {
	mock.Mock
}

func NewAccountRepository(t *testing.T) *MockAccountRepository {
	return &MockAccountRepository{
		Mock: mock.Mock{},
	}
}

func (a *MockAccountRepository) SaveAccount(ctx context.Context, account *model.Account) error {
	args := a.Called(ctx, account)
	return args.Error(0)
}

func (a *MockAccountRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := a.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (a *MockAccountRepository) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	args := a.Called(ctx, email)
	if result := args.Get(0); result != nil {
		return result.(*model.Account), args.Error(1)
	}
	return nil, args.Error(1)
}

func (a *MockAccountRepository) FindByAccountID(ctx context.Context, accountID uint64) (*model.Account, error) {
	args := a.Called(ctx, accountID)
	if result := args.Get(0); result != nil {
		return result.(*model.Account), args.Error(1)
	}
	return nil, args.Error(1)
}

func (a *MockAccountRepository) FindAllAccount(ctx context.Context) ([]model.Account, error) {
	args := a.Called(ctx)
	return args.Get(0).([]model.Account), args.Error(1)
}

func (a *MockAccountRepository) FindByAccountByStudentInfo(ctx context.Context, searchAccountInput *input.SearchAccountInput) ([]model.Account, error) {
	args := a.Called(ctx, searchAccountInput)
	return args.Get(0).([]model.Account), args.Error(1)
}

func (a *MockAccountRepository) UpdateAccountAuthority(ctx context.Context, authorityInput *input.UpdateAccountAuthorityInput) error {
	args := a.Called(ctx, authorityInput)
	return args.Error(0)
}

func (a *MockAccountRepository) DeleteAccount(ctx context.Context, account *model.Account) error {
	args := a.Called(ctx, account)
	return args.Error(0)
}

type AuthenticationRepository struct {
	mock.Mock
}

func NewAuthenticationRepository(t *testing.T) *AuthenticationRepository {
	return &AuthenticationRepository{
		Mock: mock.Mock{},
	}
}

func (a *AuthenticationRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := a.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (a *AuthenticationRepository) FindByEmail(ctx context.Context, email string) (*model.Authentication, error) {
	args := a.Called(ctx, email)
	if result := args.Get(0); result != nil {
		return result.(*model.Authentication), args.Error(1)
	}
	return nil, args.Error(1)
}

func (a *AuthenticationRepository) SaveAuthentication(ctx context.Context, auth *model.Authentication) error {
	args := a.Called(ctx, auth)
	return args.Error(0)
}

// RefreshToken
type RefreshTokenRepository struct {
	mock.Mock
}

func NewRefreshTokenRepository(t *testing.T) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		Mock: mock.Mock{},
	}
}

func (r *RefreshTokenRepository) FindRefreshTokenByRefreshToken(ctx context.Context, refreshToken string) (*model.RefreshToken, error) {
	args := r.Called(ctx, refreshToken)
	if result := args.Get(0); result != nil {
		return result.(*model.RefreshToken), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, refreshToken *model.RefreshToken) error {
	args := r.Called(ctx, refreshToken)
	return args.Error(0)
}

// AuthCode
type AuthCodeRepository struct {
	mock.Mock
}

func NewAuthCodeRepository(t *testing.T) *AuthCodeRepository {
	return &AuthCodeRepository{
		Mock: mock.Mock{},
	}
}

func (a *AuthCodeRepository) FindByEmail(ctx context.Context, email string) (*model.AuthCode, error) {
	args := a.Called(ctx, email)
	if result := args.Get(0); result != nil {
		return result.(*model.AuthCode), args.Error(1)
	}
	return nil, args.Error(1)
}

func (a *AuthCodeRepository) SaveAuthCode(ctx context.Context, authCode *model.AuthCode) error {
	args := a.Called(ctx, authCode)
	return args.Error(0)
}

// outing
type OutingRepository struct {
	mock.Mock
}

func NewOutingRepository(t *testing.T) *OutingRepository {
	return &OutingRepository{
		Mock: mock.Mock{},
	}
}

func (o *OutingRepository) SaveOutingStudnet(ctx context.Context, outing *model.Outing) error {
	args := o.Called(ctx, outing)
	return args.Error(0)
}

func (o *OutingRepository) ExistsOutingByAccountID(ctx context.Context, accountID uint64) (bool, error) {
	args := o.Called(ctx, accountID)
	return args.Bool(0), args.Error(1)
}

func (o *OutingRepository) DeleteOutingByAccountID(ctx context.Context, accountID uint64) error {
	args := o.Called(ctx, accountID)
	return args.Error(0)
}

func (o *OutingRepository) FindAllOuting(ctx context.Context) ([]model.Outing, error) {
	args := o.Called(ctx)
	if result := args.Get(0); result != nil {
		return result.([]model.Outing), args.Error(1)
	}
	return nil, args.Error(1)
}

func (o *OutingRepository) FindByOutingAccountNameContaining(ctx context.Context, name string) ([]model.Outing, error) {
	args := o.Called(ctx, name)
	if result := args.Get(0); result != nil {
		return result.([]model.Outing), args.Error(1)
	}
	return nil, args.Error(1)
}
