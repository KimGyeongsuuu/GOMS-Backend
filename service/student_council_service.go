package service

import (
	"GOMS-BACKEND-GO/global/config"
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
	"fmt"
	"net/http"

	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentCouncilService struct {
	outingUUIDRepo model.OutingUUIDRepository
	accountRepo    model.AccountRepository
	blackListRepo  model.BlackListRepository
	outingConfig   *config.OutingConfig
	outingRepo     model.OutingRepository
	lateRepo       model.LateRepository
}

func NewStudentCouncilService(
	outingUUIDRepo model.OutingUUIDRepository,
	accountRepo model.AccountRepository,
	blackListRepo model.BlackListRepository,
	outingConfig *config.OutingConfig,
	outingRepo model.OutingRepository,
	lateRepo model.LateRepository,
) model.StudentCouncilUseCase {
	return &StudentCouncilService{
		outingUUIDRepo: outingUUIDRepo,
		accountRepo:    accountRepo,
		blackListRepo:  blackListRepo,
		outingConfig:   outingConfig,
		outingRepo:     outingRepo,
		lateRepo:       lateRepo,
	}
}

func (service *StudentCouncilService) CreateOuting(ctx context.Context) (uuid.UUID, error) {

	outingUUID, err := service.outingUUIDRepo.CreateOutingUUID(ctx)
	if err != nil {
		return uuid.UUID{}, status.NewError(http.StatusInternalServerError, "create outing uuid error")
	}

	return outingUUID, nil
}

func (service *StudentCouncilService) FindAllAccount(ctx context.Context) ([]output.AccountOutput, error) {

	accounts, err := service.accountRepo.FindAllAccount(ctx)

	var accountOutputs []output.AccountOutput

	if err != nil {
		return nil, nil
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
		return nil, status.NewError(http.StatusInternalServerError, "search account error")
	}
	var accountOutputs []output.AccountOutput

	for _, account := range accounts {

		isBlackList, err := service.blackListRepo.ExistsByAccountID(ctx, account.ID)
		if err != nil {
			return nil, status.NewError(http.StatusInternalServerError, "black list error")
		}
		accountOutput := output.AccountOutput{
			AccountID:   account.ID,
			Name:        account.Name,
			Major:       account.Major,
			Grade:       account.Grade,
			ProfileURL:  account.ProfileURL,
			Authority:   account.Authority,
			IsBlackList: isBlackList,
		}

		accountOutputs = append(accountOutputs, accountOutput)

	}
	return accountOutputs, err
}

func (service *StudentCouncilService) UpdateAccountAuthority(ctx context.Context, authorityInput *input.UpdateAccountAuthorityInput) error {

	err := service.accountRepo.UpdateAccountAuthority(ctx, authorityInput)

	if err != nil {
		return status.NewError(http.StatusInternalServerError, "update account authority error")
	}

	return nil

}

func (service *StudentCouncilService) AddBlackList(ctx context.Context, accountID primitive.ObjectID) error {
	expiration := time.Duration(service.outingConfig.OutingBlacklistExp) * time.Second

	blackList := &model.BlackList{
		AccountID: accountID,
		ExpiredAt: int64(expiration.Seconds()),
	}
	service.blackListRepo.SaveBlackList(ctx, blackList)

	return nil
}

func (service *StudentCouncilService) ExcludeBlackList(ctx context.Context, accountID primitive.ObjectID) error {
	outingBlackList, err := service.blackListRepo.FindBlackListByAccountID(ctx, accountID)
	if err != nil {
		return status.NewError(http.StatusInternalServerError, fmt.Sprintf("blacklist not found for account ID: %d", accountID))
	}
	if outingBlackList == nil {
		return status.NewError(http.StatusInternalServerError, fmt.Sprintf("blacklist not found for account ID: %d", accountID))
	}

	service.blackListRepo.DeleteBlackList(ctx, outingBlackList)
	return nil
}

func (service *StudentCouncilService) DeleteOutingStudent(ctx context.Context, accountID primitive.ObjectID) error {
	exists, err := service.outingRepo.ExistsOutingByAccountID(ctx, accountID)
	if err != nil {
		return status.NewError(http.StatusInternalServerError, "delete outing student error")
	}
	if !exists {
		return status.NewError(http.StatusNotFound, "not outing student")
	}

	deleteErr := service.outingRepo.DeleteOutingByAccountID(ctx, accountID)
	if deleteErr != nil {
		return nil
	}

	return nil

}

func (service *StudentCouncilService) FindLateStudentByDate(ctx context.Context, date time.Time) ([]output.LateOutput, error) {
	lates, err := service.lateRepo.FindLateByCreatedAt(ctx, date)

	if err != nil {
		return []output.LateOutput{}, nil
	}

	var outputList []output.LateOutput
	for _, late := range lates {

		accountDomain, err := service.accountRepo.FindByAccountID(ctx, late.AccountID)

		if err != nil {
			return nil, status.NewError(http.StatusInternalServerError, "find by account id error")
		}

		outputItem := output.LateOutput{
			AccountID:  late.AccountID,
			Name:       accountDomain.Name,
			Major:      accountDomain.Major,
			Grade:      accountDomain.Grade,
			Gender:     accountDomain.Gender,
			ProfileURL: accountDomain.ProfileURL,
		}
		outputList = append(outputList, outputItem)
	}

	return outputList, nil

}
