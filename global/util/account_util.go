package util

import (
	"GOMS-BACKEND-GO/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetCurrentAccountID(c *gin.Context) (uint64, error) {
	account, ok := c.Get("account")
	if !ok || account == nil {
		return 0, fmt.Errorf("unauthorized")
	}

	accountModel, ok := account.(*model.Account)
	if !ok {
		return 0, fmt.Errorf("invalid account type")
	}
	return accountModel.ID, nil
}
