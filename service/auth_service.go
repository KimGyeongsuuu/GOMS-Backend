package service

import (
	"GOMS-BACKEND-GO/global/auth/jwt"
	"GOMS-BACKEND-GO/global/email"
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/global/util"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/rand"
)

type AuthService struct {
	accountRepo        model.AccountRepository
	tokenAdapter       jwt.GenerateToken
	tokenParser        jwt.ParseToken
	refreshTokenRepo   model.RefreshTokenRepository
	authenticationRepo model.AuthenticationRepository
	authCodeRepo       model.AuthCodeRepository
	utilPassword       util.UtilPassword
}

func NewAuthService(
	accountRepo model.AccountRepository,
	tokenAdapter jwt.GenerateToken,
	tokenParser jwt.ParseToken,
	refreshTokenRepo model.RefreshTokenRepository,
	authenticationRepo model.AuthenticationRepository,
	authCodeRepo model.AuthCodeRepository,
	utilPassword util.UtilPassword,
) model.AuthUseCase {
	return &AuthService{
		accountRepo:        accountRepo,
		tokenAdapter:       tokenAdapter,
		refreshTokenRepo:   refreshTokenRepo,
		tokenParser:        tokenParser,
		authenticationRepo: authenticationRepo,
		authCodeRepo:       authCodeRepo,
		utilPassword:       utilPassword,
	}
}

func (service *AuthService) SignUp(ctx context.Context, input input.SignUpInput) error {
	// 이메일 중복 검사
	exists, err := service.accountRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return status.NewError(http.StatusInternalServerError, "failed to check email existence")
	}
	if exists {
		return status.NewError(http.StatusConflict, "email already exists")
	}

	// 인증된 사용자 여부 검증
	// authentication, err := service.authenticationRepo.FindByEmail(ctx, input.Email)
	// if err != nil {
	// 	return status.NewError(http.StatusInternalServerError, "find authentication by email failed")
	// }
	// if authentication == nil {
	// 	return status.NewError(http.StatusNotFound, "authentication not found")
	// }
	// if !authentication.IsAuthenticated {
	// 	return status.NewError(http.StatusUnauthorized, "authentication not authenticated")
	// }

	// 패스워드 인코딩
	encodedPassword, err := service.utilPassword.EncodePassword(input.Password)
	if err != nil {
		return status.NewError(http.StatusInternalServerError, "password encode error")
	}

	account := &model.Account{
		Email:      input.Email,
		Password:   encodedPassword,
		Grade:      extractGrade(input.Email),
		Name:       input.Name,
		Major:      constant.Major(input.Major),
		Gender:     constant.Gender(input.Gender),
		ProfileURL: nil,
		Authority:  constant.ROLE_STUDENT,
		CreatedAt:  time.Now(),
	}
	return service.accountRepo.SaveAccount(ctx, account)
}

func (service *AuthService) SignIn(ctx context.Context, input input.SignInInput) (output.TokenOutput, error) {
	account, err := service.accountRepo.FindByEmail(ctx, input.Email)
	if err != nil || account == nil {
		return output.TokenOutput{}, status.NewError(http.StatusNotFound, "not found account")
	}

	isValidPassword, err := service.utilPassword.IsPasswordMatch(input.Password, account.Password)
	if err != nil || !isValidPassword {
		return output.TokenOutput{}, status.NewError(http.StatusUnauthorized, "mis match password")
	}

	tokenOutput, err := service.tokenAdapter.GenerateToken(ctx, account.ID, account.Authority)
	if err != nil {
		return output.TokenOutput{}, status.NewError(http.StatusInternalServerError, "token generate error")
	}

	return tokenOutput, nil
}

