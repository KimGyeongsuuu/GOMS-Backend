package service

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"

	"github.com/google/uuid"
)

type StudentCouncilService struct {
	outingUUIDRepo model.OutingUUIDRepository
	accountRepo    model.AccountRepository
}

func NewStudentCouncilService(outingUUIDRepo model.OutingUUIDRepository, accountRepo model.AccountRepository) model.StudentCouncilUseCase {
	return &StudentCouncilService{
		outingUUIDRepo: outingUUIDRepo,
		accountRepo:    accountRepo,
	}
}

func (service *StudentCouncilService) CreateOuting(ctx context.Context) (uuid.UUID, error) {

	outingUUID, err := service.outingUUIDRepo.CreateOutingUUID(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}

	return outingUUID, nil
}

func (service *StudentCouncilService) FindAllAccount(ctx context.Context) ([]output.AccountOutput, error) {

	accounts, err := service.accountRepo.FindAllAccount(ctx)

	var accountOutputs []output.AccountOutput

	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		accountOutput := output.AccountOutput{
			AccountID:   account.ID,
			Name:        account.Name,
			Major:       account.Major,
			Grade:       account.Grade,
			ProfileURL:  account.ProfileURL,
			Authority:   account.Authority,
			IsBlackList: false,
		}

		accountOutputs = append(accountOutputs, accountOutput)

	}

	return accountOutputs, err
}

func (service *StudentCouncilService) SearchAccount(ctx context.Context, accountInput *input.SearchAccountInput) ([]output.AccountOutput, error) {

	accounts, err := service.accountRepo.FindByAccountByStudentInfo(ctx, accountInput)

	if err != nil {
		return nil, err
	}
	var accountOutputs []output.AccountOutput

	for _, account := range accounts {
		accountOutput := output.AccountOutput{
			AccountID:   account.ID,
			Name:        account.Name,
			Major:       account.Major,
			Grade:       account.Grade,
			ProfileURL:  account.ProfileURL,
			Authority:   account.Authority,
			IsBlackList: false,
		}

		accountOutputs = append(accountOutputs, accountOutput)

	}
	return accountOutputs, err
}

func (service *StudentCouncilService) UpdateAccountAuthority(ctx context.Context, authorityInput *input.UpdateAccountAuthorityInput) error {

	err := service.accountRepo.UpdateAccountAuthority(ctx, authorityInput)

	if err != nil {
		return err
	}

	return nil

}
