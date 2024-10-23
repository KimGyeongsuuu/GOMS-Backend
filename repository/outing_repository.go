package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"

	"gorm.io/gorm"
)

type OutingRepository struct {
	db *gorm.DB
}

func NewOutingRepository(db *gorm.DB) *OutingRepository {
	return &OutingRepository{
		db: db,
	}
}

func (repository *OutingRepository) SaveOutingStudnet(ctx context.Context, outing *model.Outing) error {
	result := repository.db.WithContext(ctx).Create(outing)
	return result.Error
}

func (repository *OutingRepository) ExistsOutingByAccountID(ctx context.Context, accountID uint64) (bool, error) {
	var count int64
	result := repository.db.WithContext(ctx).Model(&model.Outing{}).Where("account_id = ?", accountID).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (repository *OutingRepository) DeleteOutingByAccountID(ctx context.Context, accountID uint64) error {
	result := repository.db.WithContext(ctx).Where("account_id = ?", accountID).Delete(&model.Outing{})
	return result.Error
}
