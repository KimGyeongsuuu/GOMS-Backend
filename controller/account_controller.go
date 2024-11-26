package controller

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	accountUseCase model.AccountUseCase
}

func NewAccountController(accountUseCase model.AccountUseCase) *AccountController {
	return &AccountController{
		accountUseCase: accountUseCase,
	}
}

func (controller *AccountController) WithDraw(ctx *gin.Context) {
	if err := controller.accountUseCase.WithDrawAccount(ctx, context.Background()); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
