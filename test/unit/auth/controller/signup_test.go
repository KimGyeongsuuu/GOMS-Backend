package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"GOMS-BACKEND-GO/controller"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/test/mocks"
)

func (suite *AuthControllerTestSuite) TestSignUp() {
	testcases := []struct {
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
			statusCode: http.StatusConflict, // conroller에서 발생할 것이라고 예상되는 상태코드
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

	for _, test := range testcases {
		suite.Run(test.name, func() {
			mockAuthUseCase := new(mocks.MockAuthUseCase) // mock 메서드 생성 및 설정
			test.on(mockAuthUseCase)

			authController := controller.NewAuthController(mockAuthUseCase) // controller를 생성하고 mock 객체 주입
			gin.SetMode(gin.TestMode)                                       // gin을 통해 라우터 설정
			router := gin.Default()
			router.POST("/signup", authController.SignUp) // 라우터에 sign up 등록

			body, err := json.Marshal(test.payload) // payload를 json로 변환하고 body에 등록
			assert.NoError(suite.T(), err)          // 에러가 발생하지 않는다는 코드

			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(suite.T(), test.statusCode, rec.Code) // 실제 발생 상태코드와, 테스트코드에서 작성해놓은 statusCode가 같다는 Equal
			mockAuthUseCase.AssertExpectations(suite.T())
		})
	}

}
