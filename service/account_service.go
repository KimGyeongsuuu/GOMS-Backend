package service

import (
	"GOMS-BACKEND-GO/global/util"
	"GOMS-BACKEND-GO/model"
	"context"

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
		return err
	}
	account, err := service.accountRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return err
	}
	return service.accountRepo.DeleteAccount(ctx, account)
}
