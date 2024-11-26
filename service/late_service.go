package service

import (
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"fmt"
	"net/http"
)

type LateService struct {
	lateRepo model.LateRepository
}

func NewLateService(lateRepo model.LateRepository) *LateService {
	return &LateService{
		lateRepo: lateRepo,
	}
}

func (service *LateService) GetTop3LateStudent(ctx context.Context) ([]output.LateOutput, error) {
	lates, err := service.lateRepo.FindTop3ByOrderByAccountDesc(ctx)
	if err != nil {
		return nil, err
	}

	var outputList []output.LateOutput
	fmt.Println(len(lates))
	for _, late := range lates {
		if late.Account == nil {
			return nil, status.NewError(http.StatusNotFound, fmt.Sprintf("late.Account is nil for late ID: %d", late.Account.ID))
		}

		output := output.LateOutput{
			AccountID:  late.Account.ID,
			Name:       late.Account.Name,
			Major:      late.Account.Major,
			Grade:      late.Account.Grade,
			Gender:     late.Account.Gender,
			ProfileURL: late.Account.ProfileURL,
		}
		outputList = append(outputList, output)
	}
	return outputList, nil
}
