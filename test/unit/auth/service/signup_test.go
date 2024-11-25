package service

import (
	"context"
	"errors"
	"testing"

	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/service"
	"GOMS-BACKEND-GO/test/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	mockAccountRepo := mocks.NewAccountRepository(t)
	mockAuthRepo := mocks.NewAuthenticationRepository(t)
	mockPasswordUtil := mocks.NewPasswordUtil(t)

	authService := service.NewAuthService(
		mockAccountRepo,
		nil,
		nil,
		nil,
		mockAuthRepo,
		nil,
		mockPasswordUtil,
	)

	testcase := []struct {
		name          string
		setupMocks    func()
		input         input.SignUpInput
		expectedError string
	}{
		{
			name: "이미 존재하는 사용자 Email 입니다.",
			setupMocks: func() {
				mockAccountRepo.On("ExistsByEmail", mock.Anything, "kskim@nurilab.com").
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
				mockAccountRepo.On("ExistsByEmail", mock.Anything, "kskim@nurilab.com").
					Return(false, nil).Once()
				mockAuthRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
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
				mockAccountRepo.On("ExistsByEmail", mock.Anything, "kskim@nurilab.com").
					Return(false, nil).Once()
				mockAuthRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Authentication{IsAuthenticated: true}, nil).Once()
				mockPasswordUtil.On("EncodePassword", "rudtn1991!").
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
				mockAccountRepo.On("ExistsByEmail", mock.Anything, "kskim@nurilab.com").
					Return(false, nil).Once()
				mockAuthRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Authentication{IsAuthenticated: true}, nil).Once()
				mockPasswordUtil.On("EncodePassword", "rudtn1991!").
					Return("encoded_password", nil).Once()
				mockAccountRepo.On("SaveAccount", mock.Anything, mock.AnythingOfType("*model.Account")).
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
		t.Run(test.name, func(t *testing.T) {
			test.setupMocks()

			err := authService.SignUp(context.Background(), test.input)

			if test.expectedError != "" { // 에러 발생 O
				assert.EqualError(t, err, test.expectedError)
			} else { // 에러 발생 X
				assert.NoError(t, err)
			}

			mockAccountRepo.AssertExpectations(t)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}
