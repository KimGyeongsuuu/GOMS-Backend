package service

import (
	"GOMS-BACKEND-GO/global/auth/jwt"
	"GOMS-BACKEND-GO/global/email"
	"GOMS-BACKEND-GO/global/util"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

type AuthService struct {
	accountRepo        model.AccountRepository
	tokenAdapter       *jwt.GenerateTokenAdapter
	tokenParser        *jwt.TokenParser
	refreshTokenRepo   model.RefreshTokenRepository
	authenticationRepo model.AuthenticationRepository
	authCodeRepo       model.AuthCodeRepository
}

func NewAuthService(
	accountRepo model.AccountRepository,
	tokenAdapter *jwt.GenerateTokenAdapter,
	refreshTokenRepo model.RefreshTokenRepository,
	tokenParser *jwt.TokenParser,
	authenticationRepo model.AuthenticationRepository,
	authCodeRepo model.AuthCodeRepository,
) model.AuthUseCase {
	return &AuthService{
		accountRepo:        accountRepo,
		tokenAdapter:       tokenAdapter,
		refreshTokenRepo:   refreshTokenRepo,
		tokenParser:        tokenParser,
		authenticationRepo: authenticationRepo,
		authCodeRepo:       authCodeRepo,
	}
}

func (service *AuthService) SignUp(ctx context.Context, input *input.SignUpInput) error {

	exists, err := service.accountRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return errors.New("failed to check email existence")
	}

	authentication, err := service.authenticationRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return errors.New("find authentication by email is failed ")
	}
	if authentication == nil {
		return errors.New("authentication not found")
	}
	if !authentication.IsAuthenticated {
		return errors.New("authentication not found")
	}

	if exists {
		return errors.New("email already exists")
	}

	encodedPassword, err := util.EncodePassword(input.Password)
	if err != nil {
		return err
	}

	account := &model.Account{
		Email:      input.Email,
		Password:   encodedPassword,
		Grade:      *extractGrade(input.Email),
		Name:       input.Name,
		Major:      constant.Major(input.Major),
		Gender:     constant.Gender(input.Gender),
		ProfileURL: nil,
		Authority:  constant.ROLE_STUDENT,
		CreatedAt:  time.Now(),
	}

	return service.accountRepo.SaveAccount(ctx, account)

}

func (service *AuthService) SignIn(ctx context.Context, input *input.SignInInput) (output.TokenOutput, error) {

	account, err := service.accountRepo.FindByEmail(ctx, input.Email)

	if err != nil {
		return output.TokenOutput{}, err
	}

	if account == nil {
		return output.TokenOutput{}, errors.New("not found account")
	}

	isValidPassword, err := util.IsPasswordMatch(input.Password, account.Password)
	if err != nil {
		return output.TokenOutput{}, err
	}

	if !isValidPassword {
		return output.TokenOutput{}, errors.New("mis match password")
	}

	tokenOutput, err := service.tokenAdapter.GenerateToken(ctx, account.ID, account.Authority)

	if err != nil {
		return output.TokenOutput{}, err
	}

	return tokenOutput, nil

}

func (service *AuthService) TokenReissue(ctx context.Context, refreshToken string) (output.TokenOutput, error) {
	parsedRefreshToken, err := service.tokenParser.ParseRefreshToken(refreshToken)
	if err != nil {
		return output.TokenOutput{}, fmt.Errorf("request header refresh token is error")
	}

	refreshTokenDomain, err := service.refreshTokenRepo.FindRefreshTokenByRefreshToken(ctx, parsedRefreshToken)
	if err != nil {
		return output.TokenOutput{}, fmt.Errorf("find refresh token by refresh token method is eror")
	}

	if service.accountRepo == nil {
		return output.TokenOutput{}, fmt.Errorf("accountRepo is not initialized")
	}

	if refreshTokenDomain.AccountID == 0 {
		return output.TokenOutput{}, fmt.Errorf("invalid AccountID in refresh token")
	}

	accountDomain, err := service.accountRepo.FindByAccountID(ctx, refreshTokenDomain.AccountID)
	if err != nil {
		return output.TokenOutput{}, err
	}
	tokenOutput, err := service.tokenAdapter.GenerateToken(ctx, accountDomain.ID, accountDomain.Authority)
	if err != nil {
		return output.TokenOutput{}, err
	}
	service.refreshTokenRepo.DeleteRefreshToken(ctx, refreshTokenDomain)
	return tokenOutput, nil

}

func (service *AuthService) SendAuthEmail(ctx context.Context, input *input.SendEmaiInput) error {
	templateName := "verification.html"
	verificationCode := generateVerificationCode()

	success, err := email.SendEmailSMTP(input.Email, nil, templateName, verificationCode)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if !success {
		return fmt.Errorf("email not sent")
	}

	exists, err := service.authenticationRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return fmt.Errorf("failed to exists by email method")
	}

	if exists {
		authentication, err := service.authenticationRepo.FindByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("failed to find by email method")
		}
		if authentication.AttemptCount > 5 {
			return fmt.Errorf("many email request (over 5 times)")
		}

		authentication.AttemptCount++

		err = service.authenticationRepo.SaveAuthentication(ctx, authentication)
		if err != nil {
			return fmt.Errorf("failed to save updated authentication: %v", err)
		}

		return nil
	}

	if !exists {
		authentication := &model.Authentication{
			Email:           input.Email,
			AttemptCount:    1,
			AuthCodeCount:   0,
			IsAuthenticated: false,
			ExpiredAt:       time.Now().Add(5 * time.Minute),
		}

		service.authenticationRepo.SaveAuthentication(ctx, authentication)

	}

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
		return fmt.Errorf("auth code domain not found")
	}

	authentication, err := service.authenticationRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("authentication domain not found")
	}

	if authentication == nil {
		return fmt.Errorf("authentication domain not found")
	}

	if authentication.AuthCodeCount > 5 {
		return fmt.Errorf("to many verify auth code request (over 5 times)")
	}

	if authCodeDomain.AuthCode != authCode {
		authentication.AuthCodeCount++
		service.authenticationRepo.SaveAuthentication(ctx, authentication)
		return fmt.Errorf("auth code not match")
	}

	authentication.IsAuthenticated = true
	service.authenticationRepo.SaveAuthentication(ctx, authentication)
	return nil
}

func generateVerificationCode() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func extractGrade(email string) *int {
	if len(email) < 3 {
		return nil
	}

	var grade int
	switch email[2] {
	case '2':
		grade = 6
	case '3':
		grade = 7
	case '4':
		grade = 8
	default:
		return nil
	}

	return &grade
}
