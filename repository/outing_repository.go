package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (repository *OutingRepository) SaveOutingStudent(ctx context.Context, outing *model.Outing) error {
	result := repository.db.WithContext(ctx).Create(outing)
	return result.Error
}

func (repository *OutingRepository) ExistsOutingByAccountID(ctx context.Context, accountID primitive.ObjectID) (bool, error) {
	var count int64
	result := repository.db.WithContext(ctx).Model(&model.Outing{}).Where("account_id = ?", accountID).Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

func (repository *OutingRepository) DeleteOutingByAccountID(ctx context.Context, accountID primitive.ObjectID) error {
	result := repository.db.WithContext(ctx).Where("account_id = ?", accountID).Delete(&model.Outing{})
	return result.Error
}

func (repository *OutingRepository) FindAllOuting(ctx context.Context) ([]model.Outing, error) {
	var outings []model.Outing
	result := repository.db.WithContext(ctx).Find(&outings)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(outings) == 0 {
		return nil, errors.New("not outings found")
	}

	return outings, nil
}

func (repository *OutingRepository) FindByOutingAccountNameContaining(ctx context.Context, name string) ([]model.Outing, error) {
	var outings []model.Outing
	result := repository.db.WithContext(ctx).
		Preload("Account").
		Joins("JOIN accounts ON accounts.id = outings.account_id").
		Where("accounts.name LIKE ?", "%"+name+"%").
		Find(&outings)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(outings) == 0 {
		return nil, errors.New("not found account")
	}

	return outings, nil
}
