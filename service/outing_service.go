package service

import (
	"GOMS-BACKEND-GO/global/error/status"
	"GOMS-BACKEND-GO/global/util"
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/output"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OutingService struct {
	outingRepo     model.OutingRepository
	accountRepo    model.AccountRepository
	outingUUIDRepo model.OutingUUIDRepository
}

func NewOutingService(outingRepo model.OutingRepository, accountRepo model.AccountRepository, outingUUIDRepo model.OutingUUIDRepository) *OutingService {
	return &OutingService{
		outingRepo:     outingRepo,
		accountRepo:    accountRepo,
		outingUUIDRepo: outingUUIDRepo,
	}
}

func (service *OutingService) OutingStudent(c *gin.Context, ctx context.Context, outingUUID uuid.UUID) error {
	accountID, err := util.GetCurrentAccountID(c)
	if err != nil {
		return err
	}

	// 유효한 외출 UUID 인지 검증
	existsOutingUUID, err := service.outingUUIDRepo.ExistsByOutingUUID(ctx, outingUUID)

	if err != nil {
		return status.NewError(http.StatusInternalServerError, "failed to outing UUID")
	}
	if !existsOutingUUID {
		return status.NewError(http.StatusBadRequest, "failed to outing UUID")
	}

	// account id를 기반으로 account 추출
	account, err := service.accountRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return status.NewError(http.StatusInternalServerError, "failed to find account")
	}

	// 이미 외출한 학생인지 검증
	existsOutingStudent, err := service.outingRepo.ExistsOutingByAccountID(ctx, account.ID)

	if existsOutingStudent {
		// 이미 외출한 학생이라면 복귀를 위한 QR 인식으로 외출 명단에서 삭제
		service.outingRepo.DeleteOutingByAccountID(ctx, accountID)
	} else {
		outing := &model.Outing{
			Account:   account,
			CreatedAt: time.Now(),
		}
		// 아직 외출을 하지 않은 학생이라면 외출자 명단에 추가
		service.outingRepo.SaveOutingStudnet(ctx, outing)

	}

	return err
}

func (service *OutingService) FindAllOutingStudent(ctx context.Context) ([]output.OutingStudentOutput, error) {
	outings, err := service.outingRepo.FindAllOuting(ctx)
	if err != nil {
		return nil, status.NewError(http.StatusInternalServerError, "failed to find account")
	}

	var outingStudentOutputs []output.OutingStudentOutput

	for _, outing := range outings {
		account, err := service.accountRepo.FindByAccountID(ctx, outing.AccountID)

		if err != nil {
			return nil, status.NewError(http.StatusInternalServerError, "failed to find account")
		}

		outingStudentOutput := output.OutingStudentOutput{
			AccountID:   account.ID,
			Name:        account.Name,
			Grade:       account.Grade,
			Major:       account.Major,
			Gender:      account.Gender,
			ProfileURL:  account.ProfileURL,
			CreatedTime: outing.CreatedAt,
		}

		outingStudentOutputs = append(outingStudentOutputs, outingStudentOutput)
	}

	return outingStudentOutputs, nil
}

func (service *OutingService) CountAllOutingStudent(ctx context.Context) (int, error) {
	outings, err := service.outingRepo.FindAllOuting(ctx)
	if err != nil {
		return 0, status.NewError(http.StatusInternalServerError, "failed to find account")
	}

	return len(outings), err

}

func (service *OutingService) SearchOutingStudent(ctx context.Context, name string) ([]output.OutingStudentOutput, error) {
	outings, err := service.outingRepo.FindByOutingAccountNameContaining(ctx, name)

	if err != nil {
		return nil, status.NewError(http.StatusInternalServerError, "failed to find account")
	}

	var outingStudentOutputs []output.OutingStudentOutput

	for _, outing := range outings {
		account, err := service.accountRepo.FindByAccountID(ctx, outing.AccountID)

		if err != nil {
			return nil, status.NewError(http.StatusInternalServerError, "failed to find account")
		}

		outingStudentOutput := output.OutingStudentOutput{
			AccountID:   account.ID,
			Name:        account.Name,
			Grade:       account.Grade,
			Major:       account.Major,
			Gender:      account.Gender,
			ProfileURL:  account.ProfileURL,
			CreatedTime: outing.CreatedAt,
		}

		outingStudentOutputs = append(outingStudentOutputs, outingStudentOutput)
	}

	return outingStudentOutputs, nil

}
