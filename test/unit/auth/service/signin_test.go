package service

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (suite *AuthServiceTestSuite) TestSignIn() {
	testcase := []struct {
		name           string
		setupMocks     func()
		input          input.SignInInput
		expectedOutput output.TokenOutput
		expectedError  string
	}{
		{
			name: "존재하지 않는 사용자 계정입니다.",
			setupMocks: func() {
				suite.mockAccountRepo.On("FindByEmail", mock.Anything, "kimks@nurilab.com").
					Return(nil, nil).Once()
			},
			input: input.SignInInput{
				Email:    "kimks@nurilab.com",
				Password: "rudtn1991!",
			},
			expectedOutput: output.TokenOutput{},
			expectedError:  "not found account",
		},
		{
			name: "비밀번호가 일치하지 않습니다.",
			setupMocks: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("rudtn1991!"), bcrypt.DefaultCost)
				suite.mockAccountRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Account{
						ID:       primitive.NewObjectID(),
						Email:    "kskim@nurilab.com",
						Password: string(hashedPassword),
					}, nil).Once()
				suite.mockPasswordUtil.On("IsPasswordMatch", "wrongPassword", mock.Anything).
					Return(false, nil).Once()
			},
			input: input.SignInInput{
				Email:    "kskim@nurilab.com",
				Password: "wrongPassword",
			},
			expectedOutput: output.TokenOutput{},
			expectedError:  "mis match password",
		},
		{
			name: "토큰 생성 중 오류 발생",
			setupMocks: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("rudtn1991!"), bcrypt.DefaultCost)
				suite.mockAccountRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Account{
						ID:       primitive.NewObjectID(),
						Email:    "kskim@nurilab.com",
						Password: string(hashedPassword),
					}, nil).Once()
				suite.mockPasswordUtil.On("IsPasswordMatch", mock.Anything, mock.Anything).
					Return(true, nil).Once()
				suite.mockTokenAdapter.On("GenerateToken", mock.Anything, uint64(1), mock.Anything).
					Return(output.TokenOutput{}, errors.New("token generate error")).Once()
			},
			input: input.SignInInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			expectedOutput: output.TokenOutput{},
			expectedError:  "token generate error",
		},
		{
			name: "로그인 성공",
			setupMocks: func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("rudtn1991!"), bcrypt.DefaultCost)
				suite.mockAccountRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Account{
						ID:       primitive.NewObjectID(),
						Email:    "kskim@nurilab.com",
						Password: string(hashedPassword),
					}, nil).Once()
				suite.mockPasswordUtil.On("IsPasswordMatch", "rudtn1991!", mock.Anything).
					Return(true, nil).Once()
				suite.mockTokenAdapter.On("GenerateToken", mock.Anything, uint64(1), mock.Anything).
					Return(output.TokenOutput{
						AccessToken:  "accessToken",
						RefreshToken: "refreshToken",
						Authority:    constant.ROLE_STUDENT,
					}, nil).Once()
			},
			input: input.SignInInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			expectedOutput: output.TokenOutput{
				AccessToken:  "accessToken",
				RefreshToken: "refreshToken",
				Authority:    constant.ROLE_STUDENT,
			},
			expectedError: "",
		},
	}
	for _, test := range testcase {
		suite.Run(test.name, func() {
			test.setupMocks()

			actualOutput, err := suite.authUsecase.SignIn(context.Background(), test.input)
			assert.Equal(suite.T(), test.expectedOutput, actualOutput)

			if test.expectedError != "" {
				assert.EqualError(suite.T(), err, test.expectedError)
			} else {
				assert.NoError(suite.T(), err)
			}

			suite.mockAccountRepo.AssertExpectations(suite.T())
			suite.mockTokenAdapter.AssertExpectations(suite.T())
			suite.mockPasswordUtil.AssertExpectations(suite.T())
		})
	}
}
