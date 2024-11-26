package util

import (
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCurrentAccountID(c *gin.Context) (uint64, error) {
	account, ok := c.Get("account")
	if !ok || account == nil {
		return 0, status.NewError(http.StatusUnauthorized, "unauthorized")
	}

	accountModel, ok := account.(*model.Account)
	if !ok {
		return 0, status.NewError(http.StatusUnauthorized, "invalid account type")
	}
	return accountModel.ID, nil
}
