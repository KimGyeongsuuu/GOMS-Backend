package model

import "context"

type BlackList struct {
	AccountID uint64
	ExpiredAt int64
}

type BlackListRepository interface {
	SaveBlackList(ctx context.Context, blackList *BlackList) error
	DeleteBlackList(ctx context.Context, blackList *BlackList) error
	FindBlackListByAccountID(ctx context.Context, accountID uint64) (*BlackList, error)
}
