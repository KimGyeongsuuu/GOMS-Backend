package service

import (
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"net/http"
)

type LateService struct {
	lateRepo    model.LateRepository
	accountRepo model.AccountRepository
}

func NewLateService(lateRepo model.LateRepository, accountRepo model.AccountRepository) *LateService {
	return &LateService{
		lateRepo:    lateRepo,
		accountRepo: accountRepo,
	}
}

func (service *LateService) GetTop3LateStudent(ctx context.Context) ([]output.LateOutput, error) {
	lates, err := service.lateRepo.FindTop3ByOrderByAccountDesc(ctx)
	if err != nil {
		return nil, status.NewError(http.StatusInternalServerError, "find top 3 by order by account desc error")
	}

	var outputList []output.LateOutput
	for _, late := range lates {

		accountDomain, err := service.accountRepo.FindByAccountID(ctx, late.AccountID)

		if err != nil {
			return nil, status.NewError(http.StatusInternalServerError, "find by account id error")
		}

		output := output.LateOutput{
			AccountID:  accountDomain.ID,
			Name:       accountDomain.Name,
			Major:      accountDomain.Major,
			Grade:      accountDomain.Grade,
			Gender:     accountDomain.Gender,
			ProfileURL: accountDomain.ProfileURL,
		}
		outputList = append(outputList, output)
	}

	return outputList, nil
}
