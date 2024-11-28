package service

import (
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/global/util"
	"GOMS-BACKEND-GO/model"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountService struct {
	accountRepo model.AccountRepository
}

func NewAccountService(accountRepo model.AccountRepository) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
	}
}

func (service *AccountService) WithDrawAccount(c *gin.Context, ctx context.Context) error {
	accountID, err := util.GetCurrentAccountID(c)
	if err != nil {
		return status.NewError(http.StatusUnauthorized, "get current account id unauthorized")
	}
	account, err := service.accountRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return status.NewError(http.StatusInternalServerError, "find by id account id error")
	}
	return service.accountRepo.DeleteAccount(ctx, account)
}
