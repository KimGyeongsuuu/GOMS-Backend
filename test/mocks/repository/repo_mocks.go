package mocks

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

// Account
type AccountRepository struct {
	mock.Mock
}

func NewAccountRepository(t *testing.T) *AccountRepository {
	return &AccountRepository{
		Mock: mock.Mock{},
	}
}

func (a *AccountRepository) SaveAccount(ctx context.Context, account *model.Account) error {
	args := a.Called(ctx, account)
	return args.Error(0)
}

func (a *AccountRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := a.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (a *AccountRepository) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	args := a.Called(ctx, email)
	if result := args.Get(0); result != nil {
		return result.(*model.Account), args.Error(1)
	}
	return nil, args.Error(1)
}

func (a *AccountRepository) FindByAccountID(ctx context.Context, accountID uint64) (*model.Account, error) {
	args := a.Called(ctx, accountID)
	if result := args.Get(0); result != nil {
		return result.(*model.Account), args.Error(1)
	}
	return nil, args.Error(1)
}

func (a *AccountRepository) FindAllAccount(ctx context.Context) ([]model.Account, error) {
	args := a.Called(ctx)
	return args.Get(0).([]model.Account), args.Error(1)
}

func (a *AccountRepository) FindByAccountByStudentInfo(ctx context.Context, searchAccountInput *input.SearchAccountInput) ([]model.Account, error) {
	args := a.Called(ctx, searchAccountInput)
	return args.Get(0).([]model.Account), args.Error(1)
}

func (a *AccountRepository) UpdateAccountAuthority(ctx context.Context, authorityInput *input.UpdateAccountAuthorityInput) error {
	args := a.Called(ctx, authorityInput)
	return args.Error(0)
}

func (a *AccountRepository) DeleteAccount(ctx context.Context, account *model.Account) error {
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

// token
type GenerateTokenAdapter struct {
	mock.Mock
}

func NewGenerateTokenAdapter() *GenerateTokenAdapter {
	return &GenerateTokenAdapter{}
}

func (m *GenerateTokenAdapter) GenerateToken(ctx context.Context, accountId uint64, authority constant.Authority) (output.TokenOutput, error) {
	args := m.Called(ctx, accountId, authority)

	if result := args.Get(0); result != nil {
		return result.(output.TokenOutput), args.Error(1)
	}
	return output.TokenOutput{}, args.Error(1)
}
