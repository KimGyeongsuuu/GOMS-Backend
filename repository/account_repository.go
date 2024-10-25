package repository

import (
	"GOMS-BACKEND-GO/model"
	"GOMS-BACKEND-GO/model/data/input"
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

func (repository *AccountRepository) SaveAccount(ctx context.Context, account *model.Account) error {
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

func (repository *AccountRepository) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account
	result := repository.db.WithContext(ctx).Where("email = ?", email).First(&account)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &account, nil
}

func (repository *AccountRepository) FindByAccountID(ctx context.Context, accountID uint64) (*model.Account, error) {
	var account model.Account
	result := repository.db.WithContext(ctx).Where("id = ?", accountID).First(&account)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &account, nil
}

func (repository *AccountRepository) FindAllAccount(ctx context.Context) ([]model.Account, error) {
	var accounts []model.Account
	result := repository.db.WithContext(ctx).Find(&accounts)

	return accounts, result.Error
}

func (repository *AccountRepository) FindByAccountByStudentInfo(ctx context.Context, searchAccountInput *input.SearchAccountInput) ([]model.Account, error) {
	var accounts []model.Account
	query := repository.db.WithContext(ctx).Model(&model.Account{})

	if searchAccountInput.Grade != nil {
		query = query.Where("grade = ?", *searchAccountInput.Grade)
	}
	if searchAccountInput.Gender != nil {
		query = query.Where("gender = ?", *searchAccountInput.Gender)
	}
	if searchAccountInput.Name != nil {
		query = query.Where("name LIKE ?", "%"+*searchAccountInput.Name+"%")
	}
	if searchAccountInput.Authority != nil {
		query = query.Where("authority = ?", *searchAccountInput.Authority)
	}
	if searchAccountInput.Major != nil {
		query = query.Where("major = ?", *searchAccountInput.Major)
	}

	query = query.Order("grade ASC")

	result := query.Find(&accounts)

	if result.Error != nil {
		return nil, result.Error
	}

	return accounts, nil
}

func (repository *AccountRepository) UpdateAccountAuthority(ctx context.Context, authorityInput *input.UpdateAccountAuthorityInput) error {
	result := repository.db.WithContext(ctx).
		Model(&model.Account{}).
		Where("id = ?", authorityInput.AccountID).
		Update("authority", authorityInput.Authority)

	if result.Error != nil {
		return result.Error
	}

	return nil

}

func (repository *AccountRepository) DeleteAccount(ctx context.Context, account *model.Account) error {
	result := repository.db.WithContext(ctx).Delete(account)
	return result.Error
}
