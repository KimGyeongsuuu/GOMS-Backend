package service

import (
	"GOMS-BACKEND-GO/global/util"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
	"context"
	"errors"
	"time"
)

type AuthService struct {
	accountRepo model.AccountRepository
}

func NewAuthService(accountRepo model.AccountRepository) model.AuthUseCase {
	return &AuthService{
		accountRepo: accountRepo,
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
		Major:      model.Major(input.Major),
		Gender:     model.Gender(input.Gender),
		ProfileURL: nil,
		Authority:  model.ROLE_STUDENT,
		CreatedAt:  time.Now(),
	}

	return service.accountRepo.CreateAccount(ctx, account)

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
