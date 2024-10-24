package service

import (
	"GOMS-BACKEND-GO/global/config"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StudentCouncilService struct {
	outingUUIDRepo            model.OutingUUIDRepository
	accountRepo               model.AccountRepository
	blackListRepo             model.BlackListRepository
	outingBlackListProperties *config.OutingBlackListProperties
	outingRepo                model.OutingRepository
}

func NewStudentCouncilService(outingUUIDRepo model.OutingUUIDRepository, accountRepo model.AccountRepository, blackListRepo model.BlackListRepository, outingBlackListProperties *config.OutingBlackListProperties, outingRepo model.OutingRepository) model.StudentCouncilUseCase {
	return &StudentCouncilService{
		outingUUIDRepo:            outingUUIDRepo,
		accountRepo:               accountRepo,
		blackListRepo:             blackListRepo,
		outingBlackListProperties: outingBlackListProperties,
		outingRepo:                outingRepo,
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

func (service *StudentCouncilService) AddBlackList(ctx context.Context, accountID uint64) error {
	expiration := time.Duration(service.outingBlackListProperties.OutingBlackListExp) * time.Second

	blackList := &model.BlackList{
		AccountID: accountID,
		ExpiredAt: int64(expiration.Seconds()),
	}
	service.blackListRepo.SaveBlackList(ctx, blackList)

	return nil
}

func (service *StudentCouncilService) ExcludeBlackList(ctx context.Context, accountID uint64) error {
	outingBlackList, err := service.blackListRepo.FindBlackListByAccountID(ctx, accountID)
	if err != nil {
		return nil
	}
	if outingBlackList == nil {
		return fmt.Errorf("blacklist not found for account ID: %d", accountID)
	}

	service.blackListRepo.DeleteBlackList(ctx, outingBlackList)
	return nil
}

func (service *StudentCouncilService) DeleteOutingStudent(ctx context.Context, accountID uint64) error {
	exists, err := service.outingRepo.ExistsOutingByAccountID(ctx, accountID)
	if err != nil {
		return nil
	}
	if !exists {
		return fmt.Errorf("not outing student")
	}

	deleteErr := service.outingRepo.DeleteOutingByAccountID(ctx, accountID)
	if deleteErr != nil {
		return nil
	}

	return nil

}
