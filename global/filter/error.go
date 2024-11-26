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
		c.Next()

		for _, e := range c.Errors {
			statusErr, ok := e.Err.(*status.Err)
			if !ok {
				statusErr = status.NewError(http.StatusInternalServerError, "internal server error")
			}

			c.AbortWithStatusJSON(statusErr.Code, gin.H{
				"message": statusErr.Message,
			})
		}
	}
}
