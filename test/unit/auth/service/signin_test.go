package service

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"GOMS-BACKEND-GO/service"
	"GOMS-BACKEND-GO/test/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestSignIn(t *testing.T) {
	mockAccountRepo := mocks.NewAccountRepository(t)
	mockTokenAdapter := mocks.NewGenerateTokenAdapter(t)
	mockPasswordUtil := mocks.NewPasswordUtil(t)

	authService := service.NewAuthService(
		mockAccountRepo,
		mockTokenAdapter,
		nil,
		nil,
		nil,
		nil,
		mockPasswordUtil,
	)

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
				mockAccountRepo.On("FindByEmail", mock.Anything, "kimks@nurilab.com").
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
				mockAccountRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Account{
						ID:       1,
						Email:    "kskim@nurilab.com",
						Password: string(hashedPassword),
					}, nil).Once()
				mockPasswordUtil.On("IsPasswordMatch", "wrongPassword", mock.Anything).
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
				mockAccountRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Account{
						ID:       1,
						Email:    "kskim@nurilab.com",
						Password: string(hashedPassword),
					}, nil).Once()
				mockPasswordUtil.On("IsPasswordMatch", mock.Anything, mock.Anything).
					Return(true, nil).Once()
				mockTokenAdapter.On("GenerateToken", mock.Anything, uint64(1), mock.Anything).
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
				mockAccountRepo.On("FindByEmail", mock.Anything, "kskim@nurilab.com").
					Return(&model.Account{
						ID:       1,
						Email:    "kskim@nurilab.com",
						Password: string(hashedPassword),
					}, nil).Once()
				mockPasswordUtil.On("IsPasswordMatch", "rudtn1991!", mock.Anything).
					Return(true, nil).Once()
				mockTokenAdapter.On("GenerateToken", mock.Anything, uint64(1), mock.Anything).
					Return(output.TokenOutput{
						AccessToken:     "accessToken",
						RefreshToken:    "refreshToken",
						AccessTokenExp:  "2024-01-01 00:00:00",
						RefreshTokenExp: "2024-01-01 00:00:00",
						Authority:       constant.ROLE_STUDENT,
					}, nil).Once()
			},
			input: input.SignInInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			expectedOutput: output.TokenOutput{
				AccessToken:     "accessToken",
				RefreshToken:    "refreshToken",
				AccessTokenExp:  "2024-01-01 00:00:00",
				RefreshTokenExp: "2024-01-01 00:00:00",
				Authority:       constant.ROLE_STUDENT,
			},
			expectedError: "",
		},
	}
	for _, test := range testcase {
		t.Run(test.name, func(t *testing.T) {
			test.setupMocks()

			actualOutput, err := authService.SignIn(context.Background(), test.input)
			assert.Equal(t, test.expectedOutput, actualOutput)

			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockAccountRepo.AssertExpectations(t)
			mockTokenAdapter.AssertExpectations(t)
			mockPasswordUtil.AssertExpectations(t)
		})
	}
}
