package controller

import (
	"GOMS-BACKEND-GO/controller"
	"GOMS-BACKEND-GO/test/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AuthControllerTestSuite struct {
	suite.Suite
	mockAuthUseCase *mocks.MockAuthUseCase
	authController  *controller.AuthController
	router          *gin.Engine
}

func TestAuthControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}

func (suite *AuthControllerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.mockAuthUseCase = new(mocks.MockAuthUseCase)
	suite.authController = controller.NewAuthController(suite.mockAuthUseCase)
	suite.router = gin.Default()
	suite.router.POST("/signin", suite.authController.SignIn)
	suite.router.POST("/signup", suite.authController.SignUp)
}
