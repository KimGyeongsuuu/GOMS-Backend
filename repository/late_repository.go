package repository

import (
	"GOMS-BACKEND-GO/model"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type LateRepository struct {
	db *gorm.DB
}

func NewLateRepository(db *gorm.DB) *LateRepository {
	return &LateRepository{
		db: db,
	}
}

func (repository *LateRepository) FindTop3ByOrderByAccountDesc(ctx context.Context) ([]model.Late, error) {
	var lates []model.Late

	err := repository.db.WithContext(ctx).
		Preload("Account").
		Model(&model.Late{}).
		Select("lates.*, COUNT(accounts.id) as account_count").
		Joins("JOIN accounts ON accounts.id = lates.account_id").
		Group("accounts.id").
		Order("account_count DESC").
		Limit(3).
		Find(&lates).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch top 3 lates: %w", err)
	}

	return lates, nil
}

func (repository *LateRepository) FindLateByCreatedAt(ctx context.Context, date time.Time) ([]model.Late, error) {
	var lates []model.Late

	startOfDay := date.Truncate(24 * time.Hour)
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := repository.db.WithContext(ctx).
		Preload("Account").
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Find(&lates).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch late students: %w", err)
	}

	return lates, nil
}
