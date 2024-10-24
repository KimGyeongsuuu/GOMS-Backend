package service

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"fmt"
)

type LateService struct {
	lateRepo model.LateRepository
}

func NewLateService(lateRepo model.LateRepository) *LateService {
	return &LateService{
		lateRepo: lateRepo,
	}
}

func (service *LateService) GetTop3LateStudent(ctx context.Context) ([]output.LateTop3Output, error) {
	lates, err := service.lateRepo.FindTop3ByOrderByAccountDesc(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Fetched lates not if: %+v\n", lates)

	var outputList []output.LateTop3Output
	fmt.Println(len(lates))
	for _, late := range lates {
		if late.Account == nil {
			fmt.Println(late.Account)
			return nil, fmt.Errorf("late.Account is nil for late ID: %d", late.Account.ID)
		}

		output := output.LateTop3Output{
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
