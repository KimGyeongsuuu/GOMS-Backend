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
	"GOMS-BACKEND-GO/test/mocks"
)

func TestSignUp(t *testing.T) {
	testcase := []struct {
		name       string
		payload    input.SignUpInput
		on         func(mockAuthUseCase *mocks.MockAuthUseCase)
		statusCode int
	}{
		{
			name: "이미 존재하는 사용자 Email 입니다",
			payload: input.SignUpInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			on: func(mockAuthUseCase *mocks.MockAuthUseCase) {
				mockAuthUseCase.On("SignUp", mock.Anything, mock.AnythingOfType("input.SignUpInput")).
					Return(errors.New("email already exists"))
			},
			statusCode: http.StatusConflict,
		},
		{
			name: "인증되지 않은 사용자 Email 입니다.",
			payload: input.SignUpInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			on: func(mockAuthUseCase *mocks.MockAuthUseCase) {
				mockAuthUseCase.On("SignUp", mock.Anything, mock.AnythingOfType("input.SignUpInput")).
					Return(errors.New("authentication not found"))
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "회원가입 성공",
			payload: input.SignUpInput{
				Email:    "kskim@nurilab.com",
				Password: "rudtn1991!",
			},
			on: func(mockAuthUseCase *mocks.MockAuthUseCase) {
				mockAuthUseCase.On("SignUp", mock.Anything, mock.AnythingOfType("input.SignUpInput")).
					Return(nil)
			},
			statusCode: http.StatusCreated,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthUseCase := new(mocks.MockAuthUseCase)
			tc.on(mockAuthUseCase)

			authController := controller.NewAuthController(mockAuthUseCase)
			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.POST("/signup", authController.SignUp)

			body, err := json.Marshal(tc.payload)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.statusCode, rec.Code)
			mockAuthUseCase.AssertExpectations(t)
		})
	}

}