func (service *AuthService) TokenReissue(ctx context.Context, refreshToken string) (output.TokenOutput, error) {
	parsedRefreshToken, err := service.tokenParser.ParseRefreshToken(refreshToken)
	if err != nil {
		return output.TokenOutput{}, status.NewError(http.StatusBadRequest, "request header refresh token is error")
	}

	refreshTokenDomain, err := service.refreshTokenRepo.FindRefreshTokenByRefreshToken(ctx, parsedRefreshToken)
	if err != nil {
		return output.TokenOutput{}, status.NewError(http.StatusNotFound, "find refresh token by refresh token error")
	}

	if service.accountRepo == nil {
		return output.TokenOutput{}, status.NewError(http.StatusInternalServerError, "accountRepo is not initialized")
	}
	if refreshTokenDomain.AccountID == primitive.NilObjectID {
		return output.TokenOutput{}, status.NewError(http.StatusBadRequest, "invalid AccountID in refresh token")
	}

	fmt.Println("refresh token domain account id", refreshTokenDomain.AccountID)
	accountDomain, err := service.accountRepo.FindByAccountID(ctx, refreshTokenDomain.AccountID)
	if err != nil {
		return output.TokenOutput{}, status.NewError(http.StatusNotFound, "account not found")
	}

	tokenOutput, err := service.tokenAdapter.GenerateToken(ctx, accountDomain.ID, accountDomain.Authority)
	if err != nil {
		return output.TokenOutput{}, status.NewError(http.StatusInternalServerError, "token generation failed")
	}

	service.refreshTokenRepo.DeleteRefreshToken(ctx, refreshTokenDomain)

	return tokenOutput, nil
}

func (service *AuthService) SendAuthEmail(ctx context.Context, input input.SendEmaiInput) error {
	verificationCode := generateVerificationCode()
	fmt.Println(verificationCode)
	emailBody := fmt.Sprintf(
		`<p>인증번호 : <strong>%s</strong></p>`,
		verificationCode,
	)

	success, err := email.SendEmailSMTP(input.Email, emailBody, verificationCode)
	if err != nil {
		return status.NewError(http.StatusInternalServerError, "failed to send email")
	}
	if !success {
		return status.NewError(http.StatusInternalServerError, "email not sent")
	}

	exists, err := service.authenticationRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return status.NewError(http.StatusInternalServerError, "failed to check if email exists")
	}

	if exists {
		authentication, err := service.authenticationRepo.FindByEmail(ctx, input.Email)
		if err != nil {
			return status.NewError(http.StatusInternalServerError, "failed to find authentication")
		}

		// 이메일 전송 횟수 초과
		if authentication.AttemptCount > 5 {
			return status.NewError(http.StatusTooManyRequests, "many email request (over 5 times)")
		}

		// 전송횟수 증가
		authentication.AttemptCount++
		err = service.authenticationRepo.SaveAuthentication(ctx, authentication)
		if err != nil {
			return status.NewError(http.StatusInternalServerError, "failed to save updated authentication")
		}
		return nil
	}

	// 이메일을 처음 전송하는 경우
	authentication := &model.Authentication{
		Email:           input.Email,
		AttemptCount:    1,
		AuthCodeCount:   0,
		IsAuthenticated: false,
		ExpiredAt:       time.Now().Add(5 * time.Minute),
	}

	service.authenticationRepo.SaveAuthentication(ctx, authentication)

	authCode := &model.AuthCode{
		Email:     input.Email,
		AuthCode:  verificationCode,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	service.authCodeRepo.SaveAuthCode(ctx, authCode)

	return nil
}

func (service *AuthService) VerifyAuthCode(ctx context.Context, email string, authCode string) error {
	authCodeDomain, err := service.authCodeRepo.FindByEmail(ctx, email)
	if err != nil {
		return status.NewError(http.StatusNotFound, "auth code domain not found")
	}

	authentication, err := service.authenticationRepo.FindByEmail(ctx, email)
	if err != nil {
		return status.NewError(http.StatusNotFound, "authentication domain not found")
	}
	if authentication == nil {
		return status.NewError(http.StatusNotFound, "authentication domain not found")
	}

	if authentication.AuthCodeCount > 5 {
		return status.NewError(http.StatusTooManyRequests, "too many verify auth code requests (over 5 times)")
	}

	if authCodeDomain.AuthCode != authCode {
		authentication.AuthCodeCount++
		service.authenticationRepo.SaveAuthentication(ctx, authentication)
		return status.NewError(http.StatusUnauthorized, "auth code not match")
	}

	authentication.IsAuthenticated = true
	service.authenticationRepo.SaveAuthentication(ctx, authentication)

	return nil
}

func generateVerificationCode() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func extractGrade(email string) int {
	var grade int
	switch email[2] {
	case '2':
		grade = 6
	case '3':
		grade = 7
	case '4':
		grade = 8
	}
	return grade
}
