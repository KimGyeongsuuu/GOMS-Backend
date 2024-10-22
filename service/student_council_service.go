package service

import (
	"GOMS-BACKEND-GO/model"
	"context"

	"github.com/google/uuid"
)

type StudentCouncilService struct {
	outingUUIDRepo model.OutingUUIDRepository
}

func NewStudentCouncilService(outingUUIDRepo model.OutingUUIDRepository) model.StudentCouncilUseCase {
	return &StudentCouncilService{
		outingUUIDRepo: outingUUIDRepo,
	}
}

func (service *StudentCouncilService) CreateOuting(ctx context.Context) (uuid.UUID, error) {

	newUUID, err := service.outingUUIDRepo.CreateOutingUUID(ctx)
	if err != nil {
		return uuid.UUID{}, err
	}

	return newUUID, nil
}
