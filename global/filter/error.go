package filter

import (
	"GOMS-BACKEND-GO/global/error/status"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorFilter struct{}

func NewErrorFilter() *ErrorFilter {
	return &ErrorFilter{}
}

func (f *ErrorFilter) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // c.Next()코드를 통해 다음 핸들러로 넘김
		// ErrorFilter의 역할은 다른 미들웨어나 핸들러에서 받은 요청을 가지고 요청 이후에 발생한 에러를 처리하기 위함이다.
		// 그렇기 때문에 다른 미들웨어, 핸들러가 먼저 요청을 처리할 수 있도록 함수가 시작되자 마자 c.Next() 통해 제어권을 넘김
		// 다른 요청이 먼저 실행이 되면 다른 요청에서의 c.Next()를 통해 마지막에 제어권을 받음

		for _, e := range c.Errors {
			statusErr, ok := e.Err.(*status.Err) // 제어권을 가지고 있을때 c.Errors를 통해 계속 Error가 생겼는지 확인하면서 Error가 발생했을때 에러 반환
			if !ok {
				statusErr = status.NewError(http.StatusInternalServerError, "internal server error")
			}

			c.AbortWithStatusJSON(statusErr.Code, gin.H{ // 반환된 에러를 가지고 에러 메세지를 커스텀 이후 상태코드와 함께 반환
				"message": statusErr.Message,
			})
		}
	}
}
