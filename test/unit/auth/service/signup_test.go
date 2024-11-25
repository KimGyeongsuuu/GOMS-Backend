package service

import (
	"context"
	"errors"

	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *AuthServiceTestSuite) TestSignUp() {

	testcase := []struct {
		name          string
		setupMocks    func()
		input         input.SignUpInput
		expectedError string
	}{
		{
			name: "이미 존재하는 사용자 Email 입니다.",
			setupMocks: func() {
				suite.mockAccountRepo.On("ExistsByEmail", mock.Anything, "kskim@nurilab.com").
					Return(true, nil).Once()
			},
			input: input.SignUpInput{
				Email:    "kskim@nurilab.com",
				Name:     "kimgyeongsu",
				Gender:   "MAN",
				Password: "rudtn1991!",
			},
			expectedError: "email already exists",
		},
		{
			name: "인증 객체가 존재하지 않습니다.",
			setupMocks: func() {
				suite.mockAccountRepo.On("ExistsByEmail", mock.Anything, "kskim@nurilab.com").
					Return(false, nil).Once()
				suite.mockAuthRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(nil, nil).Once()
			},
			input: input.SignUpInput{
				Email:    "kskim@nurilab.com",
				Name:     "kimgyeongsu",
				Gender:   "MAN",
				Password: "rudtn1991!",
			},
			expectedError: "authentication not found",
		},
		{
			name: "password 인코딩 중 오류 발생",
			setupMocks: func() {
				suite.mockAccountRepo.On("ExistsByEmail", mock.Anything, "kskim@nurilab.com").
					Return(false, nil).Once()
				suite.mockAuthRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Authentication{IsAuthenticated: true}, nil).Once()
				suite.mockPasswordUtil.On("EncodePassword", "rudtn1991!").
					Return("", errors.New("password encode error")).Once()
			},
			input: input.SignUpInput{
				Email:    "kskim@nurilab.com",
				Name:     "kimgyeongsu",
				Gender:   "MAN",
				Password: "rudtn1991!",
			},
			expectedError: "password encode error",
		},
		{
			name: "회원가입 성공",
			setupMocks: func() {
				suite.mockAccountRepo.On("ExistsByEmail", mock.Anything, "kskim@nurilab.com").
					Return(false, nil).Once()
				suite.mockAuthRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Authentication{IsAuthenticated: true}, nil).Once()
				suite.mockPasswordUtil.On("EncodePassword", "rudtn1991!").
					Return("encoded_password", nil).Once()
				suite.mockAccountRepo.On("SaveAccount", mock.Anything, mock.AnythingOfType("*model.Account")).
					Return(nil).Once()
			},
			input: input.SignUpInput{
				Email:    "kskim@nurilab.com",
				Name:     "kimgyeongsu",
				Gender:   "MAN",
				Password: "rudtn1991!",
			},
			expectedError: "",
		},
	}

	for _, test := range testcase {
		suite.Run(test.name, func() {
			test.setupMocks()
			err := suite.authUsecase.SignUp(context.Background(), test.input)
			if test.expectedError != "" {
				assert.EqualError(suite.T(), err, test.expectedError)
			} else {
				assert.NoError(suite.T(), err)
			}
			suite.mockAccountRepo.AssertExpectations(suite.T())
		})
	}
}
