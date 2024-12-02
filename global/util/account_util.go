package util

import (
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCurrentAccountID(c *gin.Context) (primitive.ObjectID, error) {
	account, ok := c.Get("account")
	if !ok || account == nil {
		return primitive.NilObjectID, status.NewError(http.StatusUnauthorized, "unauthorized")
	}

	accountModel, ok := account.(*model.Account)
	if !ok {
		return primitive.NilObjectID, status.NewError(http.StatusUnauthorized, "invalid account type")
	}
	return accountModel.ID, nil
}
