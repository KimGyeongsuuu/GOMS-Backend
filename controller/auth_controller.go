package controller

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
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

// Sign Up Router
// @Summary Sign Up for a new user
// @Description email인증 후 학과 정보를 통해 사용자 회원가입.
// @Accept json
// @Produce json
// @Param user body input.SignUpInput true "User information"
// @Success 201 {object} string "User created successfully"
// @Router /api/v1/auth/signup [post]
func (controller *AuthController) SignUp(ctx *gin.Context) {
	var input input.SignUpInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}

	if err := controller.authUseCase.SignUp(context.Background(), input); err != nil {
		ctx.Error(err)
	}
	ctx.Status(http.StatusCreated)
}

// Sign In Router
// @Summary Sign In for a user
// @Description email과 password를 통해 로그인 후 토큰 발급.
// @Accept json
// @Produce json
// @Param user body input.SignInInput true "User login credentials"
// @Success 200 {object} output.TokenOutput "로그인 성공"
// @Router /api/v1/auth/signin [post]
func (controller *AuthController) SignIn(ctx *gin.Context) {
	var input input.SignInInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}

	token, err := controller.authUseCase.SignIn(context.Background(), input)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, output.TokenOutput{
		AccessToken:     token.AccessToken,
		RefreshToken:    token.RefreshToken,
		AccessTokenExp:  token.AccessTokenExp,
		RefreshTokenExp: token.RefreshTokenExp,
	})
}

// Token Reissue Router
// @Summary 토큰 재발급
// @Description RefreshToken를 header로 받아서 요청.
// @Accept json
// @Produce json
// @Param RefreshToken header string true "Refresh Token"
// @Success 200 {object} output.TokenOutput "토큰 재발급 성공"
// @Router /api/v1/auth [Patch]
func (controller *AuthController) TokenReissue(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("RefreshToken")

	token, err := controller.authUseCase.TokenReissue(context.Background(), refreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, output.TokenOutput{
		AccessToken:     token.AccessToken,
		RefreshToken:    token.RefreshToken,
		AccessTokenExp:  token.AccessTokenExp,
		RefreshTokenExp: token.RefreshTokenExp,
	})
}

func (controller *AuthController) SendAuthEmail(ctx *gin.Context) {
	var input input.SendEmaiInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(err)
		return
	}
	err := controller.authUseCase.SendAuthEmail(context.Background(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller *AuthController) VerifyAuthCode(ctx *gin.Context) {

	email := ctx.Query("email")
	authCode := ctx.Query("authCode")

	err := controller.authUseCase.VerifyAuthCode(context.Background(), email, authCode)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)

}
