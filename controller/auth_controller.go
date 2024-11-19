package controller

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUseCase model.AuthUseCase
}

func NewAuthController(authUseCase model.AuthUseCase) *AuthController {
	return &AuthController{
		authUseCase: authUseCase,
	}
}

func (controller *AuthController) SignUp(ctx *gin.Context) {
	var input input.SignUpInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := controller.authUseCase.SignUp(context.Background(), &input); err != nil {
		if err.Error() == "email already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func (controller *AuthController) SignIn(ctx *gin.Context) {
	var input input.SignInInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	token, err := controller.authUseCase.SignIn(context.Background(), &input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"TokenOutput": token})
}

func (controller *AuthController) TokenReissue(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("RefreshToken")

	token, err := controller.authUseCase.TokenReissue(context.Background(), refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"TokenOutput": token})
}

func (controller *AuthController) SendAuthEmail(ctx *gin.Context) {
	var input input.SendEmaiInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	err := controller.authUseCase.SendAuthEmail(context.Background(), &input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *AuthController) VerifyAuthCode(ctx *gin.Context) {

	email := ctx.Query("email")
	authCode := ctx.Query("authCode")

	err := controller.authUseCase.VerifyAuthCode(context.Background(), email, authCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)

}
