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

	testcases := []struct {
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
			name: "이메일 인증을 하지 않는 사용자 Email 입니다.",
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
			expectedError: "", // 어떤 에러도 발생 X (회원가입 성공 !!)
		},
	}

	for _, test := range testcases {
		suite.Run(test.name, func() {
			test.setupMocks()                                                 // set up mocks 를 통해서 해당 메서드에 대한 mocking 작업 진행
			err := suite.authUsecase.SignUp(context.Background(), test.input) // test.input을 통해 테스트케이스 메서드 실행
			if test.expectedError != "" {                                     // test streuct에서 expectedError 공백이 아니면 에러가 발생하는 테스트케이스이므로 EqualError 로직 검증
				assert.EqualError(suite.T(), err, test.expectedError) // suite.T(실제 에러)와 test.expectedError(기대한 에러 값)을 비교
			} else {
				assert.NoError(suite.T(), err) //에러가 발생하지 않을 거라고 명시
			}

			// 호출 될 것이라고 기대한 mock method 모두 실행되었는지 assert
			suite.mockAccountRepo.AssertExpectations(suite.T())
			suite.mockPasswordUtil.AssertExpectations(suite.T())
			suite.mockAuthRepo.AssertExpectations(suite.T())
		})
	}
}
