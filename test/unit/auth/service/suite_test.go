package service

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/service"
	"GOMS-BACKEND-GO/test/mocks"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthServiceTestSuite struct {
	suite.Suite
	authUsecase      model.AuthUseCase
	mockAccountRepo  *mocks.MockAccountRepository
	mockAuthRepo     *mocks.MockAuthenticationRepository
	mockTokenAdapter *mocks.MockGenerateTokenAdapter
	mockPasswordUtil *mocks.MockPasswordUtil
}

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (suite *AuthServiceTestSuite) SetupSuite() {
	suite.mockAccountRepo = mocks.NewAccountRepository(suite.T())
	suite.mockAuthRepo = mocks.NewAuthenticationRepository(suite.T())
	suite.mockTokenAdapter = mocks.NewGenerateTokenAdapter(suite.T())
	suite.mockPasswordUtil = mocks.NewPasswordUtil(suite.T())

	suite.authUsecase = service.NewAuthService(
		suite.mockAccountRepo,
		suite.mockTokenAdapter,
		nil,
		nil,
		suite.mockAuthRepo,
		nil,
		suite.mockPasswordUtil,
	)
}
