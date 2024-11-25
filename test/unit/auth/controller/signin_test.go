package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"GOMS-BACKEND-GO/controller"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"GOMS-BACKEND-GO/test/mocks"
)

func TestSignIn(t *testing.T) {
	testcase := []struct {
		name           string
		payload        input.SignInInput
		on             func(mockAuthUseCase *mocks.MockAuthUseCase)
		expectedOutput output.TokenOutput
		statusCode     int
	}{
		{
			name: "존재하지 않는 사용자 계정입니다.",
			payload: input.SignInInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			on: func(mockAuthUseCase *mocks.MockAuthUseCase) {
				mockAuthUseCase.On("SignIn", mock.Anything, mock.AnythingOfType("input.SignInInput")).
					Return(errors.New("not found account")).Once()
			},
			expectedOutput: output.TokenOutput{},
			statusCode:     http.StatusNotFound,
		},
		{
			name: "비밀번호가 일치하지 않습니다.",
			payload: input.SignInInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			on: func(mockAuthUseCase *mocks.MockAuthUseCase) {
				mockAuthUseCase.On("SignIn", mock.Anything, mock.AnythingOfType("input.SignInInput")).
					Return(errors.New("mis match password")).Once()
			},
			expectedOutput: output.TokenOutput{},
			statusCode:     http.StatusUnauthorized,
		},
		{
			name: "토큰 생성 중 오류 발생",
			payload: input.SignInInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			on: func(mockAuthUseCase *mocks.MockAuthUseCase) {
				mockAuthUseCase.On("SignIn", mock.Anything, mock.AnythingOfType("input.SignInInput")).
					Return(errors.New("token generate error")).Once()
			},
			expectedOutput: output.TokenOutput{},
			statusCode:     http.StatusInternalServerError,
		},
		{
			name: "로그인 성공",
			payload: input.SignInInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			on: func(mockAuthUseCase *mocks.MockAuthUseCase) {
				mockAuthUseCase.On("SignIn", mock.Anything, mock.AnythingOfType("input.SignInInput")).
					Return(output.TokenOutput{
						AccessToken:  "accessToken",
						RefreshToken: "refreshToken",
					}, nil).Once()
			},
			expectedOutput: output.TokenOutput{
				AccessToken:  "accessToken",
				RefreshToken: "refreshToken",
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.MockAuthUseCase)
			tc.on(mockAuthUseCase)

			authController := controller.NewAuthController(mockAuthUseCase)
			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.POST("/signin", authController.SignIn)

			body, err := json.Marshal(tc.payload)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// 토큰 발급 상태코드 200 StatusOK
			if tc.statusCode == http.StatusOK {
				var result map[string]output.TokenOutput

				err := json.Unmarshal(rec.Body.Bytes(), &result)
				assert.NoError(t, err)

				actualOutput := result["TokenOutput"]

				assert.Equal(t, tc.expectedOutput, actualOutput)
			}
			assert.Equal(t, tc.statusCode, rec.Code)
			mockAuthUseCase.AssertExpectations(t)
		})
	}
}
