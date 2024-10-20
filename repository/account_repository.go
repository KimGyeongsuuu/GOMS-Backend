package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (repository *AccountRepository) CreateAccount(ctx context.Context, account *model.Account) error {
	result := repository.db.WithContext(ctx).Create(account)
	return result.Error
}

func (repository *AccountRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := repository.db.Model(&model.Account{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
