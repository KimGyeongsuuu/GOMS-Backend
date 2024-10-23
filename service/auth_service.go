package service

import (
	"GOMS-BACKEND-GO/global/auth/jwt"
	"GOMS-BACKEND-GO/global/util"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/constant"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"errors"
	"time"
)

type AuthService struct {
	accountRepo  model.AccountRepository
	tokenAdapter *jwt.GenerateTokenAdapter
}

func NewAuthService(accountRepo model.AccountRepository, tokenAdapter *jwt.GenerateTokenAdapter) model.AuthUseCase {
	return &AuthService{
		accountRepo:  accountRepo,
		tokenAdapter: tokenAdapter,
	}
}

func (service *AuthService) SignUp(ctx context.Context, input *input.SignUpInput) error {

	exists, err := service.accountRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return errors.New("failed to check email existence")
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
